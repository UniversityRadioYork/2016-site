package controllers

import (
	"log"
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

// GetAbout handles the HTTP GET request r for the About page, writing to w.
func (staticC *StaticController) GetAbout(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, staticC.config.PageContext, nil, "about.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}

// GetContact handles the HTTP GET request r for the Contact page, writing to w.
func (staticC *StaticController) GetContact(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, staticC.config.PageContext, nil, "contact.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}

// GetInvolved handles the HTTP GET request r for the Get Involved page, writing to w.
func (staticC *StaticController) GetInvolved(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, staticC.config.PageContext, nil, "getinvolved.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}

// GetCompetitions handles the HTTP GET request r for the Get Involved page, writing to w.
func (staticC *StaticController) GetCompetitions(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, staticC.config.PageContext, nil, "competitions.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}

func (staticC *StaticController) GetMusic(w http.ResponseWriter, r *http.Request) {
  err := utils.RenderTemplate(w, staticC.config.PageContext, nil, "music.tmpl")
  if err != nil {
    log.Println(err)
    return
  }
}

func (staticC *StaticController) GetCIN(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, staticC.config.PageContext, nil, "cin.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
