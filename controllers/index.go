package controllers

import (
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"log"
	"net/http"
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

	data, err := model.Get()
	if err != nil {
		log.Println(err)
		return
	}

	err = utils.RenderTemplate(w, ic.config.PageContext, data, "index.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
