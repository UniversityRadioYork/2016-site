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

// ShowController is the controller for looking up shows.
type ShowController struct {
	Controller
}

// NewShowController returns a new ShowController with the MyRadio session s
// and configuration context c.
func NewShowController(s *myradio.Session, c *structs.Config) *ShowController {
	return &ShowController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for all shows, writing to w.
func (sc *ShowController) Get(w http.ResponseWriter, r *http.Request) {

	// Do the pagination!!

	// Call the DB for the things

	// Show the things

}

// GetShow handles the HTTP GET request r for an individual show, writing to w.
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

	data := struct {
		Show    myradio.ShowMeta
		Seasons []myradio.Season
	}{
		Show:    *show,
		Seasons: seasons,
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "show.tmpl")
	if err != nil {
		log.Println(err)
		return
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

	data := struct {
		Timeslot myradio.Timeslot
	}{
		Timeslot: timeslot,
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "timeslot.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}
