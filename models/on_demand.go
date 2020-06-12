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
// On success, it returns the previous 10/11 timeslots, and nil.
// Otherwise, it returns undefined data and the error causing failure.
func (m *OnDemandModel) GetLastMixcloudTimeslots() (timeslots []myradio.Timeslot, err error) {
	timeslots, err = m.session.GetPreviousTimeslots(11)
	currrentAndNext, err := m.session.GetCurrentAndNext()
	if err != nil {
		return
	}
	// If show currently on air, remove it from previous timeslots
	if currrentAndNext.Current.Id != 0 {
		timeslots = timeslots[1:11]
	}
	return
}
