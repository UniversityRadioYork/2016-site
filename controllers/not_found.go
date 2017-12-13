package controllers

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
)

// NotFoundController is the controller for the 404 error page.
type NotFoundController struct {
	Controller
}

// NewNotFoundController returns a new NotFoundController with the MyRadio
// session s and configuration context c.
func NewNotFoundController(c *structs.Config) *NotFoundController {
	return &NotFoundController{Controller{config: c}}
}

// Get handles the HTTP GET request r for the 404 page, writing to w.
func (sc *NotFoundController) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	err := utils.RenderTemplate(w, sc.config.PageContext, nil, "404.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
