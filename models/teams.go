package models

import "github.com/UniversityRadioYork/myradio-go"

// GetInvolvedModel is the model for the GetInvolvedcontroller.
type TeamsModel struct {
	Model
}

// NewPeopleModel returns a new PeopleModel on the MyRadio session s.
func NewTeamsModel(s *myradio.Session) *TeamsModel {
	return &TeamsModel{Model{session: s}}
}

// Get gets the data required for the GetInvolved controller from MyRadio.
//
// On success, it returns all the current teams
// Otherwise, it returns undefined data and the error causing failure.
func (m *TeamsModel) Get() (teams []myradio.Team, err error) {
	currentTeams, err = m.session.GetCurrentTeams(teams)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
