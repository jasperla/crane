package manifest

import (
	"errors"
	"fmt"
)

var RequiredFields = [...]string{"name", "maintainer", "version", "architecture"}
var RequiredDepFields = [...]string{"name", "repo"}

func Validate(manifest map[interface{}]interface{}) error {
	if err := ValidateRequiredFields(manifest); err != nil {
		return err
	}

	if err := ValidateDependencyFields(manifest); err != nil {
		return err
	}

	return nil
}

func ValidateRequiredFields(manifest map[interface{}]interface{}) error {
	for _, value := range RequiredFields {
		if _, ok := manifest[value]; !ok {
			err := fmt.Sprintf("required field %q not found", value)
			return errors.New(err)
		}
	}

	// While not required, revision must be > 0
	if manifest["revision"] != nil && manifest["revision"] == "0" {
		err := fmt.Sprintf("field revision must be > 1")
		return errors.New(err)
	}

	return nil
}

func ValidateDependencyFields(manifest map[interface{}]interface{}) error {
	var err string

	dependencies := Dependencies(manifest)
	for i, d := range dependencies {
		dep := d.(map[interface{}]interface{})
		for _, value := range RequiredDepFields {
			if _, ok := dep[value]; !ok {
				if value == "name" {
					err = fmt.Sprintf("required field %q not found for dependency #%d", value, i+1)
				} else {
					err = fmt.Sprintf("required field %q not found for dependency %q",
						value, dep["name"])
				}

				return errors.New(err)
			}
		}
	}
	return nil
}
