package fs

import (
	// "fmt"
	"io"
	"os"

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
		if err := os.Mkdir(dest, 0755); err != nil {
			log.PrFatal("Could not create directory %s: %s", source, err)
		}
	} else {
		if err := CopyFile(source, dest); err != nil {
			log.PrFatal("Could not install file %s: %s", source, err)
		}
	}

	return
}

// CopyFile copies a single file
func CopyFile(source string, dest string) (err error) {
	log.PrInfo("source:%s, dest:%s")
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
