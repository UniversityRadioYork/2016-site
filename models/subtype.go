package models

import "github.com/UniversityRadioYork/myradio-go"

// SubtypeModel is the model for the Subtype Controller
type SubtypeModel struct {
	Model
}

// NewSubtypeModel returns a new SubtypeModel
func NewSubtypeModel(s *myradio.Session) *SubtypeModel {
	return &SubtypeModel{Model{session: s}}
}

// Get returns a myradio ShowSeasonSubtype given an alias for the class
func (m *SubtypeModel) Get(alias string) (subtype myradio.ShowSeasonSubtype, err error) {
	subtype, err = m.session.GetShowSubtype(alias)
	if err != nil {
		return
	}
	return
}

// GetAll returns all myradio ShowSeasonSubtypes
func (m *SubtypeModel) GetAll() (subtypes []myradio.ShowSeasonSubtype, err error) {
	subtypes, err = m.session.GetAllShowSubtypes()
	if err != nil {
		return
	}

	return
}
