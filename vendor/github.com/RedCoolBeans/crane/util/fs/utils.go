package fs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Check if we can read the given 'path' denoting a 'what'
func CanReadFile(path string, what string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err := fmt.Sprintf("%s %s does not exist", what, path)
			return errors.New(err)
		} else {
			err := fmt.Sprintf("Could not open %s for reading: %s", path, err)
			return errors.New(err)
		}
	}

	return nil
}

func CanReadDir(path string, what string) error {
	if destDir, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) || !destDir.IsDir() {
			err := fmt.Sprintf("%s %s does not exist", what, path)
			return errors.New(err)
		} else {
			err := fmt.Sprintf("Could not open %s for reading: %s", path, err)
			return errors.New(err)
		}
	}

	return nil
}

// CreateTempDir is essentially a wrapper around ioutil.TempDir.
func CreateTempDir() (string, error) {
	temp, err := ioutil.TempDir("/tmp", "crane-")
	if err != nil {
		e := fmt.Sprintf("Could not create temporary directory: %s", err)
		return "", errors.New(e)
	}

	return temp, nil
}

func FileList(path string) ([]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*", path))
	if err != nil {
		e := fmt.Sprintf("Could not get filelist: %s", err)
		return nil, errors.New(e)
	}

	if len(files) == 0 {
		e := fmt.Sprintf("No files found in %s", path)
		return nil, errors.New(e)
	}

	return files, nil
}
