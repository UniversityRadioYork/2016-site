package main

import (
	"net/http"

	"github.com/UniversityRadioYork/2016-site/controllers"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// Server is the type of the main 2016site web application.
type Server struct {
	*negroni.Negroni
}

// NewServer creates a 2016site server based on the config c.
func NewServer(c *structs.Config) (*Server, error) {

	s := Server{negroni.Classic()}

	session, err := myradio.NewSessionFromKeyFile()

	if err != nil {
		return &s, err
	}

	router := mux.NewRouter().StrictSlash(true)

	getRouter := router.Methods("GET").Subrouter()
	postRouter := router.Methods("POST").Subrouter()

	// Routes go in here
	nfc := controllers.NewNotFoundController(c)
	router.NotFoundHandler = http.HandlerFunc(nfc.Get)

	ic := controllers.NewIndexController(session, c)
	getRouter.HandleFunc("/", ic.Get)

	sc := controllers.NewSearchController(session, c)
	getRouter.HandleFunc("/search/", sc.Get)

	showC := controllers.NewShowController(session, c)
	//	getRouter.HandleFunc("/schedule/shows", showC.Get) // @TODO: Implement this
	getRouter.HandleFunc("/schedule/shows/{id:[0-9]+}/", showC.GetShow).Name("show")
	getRouter.HandleFunc("/schedule/shows/timeslots/{id:[0-9]+}/", showC.GetTimeslot).Name("timeslot")
	getRouter.HandleFunc("/schedule/shows/seasons/{id:[0-9]+}/", showC.GetSeason).Name("season")

	// NOTE: NewScheduleWeekController assumes 'timeslot' is installed BEFORE it is called.
	schedWeekC := controllers.NewScheduleWeekController(session, getRouter, c)
	getRouter.HandleFunc("/schedule/thisweek/", schedWeekC.GetThisWeek).Name("schedule-thisweek")

	getRouter.HandleFunc("/schedule/{year:[1-9][0-9][0-9][0-9]}/w{week:[0-5]?[0-9]}/", schedWeekC.GetByYearWeek).Name("schedule-week")
	// This route exists so that day schedule links from the previous website aren't broken.
	getRouter.HandleFunc("/schedule/{year:[1-9][0-9][0-9][0-9]}/w{week:[0-5]?[0-9]}/{day:[1-7]}/", schedWeekC.GetByYearWeek).Name("schedule-week-day-compat")

	pc := controllers.NewPeopleController(session, c)
	getRouter.HandleFunc("/people/{id:[0-9]+}/", pc.Get)

	teamC := controllers.NewTeamController(session, c)
	getRouter.HandleFunc("/teams/", teamC.GetAll)
	getRouter.HandleFunc("/teams/{id:[0-9]+}/", teamC.Get)

	getinvolvedC := controllers.NewGetInvolvedController(session, c)
	getRouter.HandleFunc("/getinvolved/", getinvolvedC.Get)

	signupC := controllers.NewSignUpController(session, c)
	postRouter.HandleFunc("/signup/", signupC.Post)

	staticC := controllers.NewStaticController(c)
	getRouter.HandleFunc("/about/", staticC.GetAbout)
	getRouter.HandleFunc("/contact/", staticC.GetContact)
	getRouter.HandleFunc("/competitions/", staticC.GetCompetitions)

	// End routes

	s.UseHandler(router)

	return &s, nil

}
