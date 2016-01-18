package web

import (
	"github.com/UniversityRadioYork/2016-site/controllers"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	*negroni.Negroni
}

/**
 * @TODO: This is where I would pass in the config options, if I had any :(
 */
func NewServer(o map[string]string) *Server {

	s := Server{negroni.Classic()}

	router := mux.NewRouter()

	getRouter := router.Methods("GET").Subrouter()

	// Routes go in here
	nfc := controllers.NewNotFoundController()
	router.NotFoundHandler = http.HandlerFunc(nfc.Get)

	ic := controllers.NewIndexController()
	getRouter.HandleFunc("/", ic.Get)

	// End routes

	s.UseHandler(router)

	return &s

}
