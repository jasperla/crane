package fs

import (
	"fmt"
	"io"
	"os"
	"path"

	log "github.com/RedCoolBeans/crane/util/logging"
)

// Install is the single entry-point to CopyDir and CopyFile.
// Special care is needed to ensure can copy a file to an
// unnamed destination e.g. copy /file to /dest/ , we need
// to append the /file to the destination.
func Install(source string, dest string) (err error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		if err := CopyDir(source, dest); err != nil {
			log.PrFatal("Could not install directory %s: %s", source, err)
		}
	} else {
		fileName := path.Base(source)
		dest = fmt.Sprintf("%s/%s", dest, fileName)
		if err := CopyFile(source, dest); err != nil {
			log.PrFatal("Could not install file %s: %s", source, err)
		}
	}

	return
}

// CopyDir recursively copies the source directory into a destination.
func CopyDir(source string, dest string) (err error) {
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Create target directory
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := fmt.Sprintf("%s/%s", source, obj.Name())
		destinationfilepointer := fmt.Sprintf("%s/%s", dest, obj.Name())

		if obj.IsDir() {
			// recursively create subdirs
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				log.PrFatal("%s", err)
			}
		} else {
			// perform actual copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				log.PrFatal("%s", err)
			}
		}
	}

	return
}

// CopyFile copies a single file
func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	if _, err = io.Copy(destfile, sourcefile); err != nil {
		return err
	}

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if err = os.Chmod(dest, sourceinfo.Mode()); err != nil {
		return err
	}

	return
}
