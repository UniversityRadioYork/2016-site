package controllers

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// GetInvolvedController is the controller for the get involved page.
type GetInvolvedController struct {
	Controller
}

// NewGetInvolvedController returns a new GetInvolvedController with the MyRadio
// session s and configuration context c.
func NewGetInvolvedController(s *myradio.Session, c *structs.Config) *GetInvolvedController {
	return &GetInvolvedController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the get involved, writing to w.
func (gic *GetInvolvedController) Get(w http.ResponseWriter, r *http.Request) {

	tm := models.NewTeamsModel(gic.session)

	teams, err := tm.Get()

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	numTeams := len(teams)

	data := struct {
		Teams    []myradio.Team
		NumTeams int
	}{
		Teams:    teams,
		NumTeams: numTeams,
	}

	err = utils.RenderTemplate(w, gic.config.PageContext, data, "getinvolved.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}
