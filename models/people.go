package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

type PeopleModel struct {
	Model
}

func NewPeopleModel(s *myradio.Session) *PeopleModel {
	return &PeopleModel{Model{session:s}}
}

func (m *PeopleModel) Get(id int) (name, bio string, officerships []myradio.Officership, pic myradio.Photo, err error) {
	bio, err = m.session.GetUserBio(id)
	if (err != nil) {
		return
	}
	name, err = m.session.GetUserName(id)
	if (err != nil) {
		return
	}
	officerships, err = m.session.GetUserOfficerships(id)
	if (err != nil) {
		return
	}
	pic, err = m.session.GetUserProfilePhoto(id)
	if (err != nil) {
		return
	}
	return
}
