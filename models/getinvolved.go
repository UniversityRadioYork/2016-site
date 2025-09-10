package models

import (
	"github.com/BurntSushi/toml"
	"github.com/UniversityRadioYork/myradio-go"
)

// GetInvolvedModel is the model for getting data for the getinvolved controller
type GetInvolvedModel struct {
	Model
}

// FAQ contains all the FAQ objects, containing a question and an answer
type FAQ struct {
	FAQs []struct {
		Question string `toml:"question"`
		Answer   string `toml:"answer"`
		SeeMore  string `toml:"seemore"`
	} `toml:"faqs"`
}

// NewGetInvolvedModel returns a new GetInvolvedModel on the MyRadio session s.
func NewGetInvolvedModel(s *myradio.Session) *GetInvolvedModel {
	return &GetInvolvedModel{Model{session: s}}
}

func (m *GetInvolvedModel) Get() (faq *FAQ, err error) {
	// Decodes the FAQ toml file into faq
	_, err = toml.DecodeFile("faqs.toml", &faq)

	return
}
