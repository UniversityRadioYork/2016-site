package utils

import (
	"github.com/UniversityRadioYork/2016-site/structs"
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"path/filepath"
)

func GetOptionsFromFile(path string) (*structs.Options, error) {

	o := structs.Options{}
	or := &o

	absPath, _ := filepath.Abs(path)

	b, err := ioutil.ReadFile(absPath);

	if err != nil {
		return or, err;
	}

	s := string(b)

	_, err = toml.Decode(s, or);

	if err != nil {
		return or, err;
	}

	return or, nil;

}
