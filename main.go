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
	options, err := utils.GetOptionsFromFile("./config.toml")

	s, err := web.NewServer(options)

	if err != nil {
		log.Fatal(err)
	}

	l := fmt.Sprintf("%s:%d", options.Server.Address, options.Server.Port)

	log.Printf("Listening on: %s", l)

	//@TODO: All of this in a config!
	graceful.Run(l, time.Duration(options.Server.Timeout), s)

}
