package main

import (
	"github.com/UniversityRadioYork/2016-site/web"
	"github.com/stretchr/graceful"
)

func main() {

	// Get the options from config.json

	//@TODO: Pass in the options!
	s := web.NewServer(map[string]string{})

	//@TODO: All of this in a config!
	graceful.Run(":3000", 0, s)

}
