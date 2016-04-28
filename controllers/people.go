package controllers

import (
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/2016-site/structs"
	"log"
	"net/http"
	"github.com/UniversityRadioYork/2016-site/models"
	"html/template"
	"github.com/gorilla/mux"
	"strconv"
)

type PeopleController struct {
	Controller
}

func NewPeopleController(s *myradio.Session, c *structs.Config) *PeopleController {
	return &PeopleController{Controller{session:s, config:c}}
}

func (pc *PeopleController) Get(w http.ResponseWriter, r *http.Request) {

	pm := models.NewPeopleModel(pc.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	name, bio, officerships, pic, err := pm.Get(id)

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	// Render Template
	td := structs.Globals{
		PageContext: pc.config.PageContext,
		PageData: struct {
			Name         string
			Bio          string
			Officerships []myradio.Officership
			ProfilePicture myradio.Photo
		}{
			Name: name,
			Bio: bio,
			Officerships:officerships,
			ProfilePicture: pic,
		},
	}

	t := template.New("base.tmpl") // Create a template.
	t, err = t.ParseFiles(
		"views/partials/header.tmpl",
		"views/partials/footer.tmpl",
		"views/elements/navbar.tmpl",
		"views/partials/base.tmpl",
		"views/people.tmpl",
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
