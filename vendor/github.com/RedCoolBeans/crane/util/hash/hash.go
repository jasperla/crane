package hash

import (
	"fmt"

	"github.com/RedCoolBeans/crane/util/logging"
	"github.com/RedCoolBeans/crane/util/manifest"
)

func Verify(contents []interface{}, filepath string, algo string, strict bool) bool {
	manifestHash := manifest.HashFor(contents, filepath, algo)
	if manifestHash == "" {
		// If we're in 'strict' mode, require a hash for each file or bail out.
		// Otherwise just return a non-match and let the caller decide what to do.
		if strict {
			logging.PrError("Could not find hash for %s", filepath)
		} else {
			return false
		}
	}

	if fileHash, err := FileSha256(filepath); err != nil {
		logging.PrError("Could not calculate hash for %s: %v", filepath, err)
	} else {
		return (manifestHash == fmt.Sprintf("%x", fileHash))
	}

	return false
}
