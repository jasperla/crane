package manifest

import (
	"path"
	"path/filepath"
)

// IgnorePatterns() takes a Manifest and returns the array of ignore patterns
func IgnorePatterns(manifest map[interface{}]interface{}) []interface{} {
	var ignore_patterns []interface{}

	if manifest["ignore"] != nil {
		ignore_patterns = manifest["ignore"].([]interface{})
	} else {
		ignore_patterns = make([]interface{}, 0)
	}

	return ignore_patterns
}

// Checks if `file` is marked as to ignore by any of the patterns
func IsIgnored(patterns []interface{}, file string) bool {
	for _, pattern := range patterns {
		// Hot path, direct match.
		if file == pattern {
			return true
		}

		// See if the parent directory is ignored.
		// NB: path.Dir() loses the trailing '/'
		if path.Dir(file) == pattern || path.Dir(file)+"/" == pattern {
			return true
		}

		// Now walk up the parents of `base` to see if any of those directories are
		// ignored. Consider hitting '/' the stop condition.
		if file == "/" {
			return false
		}
		if IsIgnored(patterns, path.Dir(file)) {
			return true
		}

		// Finally resort to globbing; simply return false in case of errors.
		if matched, err := filepath.Match(pattern.(string), file); err != nil {
			return false
		} else if matched {
			return true
		}
	}

	return false
}
