package controllers

/*
 @TODO: Change the default methods to render a mustache template, or log into DB??
*/

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/myradio-go/api"
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

// handleError converts an error into a descriptive error page. It takes care of logging it with the given context.
func (c *Controller) handleError(w http.ResponseWriter, r *http.Request, err error, context string) {
	var httpErr utils.HTTPError
	if errors.As(err, &httpErr) {
		switch httpErr.Status {
		case 404, 500: // these we have templates for
			w.WriteHeader(httpErr.Status)
			utils.RenderTemplate(w, c.config.PageContext, nil, fmt.Sprintf("%d.tmpl", httpErr.Status))
			return
		default:
			http.Error(w, httpErr.Message, httpErr.Status)
			return
		}
	}

	pc, file, line, ok := runtime.Caller(1)
	if ok {
		fn := runtime.FuncForPC(pc)
		context = fmt.Sprintf("%s at %s:%d (%s)", context, file, line, fn.Name())
	}

	var apiErr api.Error
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case http.StatusNotFound:
			w.WriteHeader(404)
			utils.RenderTemplate(w, c.config.PageContext, nil, "404.tmpl")
			return
		case http.StatusForbidden:
			// 2016-site should never hit this, it's likely a misconfiguration of our API key's permissions
			log.Printf("Received 403 from MyRadio API [%s]: %v", context, err)
		default:
			log.Printf("Unexpected MyRadio API error [%s]: %v", context, err)
		}
	} else {
		log.Printf("Unexpected error [%s]: %v", context, err)
	}

	w.WriteHeader(500)
	utils.RenderTemplate(w, c.config.PageContext, nil, "500.tmpl")
}
