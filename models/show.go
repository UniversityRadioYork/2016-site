package models

import (
	"fmt"

	"github.com/UniversityRadioYork/myradio-go"
)

// ShowModel is the model for the Show controller.
type ShowModel struct {
	Model
}

// NewShowModel returns a new ShowModel on the MyRadio session s.
func NewShowModel(s *myradio.Session) *ShowModel {
	// @TODO: Pass in the config options
	return &ShowModel{Model{session: s}}
}

//
//func (m *ShowModel) Get(term string) (*[]myradio.ShowMeta, error) {
//
////	return m.session.(term);
//
//}

// GetShow gets the show with show ID id.
//
// On success, it returns the show's metadata, season list, and nil.
// Otherwise, it returns undefined data and the error causing failure.
func (m *ShowModel) GetShow(id int) (show *myradio.ShowMeta, seasons []myradio.Season, creditsToUsers map[string][]myradio.User, err error) {
	show, err = m.session.GetShow(id)
	if err != nil {
		return
	}

	creditsToUsers, err = m.session.GetCreditsToUsers(id, false)

	creditsToUsers = PluralCredits(creditsToUsers)

	if err != nil {
		fmt.Println(err)
	}

	seasons, err = m.session.GetSeasons(id)
	return
}

// GetTimeslot gets the timeslot with ID id.
//
// On success, it returns the timeslot information, the tracklist and nil.
// Otherwise, it returns undefined data and the error causing the failure.
func (m *ShowModel) GetTimeslot(id int) (timeslot myradio.Timeslot, tracklist []myradio.TracklistItem, creditsToUsers map[string][]myradio.User, err error) {
	timeslot, err = m.session.GetTimeslot(id)
	if err != nil {
		return
	}

	creditsToUsers, err = m.session.GetCreditsToUsers(id, true)

	tracklist, err = m.session.GetTrackListForTimeslot(id)
	return
}

// GetSeason gets the show season with season ID id.
//
// On success, it returns the season information, timeslots and nil.
// Otherwise, it returns undefined data and the error causing the failure.
func (m *ShowModel) GetSeason(id int) (season myradio.Season, timeslots []myradio.Timeslot, err error) {
	season, err = m.session.GetSeason(id)
	if err != nil {
		return
	}
	timeslots, err = m.session.GetTimeslotsForSeason(id)
	return
}

// PluralCredits takes a list of credits and makes roles plural if necessary
//
// On success, it returns the map with pluralised role names
func PluralCredits(dict map[string][]myradio.User) map[string][]myradio.User {
	new_dict := make(map[string][]myradio.User)
	for k, v := range dict {
		if len(v) > 1 {
			new_dict[k+"s"] = v
		} else {
			new_dict[k] = v
		}
	}
	return new_dict
}
