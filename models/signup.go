package models

import (
	"log"
	"regexp"

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
func (m *SignUpModel) Post(formParams map[string][]string) (feedback []string, err error) {
	//Validate that necessary params are present
	log.Println("Length of fname list:")
	if formParams["fname"][0] == "" {
		feedback = append(feedback, "You need to provide your first name")
	}
	if formParams["sname"][0] == "" {
		feedback = append(feedback, "You need to provide your second name")
	}
	if formParams["eduroam"][0] == "" {
		feedback = append(feedback, "You need to provide your york email")
	} else {
		match, _ := regexp.MatchString("^[a-z]{1,6}[0-9]{1,6}$", formParams["eduroam"][0])
		if !match {
			feedback = append(feedback, "The @york.ac.uk email you provided seems invalid")
		}
	}
	if formParams["phone"][0] == "" {
		delete(formParams, "phone")
	}
	if len(feedback) == 0 {
		err = m.session.CreateNewUser(formParams)
		if err != nil {
			feedback = append(feedback, "Oops. Something went wrong on our end.")
		}
	}
	return
}
