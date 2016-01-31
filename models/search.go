package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

type SearchModel struct {
	Model
}

// @TODO: Pass in the config options
func NewSearchModel(s *myradio.Session) *SearchModel {
	return &SearchModel{Model{session:s}}
}

func (m *SearchModel) Get() (myradio.CurrentAndNext, error) {

	data, err := m.session.GetCurrentAndNext();

	return *data, err;

}
