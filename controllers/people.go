package controllers

import (
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// PeopleController is the controller for the user bio page.
type PeopleController struct {
	Controller
}

// NewPeopleController returns a new PeopleController with the MyRadio session s
// and configuration context c.
func NewPeopleController(s *myradio.Session, c *structs.Config) *PeopleController {
	return &PeopleController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the user bio page, writing to w.
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

	data := struct {
		Name           string
		Bio            string
		Officerships   []myradio.Officership
		ProfilePicture myradio.Photo
	}{
		Name:           name,
		Bio:            bio,
		Officerships:   officerships,
		ProfilePicture: pic,
	}

	err = utils.RenderTemplate(w, pc.config.PageContext, data, "people.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}
