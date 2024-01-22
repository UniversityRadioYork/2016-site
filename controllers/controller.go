package controllers

/*
 @TODO: Change the default methods to render a mustache template, or log into DB??
*/

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// ControllerInterface is the interface to which controllers adhere.
type ControllerInterface interface {
	Get()     //method = GET processing
	Post()    //method = POST processing
	Delete()  //method = DELETE processing
	Put()     //method = PUT handling
	Head()    //method = HEAD processing
	Patch()   //method = PATCH treatment
	Options() //method = OPTIONS processing
}

// Controller is the base type of controllers in the 2016site architecture.
type Controller struct {
	session *myradio.Session
	config  *structs.Config
}

// Get handles a HTTP GET request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (*Controller) Get(w http.ResponseWriter, _ *http.Request) {
	notAllowed(w)
}

// Post handles a HTTP POST request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (*Controller) Post(w http.ResponseWriter, _ *http.Request) {
	notAllowed(w)
}

// Delete handles a HTTP DELETE request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (*Controller) Delete(w http.ResponseWriter, _ *http.Request) {
	notAllowed(w)
}

// Put handles a HTTP PUT request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (*Controller) Put(w http.ResponseWriter, _ *http.Request) {
	notAllowed(w)
}

// Head handles a HTTP HEAD request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (*Controller) Head(w http.ResponseWriter, _ *http.Request) {
	notAllowed(w)
}

// Patch handles a HTTP PATCH request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (*Controller) Patch(w http.ResponseWriter, _ *http.Request) {
	notAllowed(w)
}

// Options handles a HTTP OPTIONS request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (*Controller) Options(w http.ResponseWriter, _ *http.Request) {
	notAllowed(w)
}

// notAllowed sends a HTTP Method Not Allowed error.
func notAllowed(w http.ResponseWriter) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// notFound renders a 404 Not Found error originating from this controller.
// It takes an optional error that represents the failure to find something,
// and passes it into the template.
//
// This shouldn't be confused with NotFoundController, which handles 404s
// originating from the user asking for a URL that doesn't exist.
func (c *Controller) notFound(w http.ResponseWriter, err error) {
	// TODO(@MattWindsor91): what about non-404 errors?  how do we handle those?
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusNotFound)
	c.renderTemplate(w, err, "404.tmpl")
}

// renderTemplate renders a template from this Controller's page context.
func (c *Controller) renderTemplate(w http.ResponseWriter, data interface{}, mainTmpl string, addTmpls ...string) {
	if err := utils.RenderTemplate(w, c.config.PageContext, data, mainTmpl, addTmpls...); err != nil {
		// TODO(@MattWindsor91): handle error more gracefully
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
