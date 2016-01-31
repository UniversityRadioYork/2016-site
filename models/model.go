package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

type ModelInterface interface {
	Get() (data *interface{}, err error) // @TODO: Refactor this to a more appropriate name
}

type Model struct {
	session *myradio.Session
}
