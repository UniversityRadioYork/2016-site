package controllers

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
)

// StaticController is the controller for the 404 error page.
type StaticController struct {
	Controller
}

// NewStaticController returns a new StaticController with the MyRadio
// session s and configuration context c.
func NewStaticController(c *structs.Config) *StaticController {
	return &StaticController{Controller{config: c}}
}

// Get handles the HTTP GET request r for the 404 page, writing to w.
func (staticC *StaticController) Get(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, staticC.config.PageContext, nil, "about.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
