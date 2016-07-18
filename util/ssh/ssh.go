package ssh

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	fs "github.com/RedCoolBeans/crane/util/fs"
)

func Init(sshOptions *SshOptions, repository string, cargo string) error {
	if err := ValidKey(sshOptions.Sshkey, "Private key"); err != nil {
		return err
	}

	sshOptions.Sshpubkey = PubKey(sshOptions.Sshkey)

	if err := ValidKey(sshOptions.Sshpubkey, "Public key"); err != nil {
		return err
	}

	sshOptions.Sshrepo = repository + cargo

	sshOptions.Sshuser = FindUserName(repository)

	return nil
}

func ValidKey(path string, description string) error {
	if len(strings.TrimSpace(path)) < 1 {
		err := fmt.Sprintf("No SSH key specified")
		return errors.New(err)
	}

	if strings.HasPrefix(path, "~") || strings.Contains(path, "..") {
		err := fmt.Sprintf("Path to %s must be absolute, is %s", description, path)
		return errors.New(err)
	}

	if err := fs.CanReadFile(path, description); err != nil {
		return err
	}

	return nil
}

func PubKey(privkey string) string {
	return fmt.Sprintf("%s.pub", privkey)
}

func FindUserName(repository string) string {
	var user string

	// If the repo URI looks like it might contain a username,
	// extract it. Default to 'git' otherwise.
	if strings.Contains(repository, "@") {
		getUser := regexp.MustCompile(`(?U)(\w+)@`)
		username := getUser.FindStringSubmatch(repository)
		user = username[1]
	} else {
		user = "git"
	}

	return user
}
