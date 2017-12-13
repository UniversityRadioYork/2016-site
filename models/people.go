package models

import "github.com/UniversityRadioYork/myradio-go"

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
// On success, it returns the users name, bio, a list of officerships, their photo if they have one and nil
// Otherwise, it returns undefined data and the error causing failure.
func (m *PeopleModel) Get(id int) (user *myradio.User, officerships []myradio.Officership, credits []myradio.ShowMeta, currentAndNext *myradio.CurrentAndNext, err error) {
	user, err = m.session.GetUser(id)
	if err != nil {
		return
	}

	officerships, err = m.session.GetUserOfficerships(id)
	if err != nil {
		return
	}

	credits, err = m.session.GetUserShowCredits(id)
	if err != nil {
		return
	}

	currentAndNext, err = m.session.GetCurrentAndNext()
	if err != nil {
		return
	}

	return
}
