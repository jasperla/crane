package git

import (
	"errors"
	"fmt"
	"os"

	git2go "gopkg.in/libgit2/git2go.v24"
)

func Clone(repository string, branch string, tempdir string, options git2go.CloneOptions) error {
	_, err := git2go.Clone(repository, tempdir, &options)
	if err != nil {
		e := fmt.Sprintf("Could not clone %s (%s) into %s: %s\n    Are you using a password protected SSH key without -sshpass?", repository, branch, tempdir, err)
		return errors.New(e)
	}

	return nil
}

func RemoveDotGit(tempdir string) error {
	// XXX: Find the right glob to exclude .git
	tempGit := fmt.Sprintf("%s/.git", tempdir)

	if err := os.RemoveAll(tempGit); err != nil {
		e := fmt.Sprintf("Could not remove %s: %s", tempGit, err)
		return errors.New(e)
	}

	return nil
}
