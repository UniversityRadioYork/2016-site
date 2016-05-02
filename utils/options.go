package utils

import (
	"github.com/BurntSushi/toml"
	"github.com/UniversityRadioYork/2016-site/structs"
	"io/ioutil"
	"path/filepath"
)

// GetConfigFromFile reads the website config from the given path.
//
// path is a filepath, relative to the current working directory, of a
// TOML file marshallable to a structs.Config struct.
//
// Returns a config struct and nil if the config read was successful,
// and an undefined value and non-nil otherwise.
func GetConfigFromFile(path string) (*structs.Config, error) {

	c := structs.Config{}
	cr := &c

	absPath, _ := filepath.Abs(path)

	b, err := ioutil.ReadFile(absPath)

	if err != nil {
		return cr, err
	}

	s := string(b)

	_, err = toml.Decode(s, cr)
	return cr, err
}
