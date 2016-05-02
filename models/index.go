package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

type IndexModel struct {
	Model
}

// @TODO: Pass in the config options
func NewIndexModel(s *myradio.Session) *IndexModel {
	return &IndexModel{Model{session: s}}
}

func (m *IndexModel) Get() (*myradio.CurrentAndNext, error) {

	return m.session.GetCurrentAndNext()

}
