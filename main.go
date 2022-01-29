package main

import (
	"fmt"
	"log"
	"time"

	"github.com/UniversityRadioYork/2016-site/config"
	"github.com/UniversityRadioYork/2016-site/controller"
	"github.com/stretchr/graceful"
)

func main() {
	log.SetFlags(log.Llongfile)
	config, err := config.New("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	s, err := config.Session()
	if err != nil {
		log.Fatal(err)
	}
	c := controller.Controller{
		Session: s,
		Config:  config,
	}
	l := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	log.Printf("Listening on: %s", l)
	graceful.Run(l, time.Duration(config.Server.Timeout), c.Handler())
}
