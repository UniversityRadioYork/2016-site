package controllers

import (
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/cbroglie/mustache"
	"log"
	"net/http"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/2016-site/structs"
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

	td := struct {
		Globals structs.Globals
		Local   myradio.CurrentAndNext
	}{
		Local: *data,
		Globals: structs.Globals{
			PageContext: ic.config.PageContext,
			PageData: nil,
		},
	}

	output, err := mustache.RenderFile("views/index.mustache", td)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte(output))

}
