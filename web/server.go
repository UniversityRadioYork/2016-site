package web

import (
	"github.com/UniversityRadioYork/2016-site/controllers"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Server struct {
	*negroni.Negroni
}

func NewServer() *Server {

	s := Server{negroni.Classic()}

	router := mux.NewRouter()

	// Routes go in here

	notFoundController := controllers.NewNotFoundController(router)
	notFoundController.Register()

	indexController := controllers.NewIndexController(router)
	indexController.Register("/")

	// End routes

	s.UseHandler(router)

	return &s

}
