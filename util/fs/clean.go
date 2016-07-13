package fs

import (
	"os"

	log "github.com/RedCoolBeans/crane/util/logging"
	"github.com/kardianos/osext"
)

func CleanSelf(path string, verbose bool) {
	myname, err := osext.Executable()
	if err != nil {
		log.PrFatal("Could not get own executable name")
	}

	log.PrInfo("Removing self: %s", myname)
	if err := os.Remove(myname); err != nil {
		log.PrFatal("Could not remove %s", myname)
	}

	log.PrInfo("Removing: %s", path)
	if err := os.RemoveAll(path); err != nil {
		log.PrFatal("Could not remove %s", path)
	}
}

func CleanTempDir(temp string) {
	if err := os.RemoveAll(temp); err != nil {
		log.PrFatal("Could not remove temporary %s: %s", temp, err)
	}
}
