package main

import (
	"github.com/UniversityRadioYork/2016-site/web"
	"github.com/stretchr/graceful"
)

func main() {

	s := web.NewServer()

	graceful.Run(":3000", 0, s)

}
