package controllers

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// OnDemandController is the controller for the URY On demand (URY on Tap) pages.
type OnDemandController struct {
	Controller
}

// NewOnDemandController returns a new OnDemandController with the MyRadio session s
// and configuration context c.
func NewOnDemandController(s *myradio.Session, c *structs.Config) *OnDemandController {
	return &OnDemandController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the  page, writing to w.
func (onDemandC *OnDemandController) Get(w http.ResponseWriter, r *http.Request) {

	PodcastsM := models.NewPodcastModel(onDemandC.session)

	latestPodcasts, err := PodcastsM.GetAllPodcasts(10, 0)

	OnDemandM := models.NewOnDemandModel(onDemandC.session)

	latestTimeslots, err := OnDemandM.GetLastMixcloudTimeslots()

	if err != nil {
		log.Println(err)
		utils.RenderTemplate(w, onDemandC.config.PageContext, err, "404.tmpl")
		return
	}

	data := struct {
		LatestPodcasts  []myradio.Podcast
		LatestTimeslots []myradio.Timeslot
	}{
		LatestPodcasts:  latestPodcasts,
		LatestTimeslots: latestTimeslots,
	}

	err = utils.RenderTemplate(w, onDemandC.config.PageContext, data, "on_demand.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}
