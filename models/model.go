package models

import (
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/myradio-go"
)

type ModelInterface interface {
	Get() (data structs.Response, err error) // @TODO: Refactor this to a more appropriate name
}

type Model struct {
	session *myradio.Session
}
