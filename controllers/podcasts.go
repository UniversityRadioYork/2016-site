package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/myradio-go"
)

// PodcastController is the controller for the URYPlayer Podcast pages.
type PodcastController struct {
	Controller
	model *models.PodcastModel
}

// NewPodcastController returns a new PodcastController with the MyRadio session s
// and configuration context c.
func NewPodcastController(s *myradio.Session, c *structs.Config) *PodcastController {
	return &PodcastController{
		Controller: Controller{session: s, config: c},
		model:      models.NewPodcastModel(s),
	}
}

// GetAllPodcasts handles the HTTP GET request r for the all postcasts page, writing to w.
func (podcastsC *PodcastController) GetAllPodcasts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pageNumber, _ := strconv.Atoi(vars["page"])
	if pageNumber == 0 {
		pageNumber = 1
	}
	pageNumberPrev := pageNumber - 1
	pageNumberNext := pageNumber + 1

	//podcast page offset is indexed from 0, URL's are from 1.
	podcasts, err := podcastsC.model.GetAllPodcasts(10, pageNumber-1)
	if podcasts == nil {
		podcastsC.notFound(w, err)
		return
	}

	//see if it's possible to load another podcast for a possible next page.
	nextPodcasts, _ := podcastsC.model.GetAllPodcasts(1, pageNumber)

	var pageNext = false
	if nextPodcasts != nil {
		pageNext = true
	}
	if err != nil {
		podcastsC.notFound(w, err)
		return
	}

	data := struct {
		PageNumberPrev int
		PageNumber     int
		PageNumberNext int
		PageNext       bool
		Podcasts       []myradio.Podcast
	}{
		PageNumberPrev: pageNumberPrev,
		PageNumber:     pageNumber,
		PageNumberNext: pageNumberNext,
		PageNext:       pageNext,
		Podcasts:       podcasts,
	}

	podcastsC.renderTemplate(w, data, "podcasts.tmpl")
}

// Get handles the HTTP GET request r for a singular podcast page, writing to w.
func (podcastsC *PodcastController) Get(w http.ResponseWriter, r *http.Request) {
	podcast, err := podcastsC.getPodcast(r)
	if podcast == nil || err != nil {
		// TODO(@MattWindsor91): what if the error is not 404?
		podcastsC.notFound(w, err)
		return
	}

	podcastsC.renderPodcast(w, podcast, "podcast.tmpl")
}

// GetEmbed handles the HTTP GET request r for a singular podcast embed, writing to w.
func (podcastsC *PodcastController) GetEmbed(w http.ResponseWriter, r *http.Request) {
	podcast, err := podcastsC.getPodcast(r)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	// No error, but podcast is not available
	if podcast == nil {
		http.NotFound(w, r)
		return
	}

	podcastsC.renderPodcast(w, podcast, "podcast_player.tmpl")
}

func (podcastsC *PodcastController) getPodcast(r *http.Request) (*myradio.Podcast, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return nil, err
	}

	return podcastsC.model.Get(id)
}

func (podcastsC *PodcastController) renderPodcast(w http.ResponseWriter, podcast *myradio.Podcast, tmpl string) {
	data := struct {
		Podcast *myradio.Podcast
	}{
		Podcast: podcast,
	}
	podcastsC.renderTemplate(w, data, tmpl)
}
