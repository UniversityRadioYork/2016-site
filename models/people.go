package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

// PeopleModel is the model for the People controller.
type PeopleModel struct {
	Model
}

// NewPeopleModel returns a new PeopleModel on the MyRadio session s.
func NewPeopleModel(s *myradio.Session) *PeopleModel {
	return &PeopleModel{Model{session: s}}
}

// Get gets the data required for the People controller from MyRadio.
//
// On success, it returns the users name, bio, a list of officerships, their photo and nil
// Otherwise, it returns undefined data and the error causing failure.
func (m *PeopleModel) Get(id int) (name, bio string, officerships []myradio.Officership, pic myradio.Photo, err error) {
	bio, err = m.session.GetUserBio(id)
	if err != nil {
		return
	}
	name, err = m.session.GetUserName(id)
	if err != nil {
		return
	}
	officerships, err = m.session.GetUserOfficerships(id)
	if err != nil {
		return
	}
	pic, err = m.session.GetUserProfilePhoto(id)
	if err != nil {
		return
	}
	return
}
