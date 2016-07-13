package fs

import (
	"os"

	log "github.com/RedCoolBeans/crane/util/logging"
	"github.com/kardianos/osext"
)

func CleanSelf(verbose bool) {
	log.PrVerbose(verbose, "Removing self")
	myname, err := osext.Executable()
	if err != nil {
		log.PrFatal("Could not get own executable name")
	}

	if err := os.Remove(myname); err != nil {
		log.PrFatal("Could not remove %s", myname)
	}
}

func CleanAll(path string, verbose bool) {
	log.PrVerbose(verbose, "Removing %s", path)

	if err := os.RemoveAll(path); err != nil {
		log.PrFatal("Could not remove %s", path)
	}
}

func CleanTempDir(temp string) {
	if err := os.RemoveAll(temp); err != nil {
		log.PrFatal("Could not remove temporary %s: %s", temp, err)
	}
}
