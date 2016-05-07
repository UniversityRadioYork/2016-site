package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

// ShowModel is the model for the Show controller.
type ShowModel struct {
	Model
}

// NewShowModel returns a new ShowModel on the MyRadio session s.
func NewShowModel(s *myradio.Session) *ShowModel {
	// @TODO: Pass in the config options
	return &ShowModel{Model{session: s}}
}

//
//func (m *ShowModel) Get(term string) (*[]myradio.ShowMeta, error) {
//
////	return m.session.(term);
//
//}

// GetShow gets the show with show ID id.
//
// On success, it returns the show's metadata, season list, and nil.
// Otherwise, it returns undefined data and the error causing failure.
func (m *ShowModel) GetShow(id int) (*myradio.ShowMeta, []myradio.Season, error) {
	show, err := m.session.GetShow(id)
	if err != nil {
		return nil, nil, err
	}

	seasons, err := m.session.GetSeasons(id)

	return show, seasons, err
}

// GetTimeslot gets the timeslot with ID id.
//
// On success, it returns the timeslot information, the tracklist and nil.
// Otherwise, it returns undefined data and the error causing the failure.
func (m *ShowModel) GetTimeslot(id int) (timeslot myradio.Timeslot,tracklist []myradio.TracklistItem, err error) {
	timeslot, err = m.session.GetTimeslot(id)
	if err != nil {
		return
	}
	tracklist, err = m.session.GetTrackListForTimeslot(id)
	return
}
