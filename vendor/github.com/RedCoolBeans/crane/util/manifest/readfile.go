package manifest

import (
	"io/ioutil"

	"github.com/RedCoolBeans/crane/util"
	"gopkg.in/yaml.v2"
)

func ReadFile(file string) map[interface{}]interface{} {
	data, err := ioutil.ReadFile(file)
	util.Check(err, true)

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(data), &m)
	util.Check(err, true)

	return m
}
