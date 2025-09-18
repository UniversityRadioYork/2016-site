package models

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/UniversityRadioYork/myradio-go"
)

// SignUpModel is the model for getting team data
type SignUpModel struct {
	Model
}

type SignUpTrainingSession struct {
	myradio.TrainingSessionForSignup
}

func (s *SignUpTrainingSession) RemainingSpaces() int {
	return s.MaxParticipants - s.AttendeeCount
}

func (s *SignUpTrainingSession) TooLateToSignUp() bool {
	if s.SignupCutoffHours == 0 {
		return false
	}
	now := time.Now()
	remainingTime := s.StartTime().Sub(now)
	return remainingTime.Hours() < float64(s.SignupCutoffHours)
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
func (m *SignUpModel) Get() (colleges []myradio.College, numTeams int, teamInterestLists map[int]*myradio.Team, trainings []SignUpTrainingSession, err error) {
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

	allTrainings, err := m.session.GetFutureTrainingSessionsForSignup()
	if err != nil {
		return
	}

	trainings = make([]SignUpTrainingSession, 0, len(allTrainings))
	for _, training := range allTrainings {
		training := SignUpTrainingSession{
			training,
		}
		if training.PresenterStatusID == "Studio Trained" &&
			training.AttendeeCount != training.MaxParticipants &&
			!training.TooLateToSignUp() {
			trainings = append(trainings, training)
		}
	}

	return colleges, numTeams, teamInterestLists, trainings, nil
}

// Post posts the data from the sign up form to the api
//
// Returns an error or lack thereof based on success
func (m *SignUpModel) Post(formParams map[string][]string) (createdNewUser bool, trainingSignupResult int, err error) {
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
	trainingSessionID, ok := formParams["sessionid"]
	if ok {
		if trainingSessionID[0] == "!!unavailable" {
			_, err = m.session.AddToWaitingList(1, user.MemberID)
			if err != nil {
				return
			}

			trainingSignupResult = -5
		} else if trainingSessionID[0] != "" {
			var demoID int
			demoID, err = strconv.Atoi(trainingSessionID[0])
			if err != nil {
				log.Println(err)
				return
			}
			trainingSignupResult, err = m.session.AddAttendeeToDemo(demoID, user.MemberID)
			if err != nil {
				return
			}
		}
	}
	createdNewUser = true
	return
}
