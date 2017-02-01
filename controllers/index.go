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

	currentAndNext, banners, err := model.Get()

	if err != nil {
		log.Println(err)
		return
	}
	data := struct {
		CurrentAndNext *myradio.CurrentAndNext
		Banners        []myradio.Banner
		Sched1         []myradio.Show
		Sched2         []myradio.Show
		Sched3         []myradio.Show
		Sched4         []myradio.Show
	}{
		CurrentAndNext: currentAndNext,
		Banners:        banners,
		Sched1:         currentAndNext.Next[:2],
		Sched2:         currentAndNext.Next[2:4],
		Sched3:         currentAndNext.Next[4:8],
		Sched4:         currentAndNext.Next[8:12],
	}

	err = utils.RenderTemplate(w, ic.config.PageContext, data, "index.tmpl", "elements/current_and_next.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
