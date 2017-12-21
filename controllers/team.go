package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
)

// TeamController is the controller for the team information pages.
type TeamController struct {
	Controller
}

// NewTeamController returns a new TeamController with the MyRadio session s
// and configuration context c.
func NewTeamController(s *myradio.Session, c *structs.Config) *TeamController {
	return &TeamController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the user bio page, writing to w.
func (teamC *TeamController) Get(w http.ResponseWriter, r *http.Request) {

	teamM := models.NewTeamModel(teamC.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	currentAndNext, err := teamM.Get(id)

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	data := struct {
		CurrentAndNext *myradio.CurrentAndNext
	}{
		CurrentAndNext: currentAndNext,
	}

	err = utils.RenderTemplate(w, teamC.config.PageContext, data, "team.tmpl", "elements/current_and_next.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}

// GetAll handles the HTTP GET request r for the all teams page, writing to w.
func (teamC *TeamController) GetAll(w http.ResponseWriter, r *http.Request) {

	teamM := models.NewTeamModel(teamC.session)

	teams, err := teamM.GetAll()

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	data := struct {
		Teams []myradio.Team
	}{
		Teams: teams,
	}

	err = utils.RenderTemplate(w, teamC.config.PageContext, data, "teams.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}
