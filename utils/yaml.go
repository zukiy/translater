package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// MustReadYAML read and parse file config
func MustReadYAML(path string, config interface{}) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, config)
	if err != nil {
		return err
	}

	return nil
}
