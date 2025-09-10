package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/UniversityRadioYork/myradio-go"
)

// SignUpModel is the model for getting team data
type SignUpModel struct {
	Model
}

// NewSignUpModel returns a new SignUpModel on the MyRadio session s.
func NewSignUpModel(s *myradio.Session) *SignUpModel {
	return &SignUpModel{Model{session: s}}
}

// Get gets the data required for the GetInvolved controller from MyRadio.
//
// On success, it returns all the current teams, and a map from listID to
//
//	the team associated with that list
//
// Otherwise, it returns undefined data and the error causing failure.
func (m *SignUpModel) Get() (colleges []myradio.College, numTeams int, teamInterestLists map[int]*myradio.Team, trainings []myradio.TrainingSession, err error) {
	// Get a list of the colleges and IDs
	colleges, err = m.session.GetColleges()
	if err != nil {
		return
	}
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

	allTrainings, err := m.session.GetFutureTrainingSessions()
	if err != nil {
		return
	}

	trainings = make([]myradio.TrainingSession, 0, len(allTrainings))
	for _, training := range allTrainings {
		if training.PresenterStatusID == "Studio Trained" {
			trainings = append(trainings, training)
		}
	}

	return colleges, numTeams, teamInterestLists, trainings, nil
}

// Post posts the data from the sign up form to the api
//
// Returns an error or lack thereof based on success
func (m *SignUpModel) Post(formParams map[string][]string) (createdNewUser bool, err error) {
	user, err := m.session.CreateOrActivateUser(formParams)
	if err != nil {
		log.Println(err)
		return
	}
	if user == nil {
		return
	}
	for _, listID := range formParams["interest"] {
		LID, err := strconv.Atoi(listID)
		if err != nil {
			log.Println(err)
			continue
		}
		err = m.session.OptIn(user.MemberID, LID)
		if err != nil {
			fmt.Printf("Failed to subscribe to list %d:", LID)
			log.Println(err)
		}
	}
	createdNewUser = true
	return
}
