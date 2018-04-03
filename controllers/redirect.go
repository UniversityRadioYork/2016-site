package controllers

import (
	"net/http"

	"github.com/UniversityRadioYork/2016-site/structs"
)

// RedirectController is the controller for redirecting URLs to new ones.
type RedirectController struct {
	Controller
}

// NewRedirectController returns a new RedirectController with the MyRadio
// session s and configuration context c.
func NewRedirectController(c *structs.Config) *RedirectController {
	return &RedirectController{Controller{config: c}}
}

// Redirect handles the HTTP GET request r for the page to redirect, writing to w.
func (sc *RedirectController) Redirect(w http.ResponseWriter, r *http.Request, location string, code int) {
	http.Redirect(w, r, location, code)
}
