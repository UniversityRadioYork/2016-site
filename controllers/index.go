package controllers

import (
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/cbroglie/mustache"
	"log"
	"net/http"
	"github.com/UniversityRadioYork/myradio-go"
)

type IndexController struct {
	Controller
}

func NewIndexController(s *myradio.Session) *IndexController {
	return &IndexController{Controller{session:s}}
}

type TemplateData struct {
	Global map[string]string // @TODO: Remove this!! PASS IN CONFIGS
	Local  myradio.CurrentAndNext
}

func (sc *IndexController) Get(w http.ResponseWriter, r *http.Request) {

	// This is where any form params would be parsed

	model := models.NewIndexModel(sc.session)

	data, err := model.Get()

	if err != nil {
		log.Println(err)
		return
	}

	td := TemplateData{Local: data, Global: map[string]string{"name": "University Radio York"}}

	output, err := mustache.RenderFile("views/index.mustache", td)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte(output))

}
