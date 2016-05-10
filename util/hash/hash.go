package hash

import (
	"fmt"

	"github.com/RedCoolBeans/crane/util/logging"
	"github.com/RedCoolBeans/crane/util/manifest"
)

func Verify(contents []interface{}, filepath string, algo string) bool {
	manifestHash := manifest.HashFor(contents, filepath, algo)
	if manifestHash == "" {
		logging.PrError("Could not find hash for %s", filepath)
	}

	if fileHash, err := FileSha256(filepath); err != nil {
		logging.PrError("Could not calculate hash for %s: %v", filepath, err)
	} else {
		return (manifestHash == fmt.Sprintf("%x", fileHash))
	}

	return false
}
