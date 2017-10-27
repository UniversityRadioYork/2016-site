package models

import "github.com/UniversityRadioYork/myradio-go"

// TeamsModel is the model for getting team data
type TeamsModel struct {
	Model
}

// NewTeamsModel returns a new TeamsModel on the MyRadio session s.
func NewTeamsModel(s *myradio.Session) *TeamsModel {
	return &TeamsModel{Model{session: s}}
}

// Get gets the data required for the GetInvolved controller from MyRadio.
//
// On success, it returns all the current teams
// Otherwise, it returns undefined data and the error causing failure.
func (m *TeamsModel) Get() (currentTeams []myradio.Team, err error) {
	currentTeams, err = m.session.GetCurrentTeams()
	return
}
