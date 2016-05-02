package controllers

import (
	"github.com/UniversityRadioYork/2016-site/models"
	"log"
	"net/http"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
)

type IndexController struct {
	Controller
}

func NewIndexController(s *myradio.Session, c *structs.Config) *IndexController {
	return &IndexController{Controller{session:s, config: c}}
}

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
