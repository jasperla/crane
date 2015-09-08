package manifest

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	u "github.com/RedCoolBeans/crane/util/utils"
)

func ReadFile(file string) map[interface{}]interface{} {
	data, err := ioutil.ReadFile(file)
	u.Check(err, true)

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(data), &m)
	u.Check(err, true)

	return m
}
