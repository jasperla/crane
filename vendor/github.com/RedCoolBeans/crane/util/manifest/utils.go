package manifest

import "fmt"

func VersionString(manifest map[interface{}]interface{}) string {
	var version string

	if manifest["revision"] != nil {
		version = fmt.Sprintf("%s rev. %s", manifest["version"], manifest["revision"])
	} else {
		version = manifest["version"].(string)
	}

	return version
}
