package models

import myradio "github.com/UniversityRadioYork/myradio-go"

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
	allpodcasts, err := m.session.GetAllPodcasts(number, page, false)
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
// It does not retrieve a podcast if the podcast is unpublished.
func (m *PodcastModel) Get(id int) (*myradio.Podcast, error) {
	pod, err := m.session.GetPodcastWithShow(id)
	if err != nil {
		return nil, err
	}

	// we should not show unpublished podcasts, regardless of situation
	// (https://github.com/UniversityRadioYork/2016-site/issues/269)
	if pod.Status != "Published" {
		return nil, nil
	}

	return pod, nil
}
