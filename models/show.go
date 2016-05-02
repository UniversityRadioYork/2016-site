package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

type ShowModel struct {
	Model
}

// @TODO: Pass in the config options
func NewShowModel(s *myradio.Session) *ShowModel {
	return &ShowModel{Model{session: s}}
}

//
//func (m *ShowModel) Get(term string) (*[]myradio.ShowMeta, error) {
//
////	return m.session.(term);
//
//}

func (m *ShowModel) GetShow(id int) (*myradio.ShowMeta, []myradio.Season, error) {

	show, err := m.session.GetShow(id)

	if err != nil {
		return nil, nil, err
	}

	seasons, err := m.session.GetSeasons(id)

	return show, seasons, err

}
