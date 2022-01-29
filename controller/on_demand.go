package controller

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	view "github.com/UniversityRadioYork/2016-site/views"
	"github.com/UniversityRadioYork/myradio-go"
)

func (c *Controller) HandleOnDemand() http.HandlerFunc {
	v := view.View{
		MainTmpl: "on_demand.tmpl",
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		PodcastsM := models.NewPodcastModel(c.Session)

		latestPodcasts, err := PodcastsM.GetAllPodcasts(10, 0)

		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}

		OnDemandM := models.NewOnDemandModel(c.Session)

		latestTimeslots, err := OnDemandM.GetLastMixcloudTimeslots()

		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}

		data := struct {
			LatestPodcasts  []myradio.Podcast
			LatestTimeslots []myradio.Timeslot
		}{
			LatestPodcasts:  latestPodcasts,
			LatestTimeslots: latestTimeslots,
		}

		err = v.Render(w, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
