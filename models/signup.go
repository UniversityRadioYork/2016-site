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

// Post posts the data from the sign up form to the api
//
// Returns an error or lack thereof based on success
func (m *SignUpModel) Post(formParams map[string][]string) (err error) {
	user, err := m.session.CreateNewUser(formParams)
	if err != nil {
		log.Println(err)
		return
	}
	for _, listID := range formParams["interest"] {
		LID, err := strconv.Atoi(listID)
		if err != nil {
			fmt.Printf("Failed to subscribe to list %d:", LID)
			fmt.Println(err)
			continue
		}
		err = m.session.OptIn(user.MemberID, LID)
		if err != nil {
			fmt.Printf("Failed to subscribe to list %d:", LID)
			fmt.Println(err)
		}
	}
	return
}
