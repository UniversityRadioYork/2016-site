package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/web"
	"github.com/stretchr/graceful"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Llongfile)

	config := &structs.Config{}
	_, err := toml.DecodeFile("config.toml", config)
	if err != nil {
		log.Fatal(err)
	}

	s, err := web.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	l := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	log.Printf("Listening on: %s", l)
	graceful.Run(l, time.Duration(config.Server.Timeout), s)
}
