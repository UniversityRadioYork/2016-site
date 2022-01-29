package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/UniversityRadioYork/2016-site/models"
	view "github.com/UniversityRadioYork/2016-site/views"
	myradio "github.com/UniversityRadioYork/myradio-go"
)

func (c *Controller) HandlePodcasts() http.HandlerFunc {
	v := view.View{
		MainTmpl: "podcasts.tmpl",
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		podcastm := models.NewPodcastModel(c.Session)

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
			c.HandleNotFound(w, r)
		}
		//see if it's possible to load another podcast for a possible next page.
		nextPodcasts, _ := podcastm.GetAllPodcasts(1, pageNumber)

		var pageNext = false
		if nextPodcasts != nil {
			pageNext = true
		}
		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
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

		err = v.Render(w, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (c *Controller) HandlePodcast() http.HandlerFunc {
	v := view.View{
		MainTmpl: "podcast.tmpl",
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {

		podcastm := models.NewPodcastModel(c.Session)

		vars := mux.Vars(r)

		id, _ := strconv.Atoi(vars["id"])

		podcast, err := podcastm.Get(id)

		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}

		if podcast.Status != "Published" {
			c.HandleNotFound(w, r)
			return
		}

		data := struct {
			Podcast *myradio.Podcast
		}{
			Podcast: podcast,
		}

		err = v.Render(w, data)

		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (c *Controller) HandlePodcastEmbed() http.HandlerFunc {
	v := view.View{
		MainTmpl: "podcast_player.tmpl",
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {

		podcastm := models.NewPodcastModel(c.Session)

		vars := mux.Vars(r)

		id, _ := strconv.Atoi(vars["id"])

		podcast, err := podcastm.Get(id)

		if err != nil {
			c.HandleNotFound(w, r)
			return
		}

		data := struct {
			Podcast *myradio.Podcast
		}{
			Podcast: podcast,
		}

		err = v.Render(w, data)

		if err != nil {
			log.Println(err)
			return
		}
	}
}
