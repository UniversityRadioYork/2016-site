package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	myradio "github.com/UniversityRadioYork/myradio-go"
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

// GetAllPodcasts handles the HTTP GET request r for the all postcasts page, writing to w.
func (podcastsC *PodcastController) GetAllPodcasts(w http.ResponseWriter, r *http.Request) {

	podcastm := models.NewPodcastModel(podcastsC.session)

	vars := mux.Vars(r)

	pageNumber, _ := strconv.Atoi(vars["page"])
	if pageNumber == 0 {
		pageNumber = 1
	}
	pageNumberPrev := pageNumber - 1
	pageNumberNext := pageNumber + 1

	//podcast page offset is indexed from 0, URL's are from 1.
	podcasts, err := podcastm.GetAllPodcasts(10, pageNumber-1)

	if podcasts == nil {
		utils.RenderTemplate(w, podcastsC.config.PageContext, err, "404.tmpl")
	}
	//see if it's possible to load another podcast for a possible next page.
	nextPodcasts, _ := podcastm.GetAllPodcasts(1, pageNumber)

	var pageNext = false
	if nextPodcasts != nil {
		pageNext = true
	}
	if err != nil {
		log.Println(err)
		utils.RenderTemplate(w, podcastsC.config.PageContext, err, "404.tmpl")
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
		log.Println(err)
		err = utils.RenderTemplate(w, podcastsC.config.PageContext, nil, "404.tmpl")
		return
	}

	if podcast.Status != "Published" {
		err = utils.RenderTemplate(w, podcastsC.config.PageContext, nil, "404.tmpl")
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

// GetEmbed handles the HTTP GET request r for a singular podcast embed, writing to w.
func (podcastsC *PodcastController) GetEmbed(w http.ResponseWriter, r *http.Request) {

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

	err = utils.RenderTemplate(w, podcastsC.config.PageContext, data, "podcast_player.tmpl")

	if err != nil {
		log.Println(err)
		return
	}

}
