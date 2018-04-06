package models

import "github.com/UniversityRadioYork/myradio-go"

// PodcastModel is the model for the Podcast controller.
type PodcastModel struct {
	Model
}

// NewPodcastModel returns a new PodcastModel on the MyRadio session s.
func NewPodcastModel(s *myradio.Session) *PodcastModel {
	return &PodcastModel{Model{session: s}}
}

// GetAllPodcasts gets the data required for the Podcast controller from MyRadio.
//
// On success, it returns the podcasts and nil
// Otherwise, it returns undefined data and the error causing failure.
func (m *PodcastModel) GetAllPodcasts(number int, page int) (podcasts []myradio.Podcast, err error) {
	//Get 'number' of podcasts from page 'page' (0 = the latest podcasts)
	allpodcasts := make([]myradio.Podcast, number)
	allpodcasts, err = m.session.GetAllPodcasts(number, page)
	if err != nil {
		return
	}

	for i, p := range allpodcasts {
		if p.Status == "Published" {
			podcasts = append(podcasts, allpodcasts[i])
		}
	}

	return
}

// Get gets the data required for the Podcast controller from MyRadio.
//
// On success, it returns the users name, bio, a list of officerships, their photo if they have one and nil
// Otherwise, it returns undefined data and the error causing failure.
func (m *PodcastModel) Get(id int) (podcast *myradio.Podcast, err error) {
	podcast, err = m.session.Get(id)
	if err != nil {
		return
	}

	return
}
