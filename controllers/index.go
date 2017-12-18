package controllers

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// IndexController is the controller for the index page.
type IndexController struct {
	Controller
}

// NewIndexController returns a new IndexController with the MyRadio session s
// and configuration context c.
func NewIndexController(s *myradio.Session, c *structs.Config) *IndexController {
	return &IndexController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the index page, writing to w.
func (ic *IndexController) Get(w http.ResponseWriter, r *http.Request) {
	// This is where any form params would be parsed
	model := models.NewIndexModel(ic.session)

	currentAndNext, banners, teams, podcasts, err := model.Get()

	if err != nil {
		log.Println(err)
		return
	}

	data := struct {
		CurrentAndNext *myradio.CurrentAndNext
		Banners        []myradio.Banner
		Teams          []myradio.Team
		Podcasts       []myradio.Podcast
	}{
		CurrentAndNext: currentAndNext,
		Banners:        banners,
		Teams:          teams,
		Podcasts:       podcasts,
	}

	err = utils.RenderTemplate(w, ic.config.PageContext, data, "index.tmpl", "elements/current_and_next.tmpl", "elements/banner.tmpl", "elements/message_box.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
