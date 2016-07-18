package fs

import (
	"io"
	"os"
	"path"

	log "github.com/RedCoolBeans/crane/util/logging"
)

// Install is the single entry-point to CopyDir and CopyFile.
// Special care is needed to ensure can copy a file to an
// unnamed destination e.g. copy /file to /dest/ , we need
// to append the /file to the destination.
func Install(fullsrc string, src string, dest string, verbose bool) (err error) {
	log.PrVerbose(verbose, "fullsrc:%s, src:%s, dest:%s", fullsrc, src, dest)

	sourceInfo, err := os.Stat(fullsrc)
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		if err := os.MkdirAll(path.Join(dest, src), 0755); err != nil {
			log.PrFatal("Could not create directory %s: %s", src, err)
		}
	} else {
		if err := CopyFile(fullsrc, path.Join(dest, src)); err != nil {
			log.PrFatal("Could not install file %s: %s", src, err)
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

	return
}
