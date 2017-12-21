package models

import "github.com/UniversityRadioYork/myradio-go"

// TeamModel is the model for the Team controller.
type TeamModel struct {
	Model
}

// NewTeamModel returns a new TeamModel on the MyRadio session s.
func NewTeamModel(s *myradio.Session) *TeamModel {
	return &TeamModel{Model{session: s}}
}

// Get gets the data required for the Team controller from MyRadio.
//
// On success, it returns the users name, bio, a list of officerships, their photo if they have one and nil
// Otherwise, it returns undefined data and the error causing failure.
func (m *TeamModel) Get(id int) (currentAndNext *myradio.CurrentAndNext, err error) {
	currentAndNext, err = m.session.GetCurrentAndNext()
	if err != nil {
		return
	}
	return
}

// GetAll gets the data required for the Team controller from MyRadio.
//
// On success, it returns the team information and nil
// Otherwise, it returns undefined data and the error causing failure.
func (m *TeamModel) GetAll() (currentAndNext *myradio.CurrentAndNext, err error) {
	currentAndNext, err = m.session.GetCurrentAndNext()
	if err != nil {
		return
	}
	return
}
