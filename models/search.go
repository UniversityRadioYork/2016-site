package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

// SearchModel is the model for the Search controller.
type SearchModel struct {
	Model
}

// NewSearchModel returns a new SearchModel on the MyRadio session s.
func NewSearchModel(s *myradio.Session) *SearchModel {
	// @TODO: Pass in the config options
	return &SearchModel{Model{session: s}}
}

// Get gets the data required for the Search controller from MyRadio.
//
// term is the string term to search for.  This is currently a show
// search.
//
// On success, it returns the search results, and nil.
// Otherwise, it returns undefined data and the error causing failure.
func (m *SearchModel) Get(term string) ([]myradio.ShowMeta, error) {

	return m.session.GetSearchMeta(term)

}
