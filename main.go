package main

import (
	"fmt"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/2016-site/web"
	"github.com/stretchr/graceful"
	"log"
	"time"
)

func main() {

	//Get the config from the config.yaml file
	config, err := utils.GetConfigFromFile("./config.toml")

	s, err := web.NewServer(config)

	if err != nil {
		log.Fatal(err)
	}

	l := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)

	log.Printf("Listening on: %s", l)

	graceful.Run(l, time.Duration(config.Server.Timeout), s)

}
