package git

import (
	"errors"
	"fmt"
	"os"

	git2go "gopkg.in/libgit2/git2go.v23"
)

func Clone(repository string, branch string, tempdir string, options git2go.CloneOptions) (*git2go.Commit, error) {
	clone, err := git2go.Clone(repository, tempdir, &options)
	if err != nil {
		e := fmt.Sprintf("Could not clone %s (%s) into %s: %s", repository, branch, tempdir, err)
		return nil, errors.New(e)
	}

	head, err := clone.Head()
	if err != nil {
		e := fmt.Sprintf("Failed to lookup HEAD: %s", err)
		return nil, errors.New(e)
	}

	headCommit, err := clone.LookupCommit(head.Target())
	if err != nil {
		e := fmt.Sprintf("Failed to lookup HEAD commit: %s", err)
		return nil, errors.New(e)
	}

	return headCommit, nil
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
