package web

import (
	"net/http"

	"github.com/UniversityRadioYork/2016-site/controllers"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Server struct {
	*negroni.Negroni
}

func NewServer(c *structs.Config) (*Server, error) {

	s := Server{negroni.Classic()}

	session, err := myradio.NewSessionFromKeyFile()

	if err != nil {
		return &s, err
	}

	router := mux.NewRouter().StrictSlash(true)

	getRouter := router.Methods("GET").Subrouter()

	// Routes go in here
	nfc := controllers.NewNotFoundController(c)
	router.NotFoundHandler = http.HandlerFunc(nfc.Get)

	ic := controllers.NewIndexController(session, c)
	getRouter.HandleFunc("/", ic.Get)

	sc := controllers.NewSearchController(session, c)
	getRouter.HandleFunc("/search/", sc.Get)

	showC := controllers.NewShowController(session, c)
	//	getRouter.HandleFunc("/schedule/shows", showC.Get) // @TODO: Implement this
	getRouter.HandleFunc("/schedule/shows/{id:[0-9]+}/", showC.GetShow)
	getRouter.HandleFunc("/schedule/shows/timeslots/{id:[0-9]+}/", showC.GetTimeslot)
	getRouter.HandleFunc("/schedule/shows/seasons/{id:[0-9]+}/", showC.GetSeason)

	pc := controllers.NewPeopleController(session, c)
	getRouter.HandleFunc("/people/{id:[0-9]+}/", pc.Get)

	teamC := controllers.NewTeamController(session, c)
	getRouter.HandleFunc("/team/{id:[0-9]+}/", teamC.Get)

	// End routes

	s.UseHandler(router)

	return &s, nil

}
