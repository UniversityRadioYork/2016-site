package utils

import (
	"github.com/UniversityRadioYork/2016-site/structs"
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"path/filepath"
)

func GetOptionsFromFile(path string) (options structs.Options, err error) {

	absPath, _ := filepath.Abs(path)

	b, err := ioutil.ReadFile(absPath);

	if err != nil {
		return;
	}
	
	s := string(b)

	_, err = toml.Decode(s, &options);

	return;

}