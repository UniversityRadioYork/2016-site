package models

import (
	"log"

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
// On success, it returns undefined (lack of an error)
// Otherwise, it returns feedback to the user and the error causing failure.
func (m *SignUpModel) Post(formParams map[string][]string) (err error) {
	user, err := m.session.CreateNewUser(formParams)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(user.MemberID)
	log.Println(formParams["interest"])
	return
}
