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

type ShowInfo struct {
	Show           *myradio.ShowMeta
	Seasons        []myradio.Season
	CreditsToUsers map[string][]myradio.User
	Podcasts       []myradio.Podcast
}

// GetShow gets the show with show ID id.
//
// On success, it returns the show's metadata, season list, and nil.
// Otherwise, it returns undefined data and the error causing failure.
func (m *ShowModel) GetShow(id int) (show *ShowInfo, err error) {
	myrShow, err := m.session.GetShow(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	creditsToUsers, err := m.session.GetCreditsToUsers(id, false)

	if err != nil {
		fmt.Println(err)
		return
	}

	seasons, err := m.session.GetSeasons(id)

	if err != nil {
		fmt.Println(err)
		return
	}

	podcasts, err := m.session.GetAllShowPodcasts(id)

	if err != nil {
		fmt.Println(err)
		return
	}

	return &ShowInfo{
		Show:           myrShow,
		Seasons:        seasons,
		CreditsToUsers: creditsToUsers,
		Podcasts:       podcasts,
	}, nil
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

	if err != nil {
		return
	}

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

func (m *ShowModel) GetPodcastRSS(id int) (string, error) {
	return m.session.GetPodcastRSS(id)
}
