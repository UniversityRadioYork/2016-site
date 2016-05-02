package controllers

import (
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/2016-site/structs"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/UniversityRadioYork/2016-site/models"
	"log"
	"strconv"
	"html/template"
)

type ShowController struct {
	Controller
}

func NewShowController(s *myradio.Session, c *structs.Config) *ShowController {
	return &ShowController{Controller{session:s, config:c}}
}

func (sc *ShowController) Get(w http.ResponseWriter, r *http.Request) {

	// Do the pagination!!

	// Call the DB for the things

	// Show the things

}

func (sc *ShowController) GetShow(w http.ResponseWriter, r *http.Request) {

	sm := models.NewShowModel(sc.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	show, seasons, err := sm.GetShow(id)

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	// Render Template
	td := structs.Globals{
		PageContext: sc.config.PageContext,
		PageData: struct {
			Show    myradio.ShowMeta
			Seasons []myradio.Season
		}{
			Show: *show,
			Seasons: seasons,
		},
	}

	t := template.New("base.tmpl") // Create a template.
	t, err = t.ParseFiles(
		"views/partials/header.tmpl",
		"views/partials/footer.tmpl",
		"views/elements/navbar.tmpl",
		"views/partials/base.tmpl",
		"views/show.tmpl",
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

func (sc *ShowController) GetTimeslot(w http.ResponseWriter, r *http.Request) {

	sm := models.NewShowModel(sc.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	timeslot, err := sm.GetTimeslot(id)

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	// Render Template
	td := structs.Globals{
		PageContext: sc.config.PageContext,
		PageData: struct {
			Timeslot myradio.Timeslot
		}{
			Timeslot: timeslot,
		},
	}

	t := template.New("base.tmpl") // Create a template.
	t, err = t.ParseFiles(
		"views/partials/header.tmpl",
		"views/partials/footer.tmpl",
		"views/elements/navbar.tmpl",
		"views/partials/base.tmpl",
		"views/timeslot.tmpl",
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
