package utils

import "github.com/RedCoolBeans/crane/util/logging"

// Check abstracts away the horrid overuse of "if err != nil {}"
func Check(err error, fatal bool) {
	if err != nil {
		if fatal {
			logging.PrFatal(err.Error())
		} else {
			logging.PrError(err.Error())
		}
	}
}

