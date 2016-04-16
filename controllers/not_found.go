package controllers

import (
	"log"
	"net/http"
	"github.com/UniversityRadioYork/2016-site/structs"
	"html/template"
)

type NotFoundController struct {
	Controller
}

func NewNotFoundController(c *structs.Config) *NotFoundController {
	return &NotFoundController{Controller{config: c}}
}

func (sc *NotFoundController) Get(w http.ResponseWriter, r *http.Request) {

	td := structs.Globals{
		PageContext: sc.config.PageContext,
		PageData: nil,
	}

	w.WriteHeader(404)

	t := template.New("404.tmpl") // Create a template.

	t, err := t.ParseFiles(
		"views/404.tmpl",
		"views/partials/footer.tmpl",
		"views/partials/header.tmpl",
		"views/elements/navbar.tmpl",
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
