package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/stretchr/graceful"
)

func main() {
	log.SetFlags(log.Llongfile)

	config := &structs.Config{}
	_, err := toml.DecodeFile("config.toml", config)
	if err != nil {
		log.Fatal(err)
	}

	if &config.Schedule.StartHour != nil {
		utils.StartHour = config.Schedule.StartHour
	}

	s, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	l := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	log.Printf("Listening on: %s", l)
	graceful.Run(l, time.Duration(config.Server.Timeout), s)
}
