package controllers

/*
 @TODO: Change the default methods to render a mustache template, or log into DB??
*/

import (
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/myradio-go"
	"net/http"
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
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Post handles a HTTP POST request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Delete handles a HTTP DELETE request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Put handles a HTTP PUT request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Put(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Head handles a HTTP HEAD request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Head(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Patch handles a HTTP PATCH request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Patch(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Options handles a HTTP OPTIONS request r, writing to w.
//
// Unless overridden, controllers refuse this method.
func (c *Controller) Options(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}
