package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/RedCoolBeans/crane/util/fs"
	g "github.com/RedCoolBeans/crane/util/git"
	log "github.com/RedCoolBeans/crane/util/logging"
	m "github.com/RedCoolBeans/crane/util/manifest"
	ssh "github.com/RedCoolBeans/crane/util/ssh"
	u "github.com/RedCoolBeans/crane/util/utils"
	"gopkg.in/libgit2/git2go.v23"
)

var (
	verbose *bool
	debug   *bool
)

func main() {
	cargo := flag.String("package", "", "Name of package to load")
	branch := flag.String("branch", "master", "Branch or version")
	destination := flag.String("destination", "/", "Destination for package on filesystem")
	repo := flag.String("repo", "https://git.cargos.io/", "URI of repository base")
	sshkey := flag.String("sshkey", "", "Path to SSH private key")
	sshpass := flag.String("sshpass", "", "SSH private key password")
	verbose = flag.Bool("verbose", false, "Enable verbose logging")
	debug = flag.Bool("debug", false, "Enable debugging (uses panic(), implies -verbose)")
	clean := flag.Bool("clean", false, "Remove crane after deployment")
	prefix := flag.String("prefix", "", "Prefix into the repository to the files")

	flag.Parse()

	// debug implies verbose
	if *debug {
		*verbose = true
	}

	if *clean {
		fs.CleanSelf(*verbose)
	}

	if len(strings.TrimSpace(*cargo)) < 1 {
		log.PrError("No package specified to load")
	}

	if err := fs.CanReadDir(*destination, "Destination directory"); err != nil {
		log.PrFatal(err.Error())
	}

	chain := m.InitDependencyChain(*cargo)

	// Everything is setup, hand-off to the main loop
	crane(*repo, *cargo, *branch, *prefix, *destination, *sshkey, *sshpass, &chain)

	if *clean {
		fs.CleanSelf(*verbose)
	}
}

func initGitOptions(sshOptions *ssh.SshOptions, branch string, repo string, cargo string) (*git.CloneOptions, string) {
	options := &git.CloneOptions{}
	options.CheckoutBranch = branch

	var cargoRepo string

	if sshOptions.Enabled {
		cargoRepo = sshOptions.Sshrepo

		// Resort to using local functions for the RemoteCallbacks.
		// This way they can resolve the sshOptions fields which would
		// otherwise be global and impossible to update when resolving
		// them for dependencies.
		var credentialsCB func(string, string, git.CredType) (git.ErrorCode, *git.Cred)
		credentialsCB = func(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
			ret, cred := git.NewCredSshKey(
				sshOptions.Sshuser,
				sshOptions.Sshpubkey,
				sshOptions.Sshkey,
				sshOptions.Sshpass)
			return git.ErrorCode(ret), &cred
		}

		var certificateCB func(*git.Certificate, bool, string) git.ErrorCode
		certificateCB = func(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
			return 0
		}

		rcbs := git.RemoteCallbacks{
			CredentialsCallback:      credentialsCB,
			CertificateCheckCallback: certificateCB,
		}

		fopts := &git.FetchOptions{
			RemoteCallbacks: rcbs,
		}

		options.FetchOptions = fopts
	} else {
		cargoRepo = fmt.Sprintf("%s/%s", repo, cargo)
	}

	return options, cargoRepo
}

// Main body, dispatched to after main() itself has finished parsing all flags.
func crane(repo string, cargo string, branch string, prefix string, destination string, sshkey string, sshpass string, chain *m.DependencyChain) {
	clonedir, err := fs.CreateTempDir()
	defer fs.CleanTempDir(clonedir)
	u.Check(err, false)
	log.PrVerbose(*verbose, "Using %s to store temporary files", clonedir)

	sshOptions := ssh.SshOptions{}
	sshOptions.Enabled = false

	if ssh.CanHandle(repo) {
		sshOptions.Enabled = true
		sshOptions.Sshkey = sshkey
		sshOptions.Sshpass = sshpass

		err = ssh.Init(&sshOptions, repo, cargo)
		u.Check(err, false)
	}

	options, cargoRepo := initGitOptions(&sshOptions, branch, repo, cargo)

	log.PrInfo("Fetching %s (%s)...", cargo, branch)
	headCommit, err := g.Clone(cargoRepo, branch, clonedir, *options)
	u.Check(err, false)

	log.PrVerbose(*verbose, "HEAD is at %s: %s", headCommit.Id(), headCommit.Summary())

	if err := g.RemoveDotGit(clonedir); err != nil {
		log.PrError(err.Error())
	}

	manifest := parseManifest(clonedir)
	log.PrInfo("Installing %s %s", manifest["name"], m.VersionString(manifest))

	parent := false
	dependencies := m.Dependencies(manifest)
	for _, d := range dependencies {
		dep := d.(map[interface{}]interface{})
		depBranch := m.DependencyBranch(dep, branch)
		if err := m.PushDependency(dep["name"].(string), chain); err != nil {
			log.PrError(err.Error())
		}

		log.PrInfo("%s depends on: %s", cargo, dep["name"])
		parent = true

		crane(dep["repo"].(string), dep["name"].(string), depBranch,
			prefix, destination, sshkey, sshpass, chain)
	}

	if parent {
		log.PrInfo("Returning to installation of %s", cargo)
	}

	files, err := fs.FileList(path.Join(clonedir, prefix))
	u.Check(err, true)

	heavyLifting(files, destination, clonedir)

	log.PrInfo("Cleaning for %s", cargo)
	fs.CleanTempDir(clonedir)
}

func heavyLifting(files []string, destination string, clonedir string) {
	// Copy /* into the destination
	for _, src := range files {
		dst := destination

		if path.Base(src) == "MANIFEST.yaml" {
			continue
		}

		// Prevent losing the first directory on directory copies.
		fileInfo, err := os.Stat(src)
		u.Check(err, true)

		re := regexp.MustCompile(clonedir + "/")

		var logmsg string
		if fileInfo.IsDir() {
			base := path.Base(src)
			dst = path.Join(destination, base)
			logmsg = fmt.Sprintf("Installing %s/ to %s/", re.ReplaceAllString(src, ""), dst)
		} else {
			logmsg = fmt.Sprintf("Installing %s to %s", re.ReplaceAllString(src, ""), dst)
		}

		log.PrVerbose(*verbose, logmsg)

		if err := fs.Install(src, dst); err != nil {
			log.PrFatal("Could not install %s into %s: %s", src, dst, err)
		}
	}
}

func parseManifest(clonedir string) map[interface{}]interface{} {
	manifestFile := path.Join(clonedir, "MANIFEST.yaml")
	if err := fs.CanReadFile(manifestFile, "MANIFEST file"); err != nil {
		log.PrError(err.Error())
	}

	manifest := m.ReadFile(manifestFile)
	if err := m.Validate(manifest); err != nil {
		log.PrError("Invalid manifest: %s", err.Error())
	}

	return manifest
}
