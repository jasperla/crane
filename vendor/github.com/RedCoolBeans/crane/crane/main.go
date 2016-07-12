package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/RedCoolBeans/crane/util"
	"github.com/RedCoolBeans/crane/util/fs"
	g "github.com/RedCoolBeans/crane/util/git"
	"github.com/RedCoolBeans/crane/util/gpg"
	"github.com/RedCoolBeans/crane/util/hash"
	log "github.com/RedCoolBeans/crane/util/logging"
	m "github.com/RedCoolBeans/crane/util/manifest"
	"github.com/RedCoolBeans/crane/util/ssh"
	"gopkg.in/libgit2/git2go.v24"
)

var (
	verbose *bool
	debug   *bool
	strict  *bool
)

// Default hashing algorithm used for verifying files
const HASH_ALGO = "sha256"

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
	strict = flag.Bool("strict", true, "Enable strict signature and checksum checking")
	pubkey := flag.String("pubkey", "pubkey.asc", "Path to GPG public key")
	signature := flag.String("sig", "MANIFEST.yaml.sig", "Path to Manifest signature")

	flag.Parse()

	// debug implies verbose
	if *debug {
		*verbose = true
	}

	// When running crane with -clean, don't remove itself while still
	// running, unless it's supposed to be the only thing to do.
	if *clean && !gotCargo(*cargo) {
		fs.CleanSelf(*verbose)
		return
	}

	if !gotCargo(*cargo) {
		log.PrError("No package specified to load")
	}

	if *strict {
		if ok, ids := gpg.Verify(*pubkey, *signature, *verbose); ok {
			log.PrInfoBegin("Signature for MANIFEST.yaml verified\n")
			log.PrInfoEnd("Signed by: %s\n", strings.Join(ids, "\n\t"))
		} else {
			log.PrError("INVALID signature for MANIFEST.yaml! Aborting.")
		}
	}

	if err := fs.CanReadDir(*destination, "Destination directory"); err != nil {
		log.PrFatal(err.Error())
	}

	chain := m.InitDependencyChain(*cargo)

	// Everything is setup, hand-off to the main loop
	crane(*repo, *cargo, *branch, *prefix, *destination, *sshkey, *sshpass, &chain)

	fs.CleanSelf(*verbose)
}

func gotCargo(cargo string) bool {
	if len(strings.TrimSpace(cargo)) < 1 {
		return false
	} else {
		return true
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
	util.Check(err, false)
	log.PrVerbose(*verbose, "Using %s to store temporary files", clonedir)

	sshOptions := ssh.SshOptions{}
	sshOptions.Enabled = false

	if ssh.CanHandle(repo) {
		sshOptions.Enabled = true
		sshOptions.Sshkey = sshkey
		sshOptions.Sshpass = sshpass

		err = ssh.Init(&sshOptions, repo, cargo)
		util.Check(err, false)
	}

	options, cargoRepo := initGitOptions(&sshOptions, branch, repo, cargo)

	log.PrInfo("Fetching %s (%s)...", cargo, branch)
	headCommit, err := g.Clone(cargoRepo, branch, clonedir, *options)
	util.Check(err, false)

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

	heavyLifting(destination, clonedir, path.Join(clonedir, prefix))

	log.PrInfo("Cleaning for %s", cargo)
	fs.CleanTempDir(clonedir)
}

func heavyLifting(destination string, clonedir string, prefix string) {
	contents := m.Contents(parseManifest(clonedir))

	files, err := fs.FileList(path.Join(clonedir, prefix))
	util.Check(err, true)

	// Copy /* into the destination
	for _, src := range files {
		dst := destination

		if path.Base(src) == "MANIFEST.yaml" {
			continue
		}

		// Prevent losing the first directory on directory copies.
		fileInfo, err := os.Stat(src)
		util.Check(err, true)

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

		// Perform checksum verification on this file. If there's a hash recorded
		// use it. If there is not and we're in strict mode, fail.
		checksum := m.HashFor(contents, src, HASH_ALGO)
		if *strict && checksum == "" {
			log.PrError("No %s checksum found in manifest for %s", HASH_ALGO, src)
		}

		if ok := hash.Verify(contents, src, HASH_ALGO, *strict); !ok {
			log.PrError("Checksum mismatch for %s (%s)", src, HASH_ALGO)
		}

		if err := fs.Install(src, dst); err != nil {
			log.PrFatal("Could not install %s into %s: %s", src, dst, err)
		}

		if mode := m.ModeFor(contents, src); mode > 0 {
			os.Chmod(dst, os.FileMode(mode))
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
