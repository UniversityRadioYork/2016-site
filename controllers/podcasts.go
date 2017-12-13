package controllers

import (
	"log"
	"net/http"
	"strconv"

	"GitHub.com/gorilla/mux"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// PodcastController is the controller for the URYPlayer Podcast pages.
type PodcastController struct {
	Controller
}

// NewPodcastController returns a new PodcastController with the MyRadio session s
// and configuration context c.
func NewPodcastController(s *myradio.Session, c *structs.Config) *PodcastController {
	return &PodcastController{Controller{session: s, config: c}}
}

// GetPodcasts handles the HTTP GET request r for the all postcasts page, writing to w.
func (podcastsC *PodcastController) GetPodcasts(w http.ResponseWriter, r *http.Request) {

	podcastm := models.NewPodcastModel(podcastsC.session)

	podcasts, err := podcastm.Get()

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	data := struct {
		Podcasts []myradio.Podcast
	}{
		Podcasts: podcasts,
	}

	err = utils.RenderTemplate(w, podcastsC.config.PageContext, data, "podcasts.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}

// Get handles the HTTP GET request r for a singular podcast page, writing to w.
func (podcastsC *PodcastController) Get(w http.ResponseWriter, r *http.Request) {

	podcastm := models.NewPodcastModel(podcastsC.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	podcast, err := podcastm.Get(id)

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	data := struct {
		Podcast *myradio.Podcast
	}{
		Podcast: podcast,
	}

	err = utils.RenderTemplate(w, podcastsC.config.PageContext, data, "podcast.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}
