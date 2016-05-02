package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

// IndexModel is the model for the Index controller.
type IndexModel struct {
	Model
}

// NewIndexModel returns a new IndexModel on the MyRadio session s.
func NewIndexModel(s *myradio.Session) *IndexModel {
	// @TODO: Pass in the config options
	return &IndexModel{Model{session: s}}
}

// Get gets the data required for the Index controller from MyRadio.
//
// On success, it returns the current and next show, and nil.
// Otherwise, it returns undefined data and the error causing failure.
func (m *IndexModel) Get() (*myradio.CurrentAndNext, error) {

	return m.session.GetCurrentAndNext()

}
