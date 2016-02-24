package structs

import (
	"github.com/UniversityRadioYork/myradio-go"
)


type Options struct {

	Server  Server    `toml:"server"`
	Globals Globals    `toml:"globals"`

}

type Server struct {

	Address string    `toml:"address"`
	Port    int        `toml:"port"`
	Timeout int        `toml:"timout"`

}

type Globals struct {

	LongName  		string    `toml:"longName"`
	ShortName 		string    `toml:"shortName"`
	CurrentAndNext 	myradio.CurrentAndNext

}
