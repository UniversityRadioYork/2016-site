package controllers

import (
	"github.com/UniversityRadioYork/2016-site/models"
	"log"
	"net/http"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/2016-site/structs"
	"html/template"
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

	td := structs.Globals{
		PageContext: ic.config.PageContext,
		PageData: data,
	}

	t := template.New("base.tmpl") // Create a template.
	t, err = t.ParseFiles(
		"views/partials/header.tmpl",
		"views/partials/footer.tmpl",
		"views/elements/navbar.tmpl",
		"views/partials/base.tmpl",
		"views/index.tmpl",
	)  // Parse template file.

	if err != nil {
		log.Println(err)
		return
	}

	err = t.Execute(w, td)  // merge.

	if err != nil {
		log.Println(err)
	}

}
