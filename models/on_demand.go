package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

// OnDemandModel is the model for the OnDemand controller.
type OnDemandModel struct {
	Model
}

// NewOnDemandModel returns a new OnDemandModel on the MyRadio session s.
func NewOnDemandModel(s *myradio.Session) *OnDemandModel {
	// @TODO: Pass in the config options
	return &OnDemandModel{Model{session: s}}
}

// GetLastMixcloudTimeslots gets the data required for the OnDemand controller from MyRadio.
//
// On success, it returns the previous 6 timeslots, and nil.
// Otherwise, it returns undefined data and the error causing failure.
func (m *OnDemandModel) GetLastMixcloudTimeslots() (timeslots []myradio.Timeslot, err error) {
	timeslots, err = m.session.GetPreviousTimeslots(6)
	if err != nil {
		return
	}

	return
}
