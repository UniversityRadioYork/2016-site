package controllers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
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
	alias := vars["alias"]
	team, heads, assistants, officers, err := teamM.Get(alias)
	if err != nil {
		teamC.handleError(w, r, err, "TeamModel.Get")
		return
	}

	data := struct {
		Team       myradio.Team
		Heads      []myradio.Officer
		Assistants []myradio.Officer
		Officers   []myradio.Officer
	}{
		Team:       team,
		Heads:      heads,
		Assistants: assistants,
		Officers:   officers,
	}
	utils.RenderTemplate(w, teamC.config.PageContext, data, "team.tmpl")
}
