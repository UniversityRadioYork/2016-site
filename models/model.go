package models

import "github.com/UniversityRadioYork/2016-site/structs"

type ModelInteface interface {
	Get() (data structs.Response, err error) // Maps to a GET request to the API
}

type Model struct {
	config map[string]interface{}
}
