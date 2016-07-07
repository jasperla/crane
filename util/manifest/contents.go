package manifest

const (
	DEFAULT_FILEMODE = 0644
	DEFAULT_DIRMODE  = 0755
)

// Contents() takes a Manifest and returns the array of contents. Omitted values
// (i.e. filemode) are ignored.
func Contents(manifest map[interface{}]interface{}) []interface{} {
	var contents []interface{}

	if manifest["contents"] != nil {
		contents = manifest["contents"].([]interface{})
	} else {
		contents = make([]interface{}, 0)
	}

	return contents
}

// Return the filemode for a given file. If no mode is found
// the default is returned returns.
func ModeFor(contents []interface{}, file string, isdir bool) int {
	for _, c := range contents {
		entry := c.(map[interface{}]interface{})
		if entry["path"].(string) == file {
			if mode, ok := entry["mode"].(int); ok {
				return mode
			} else {
				if isdir {
					return DEFAULT_DIRMODE
				} else {
					return DEFAULT_FILEMODE
				}
			}
		}
	}

	// file wasn't found
	return 0
}

// Returns the hash for a given file matching the algorithm.
func HashFor(contents []interface{}, file string, algo string) string {
	for _, c := range contents {
		entry := c.(map[interface{}]interface{})
		if entry["path"].(string) == file {
			if checksum, ok := entry[algo].(string); ok {
				return checksum
			} else {
				break
			}
		}
	}

	return ""
}
