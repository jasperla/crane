package manifest

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
