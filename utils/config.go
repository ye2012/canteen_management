package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

var (
	errNilConfig = errors.New("Config object is empty. ")
)

// LoadJSONFile gets your config from the json file,
// and fills your struct with the option
func LoadJSONFile(filename string, config interface{}) error {
	if config == nil {
		return errNilConfig
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, config)
}
