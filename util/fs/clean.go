package fs

import (
	"os"

	"github.com/kardianos/osext"
	log "github.com/RedCoolBeans/crane/util/logging"
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

func CleanTempDir(temp string) {
	if err := os.RemoveAll(temp); err != nil {
		log.PrFatal("Could not remove temporary %s: %s", temp, err)
	}
}
