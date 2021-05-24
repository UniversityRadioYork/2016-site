package models

import (
	myradio "github.com/UniversityRadioYork/myradio-go"
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
func (m *IndexModel) Get() (currentAndNext *myradio.CurrentAndNext, banners []myradio.Banner, teams []myradio.Team, podcasts []myradio.Podcast, showOnAir bool, err error) {
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

	//Get 10 podcasts from page 0 (the latest podcasts)
	allpodcasts := make([]myradio.Podcast, 10)
	allpodcasts, err = m.session.GetAllPodcasts(10, 0, false)
	if err != nil {
		return
	}

	for i, p := range allpodcasts {
		if p.Status == "Published" {
			podcasts = append(podcasts, allpodcasts[i])
		}
	}

	selectorInfo, err := m.session.GetSelectorInfo()
	if err != nil {
		return
	}
	showOnAir = !(selectorInfo.Studio == myradio.SelectorJukebox || selectorInfo.Studio == myradio.SelectorOffAir)

	return
}
