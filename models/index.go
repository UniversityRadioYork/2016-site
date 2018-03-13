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
func (m *IndexModel) Get() (currentAndNext *myradio.CurrentAndNext, banners []myradio.Banner, teams []myradio.Team, jukebox bool, err error) {
	currentAndNext, err = m.session.GetCurrentAndNext()
	if err != nil {
		return
	}
	banners, err = m.session.GetLiveBanners()
	if err != nil {
		return
	}

	teams, err = m.session.GetCurrentTeams()
	if err != nil {
		return
	}

	selectorInfo, err := m.session.GetSelectorInfo()
	if err != nil {
		return
	}
	jukebox = selectorInfo.Studio == myradio.SelectorJukebox

	return
}

// MessageModel is the model that takes the response from sending a message
type MessageModel struct {
	Model
}

// NewMessageModel returns a new MessageModel
func NewMessageModel(s *myradio.Session) *MessageModel {
	// @TODO: Pass in the config options
	return &MessageModel{Model{session: s}}
}

// Put PUTs the given message to the current show
//
// On success, it returns nil
// Otherwise, it returns the error causing failure.
func (m *MessageModel) Put(msg string) (err error) {
	currentTimeslot, err := m.session.GetCurrentTimeslot()
	if err != nil {
		return
	}
	err = m.session.PutMessage(currentTimeslot.TimeslotID, msg)
	return
}
