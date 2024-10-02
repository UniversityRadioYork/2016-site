package controllers

import (
	"net/http"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
)

// StaticController is the controller for the static pages.
type StaticController struct {
	Controller
}

// NewStaticController returns a new StaticController with the MyRadio
// session s and configuration context c.
func NewStaticController(c *structs.Config) *StaticController {
	return &StaticController{Controller{config: c}}
}

// GetContact handles the HTTP GET request r for the Contact page, writing to w.
func (staticC *StaticController) GetContact(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, staticC.config.PageContext, nil, "contact.tmpl")
}

// GetInvolved handles the HTTP GET request r for the Get Involved page, writing to w.
func (staticC *StaticController) GetInvolved(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, staticC.config.PageContext, nil, "getinvolved.tmpl")
}

// GetCompetitions handles the HTTP GET request r for the Get Involved page, writing to w.
func (staticC *StaticController) GetCompetitions(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, staticC.config.PageContext, nil, "competitions.tmpl")
}

func (staticC *StaticController) GetCIN(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, staticC.config.PageContext, nil, "cin.tmpl")
}
