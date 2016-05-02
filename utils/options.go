package utils

import (
	"github.com/BurntSushi/toml"
	"github.com/UniversityRadioYork/2016-site/structs"
	"io/ioutil"
	"path/filepath"
)

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

	if err != nil {
		return cr, err
	}

	return cr, nil

}
