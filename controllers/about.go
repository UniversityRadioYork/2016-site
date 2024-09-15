package controllers

import (
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// AboutController is the controller for the about page.
type AboutController struct {
	Controller
}

// NewAboutController returns a new AboutController with the MyRadio session s
// and configuration context c.
func NewAboutController(s *myradio.Session, c *structs.Config) *AboutController {
	return &AboutController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the about us page, writing to w.
func (aboutC *AboutController) Get(w http.ResponseWriter, r *http.Request) {
	teamM := models.NewTeamModel(aboutC.session)
	teams, err := teamM.GetAll()
	if err != nil {
		aboutC.handleError(w, r, err, "TeamModel.GetAll")
		return
	}
	data := struct {
		Teams []myradio.Team
	}{
		Teams: teams,
	}
	utils.RenderTemplate(w, aboutC.config.PageContext, data, "about.tmpl")
}
