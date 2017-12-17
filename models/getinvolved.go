package models

import "fmt"
import "github.com/UniversityRadioYork/myradio-go"

// GetInvolvedModel is the model for getting data for the getinvolved controller
type GetInvolvedModel struct {
	Model
}

// NewGetInvolvedModel returns a new GetInvolvedModel on the MyRadio session s.
func NewGetInvolvedModel(s *myradio.Session) *GetInvolvedModel {
	return &GetInvolvedModel{Model{session: s}}
}

// Get gets the data required for the GetInvolved controller from MyRadio.
//
// On success, it returns all the current teams, and a map from listID to
//     the team associated with that list
// Otherwise, it returns undefined data and the error causing failure.
func (m *GetInvolvedModel) Get() (numTeams int, teamInterestLists map[int]*myradio.Team, err error) {
	// Get the currently active teams
	currentTeams, err := m.session.GetCurrentTeams()
	if err != nil {
		return
	}
	numTeams = len(currentTeams)
	// Get the list of all existing mailing lists
	allMailingLists, err := m.session.GetAllLists()
	if err != nil {
		return
	}
	// Filter out the lists that aren't "interest" to save iterating them
	interestLists := allMailingLists[:0]
	for _, list := range allMailingLists {
		if len(list.Address) > 9 && list.Address[len(list.Address)-8:] == "interest" {
			interestLists = append(interestLists, list)
		}
	}
	//For each team, find the relevant interest list and add it to the map
	teamInterestLists = make(map[int]*myradio.Team)
	for k, team := range currentTeams {
		for _, list := range interestLists {
			if list.Address[:len(list.Address)-9] == team.Alias {
				teamInterestLists[list.Listid] = &currentTeams[k]
				break
			}
		}
	}
	for k, v := range teamInterestLists {
		fmt.Println(k)
		fmt.Println(*v)
	}

	return
}
