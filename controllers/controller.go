package controllers

/*
 @TODO: Change the default methods to render a mustache template, or log into DB??
 */

import (
	"net/http"
	"github.com/UniversityRadioYork/myradio-go"
)

type ControllerInterface interface {
	Get()                        //method = GET processing
	Post()                       //method = POST processing
	Delete()                     //method = DELETE processing
	Put()                        //method = PUT handling
	Head()                       //method = HEAD processing
	Patch()                      //method = PATCH treatment
	Options()                    //method = OPTIONS processing
}

type Controller struct {
	session *myradio.Session
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

func (c *Controller) Put(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

func (c *Controller) Head(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

func (c *Controller) Patch(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}

func (c *Controller) Options(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
}
