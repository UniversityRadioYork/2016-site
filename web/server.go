package web

import (
	"github.com/UniversityRadioYork/2016-site/controllers"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	*negroni.Negroni
}

func NewServer(o structs.Options) *Server {

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
