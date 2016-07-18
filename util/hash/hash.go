package hash

import (
	"fmt"

	"github.com/RedCoolBeans/crane/util/logging"
	"github.com/RedCoolBeans/crane/util/manifest"
)

func Verify(contents []interface{}, fullsrc string, src string, algo string, strict bool) bool {
	manifestHash := manifest.HashFor(contents, src, algo)
	if manifestHash == "" {
		// If we're in 'strict' mode, require a hash for each file or bail out.
		// Otherwise just return a non-match and let the caller decide what to do.
		if strict {
			logging.PrError("Could not find hash for %s", src)
		} else {
			return false
		}
	}

	if fileHash, err := FileSha256(fullsrc); err != nil {
		logging.PrError("Could not calculate hash for %s: %v", fullsrc, err)
	} else {
		return (manifestHash == fmt.Sprintf("%x", fileHash))
	}

	return false
}
