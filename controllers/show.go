package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
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

	show, seasons, creditsToUsers, err := sm.GetShow(id)

	// Needed so that credits are grouped by type

	data := struct {
		Show           myradio.ShowMeta
		Seasons        []myradio.Season
		CreditsToUsers map[string][]myradio.User
	}{
		Show:           *show,
		Seasons:        seasons,
		CreditsToUsers: creditsToUsers,
	}

	if err != nil {
		log.Println(err)
		utils.RenderTemplate(w, sc.config.PageContext, data, "404.tmpl")
		return
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

	timeslot, tracklist, err := sm.GetTimeslot(id)
	mixcloudavailable := false

	if strings.HasPrefix(timeslot.MixcloudStatus, "/URY1350/") {
		mixcloudavailable = true
	}
	data := struct {
		Timeslot          myradio.Timeslot
		Tracklist         []myradio.TracklistItem
		MixcloudAvailable bool
	}{
		Timeslot:          timeslot,
		Tracklist:         tracklist,
		MixcloudAvailable: mixcloudavailable,
	}

	if err != nil {
		log.Println(err)
		utils.RenderTemplate(w, sc.config.PageContext, data, "404.tmpl")
		return
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "timeslot.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}

func (sc *ShowController) GetSeason(w http.ResponseWriter, r *http.Request) {
	sm := models.NewShowModel(sc.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	season, timeslots, err := sm.GetSeason(id)

	data := struct {
		Season    myradio.Season
		Timeslots []myradio.Timeslot
	}{
		Season:    season,
		Timeslots: timeslots,
	}

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "season.tmpl")
	if err != nil {
		log.Println(err)

		return
	}

}
