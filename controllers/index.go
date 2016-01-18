package controllers

import (
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/cbroglie/mustache"
	"log"
	"net/http"
)

type IndexController struct {
	Controller
}

func NewIndexController() *IndexController {
	return &IndexController{}
}

// @TODO: Remove this!!
type TemplateData struct {
	Global map[string]string
	Local  models.NowNextPayload
}

func (sc *IndexController) Get(w http.ResponseWriter, r *http.Request) {

	// This is where any form params would be parsed

	model := models.NewIndexModel()

	data, err := model.Get()

	if err != nil {
		log.Println(err)
		return
	}

	td := TemplateData{Local: data.Payload, Global: map[string]string{"name": "University Radio York"}}

	output, err := mustache.RenderFile("views/index.mustache", td)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte(output))

}
