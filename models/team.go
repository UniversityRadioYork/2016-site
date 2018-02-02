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
func (m *TeamModel) Get(alias string) (team myradio.Team, heads []myradio.Officer, assistants []myradio.Officer, officers []myradio.Officer, err error) {
	team, err = m.session.GetTeamWithOfficers(alias)
	if err != nil {
		return
	}
	var teamID int
	teamID = int(team.TeamID)

	heads, err = m.session.GetTeamHeadPositions(teamID, nil)
	if err != nil {
		return
	}
	assistants, err = m.session.GetTeamAssistantHeadPositions(teamID, nil)
	if err != nil {
		return
	}
	officers, err = m.session.GetTeamOfficerPositions(teamID, nil)
	if err != nil {
		return
	}
	return
}

// GetAll gets the data required for the Team controller from MyRadio.
//
// On success, it returns the team information and nil
// Otherwise, it returns undefined data and the error causing failure.
func (m *TeamModel) GetAll() (teams []myradio.Team, err error) {

	teams, err = m.session.GetCurrentTeams()
	if err != nil {
		return
	}

	return
}
