package main

import (
	"fmt"
	"net/http"
	"strconv"

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

	var session *myradio.Session
	var err error
	if c.Server.MyRadioAPI == "" {
		session, err = myradio.NewSessionFromKeyFile()
	} else {
		session, err = myradio.NewSessionFromKeyFileForServer(c.Server.MyRadioAPI)
	}

	if err != nil {
		return &s, err
	}

	router := mux.NewRouter().StrictSlash(true)

	getRouter := router.Methods("GET").Subrouter()
	postRouter := router.Methods("POST").Subrouter()
	headRouter := router.Methods("HEAD").Subrouter()

	// Routes go in here
	nfc := controllers.NewNotFoundController(c)
	router.NotFoundHandler = http.HandlerFunc(nfc.Get)

	ic := controllers.NewIndexController(session, c)
	getRouter.HandleFunc("/", ic.Get)
	postRouter.HandleFunc("/", ic.Post)

	sc := controllers.NewSearchController(session, c)
	getRouter.HandleFunc("/search/", sc.Get)

	showC := controllers.NewShowController(session, c)
	getRouter.HandleFunc("/schedule/shows/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/schedule/thisweek/", 301)
	})
	getRouter.HandleFunc("/schedule/shows/{id:[0-9]+}/", showC.GetShow).Name("show")
	getRouter.HandleFunc("/schedule/shows/timeslots/{id:[0-9]+}/", showC.GetTimeslot).Name("timeslot")
	getRouter.HandleFunc("/schedule/shows/seasons/{id:[0-9]+}/", showC.GetSeason).Name("season")

	getRouter.HandleFunc("/schedule/shows/{id:[0-9]+}/podcast_rss", showC.GetPodcastRss).Name("podcast_rss")
	headRouter.HandleFunc("/schedule/shows/{id:[0-9]+}/podcast_rss", showC.GetPodcastRssHead).Name("podcast_rss_head")

	getRouter.HandleFunc("/uyco/", showC.GetUyco).Name("uyco")

	getRouter.HandleFunc("/schedule/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/schedule/thisweek/", 301)
	})
	// NOTE: NewScheduleWeekController assumes 'timeslot' is installed BEFORE it is called.
	schedWeekC := controllers.NewScheduleWeekController(session, getRouter, c)
	getRouter.HandleFunc("/schedule/thisweek/", schedWeekC.GetThisWeek).Name("schedule-thisweek")

	getRouter.HandleFunc("/schedule/{year:[1-9][0-9][0-9][0-9]}/w{week:[0-5]?[0-9]}/", schedWeekC.GetByYearWeek).Name("schedule-week")
	// This route exists so that day schedule links from the previous website aren't broken.
	getRouter.HandleFunc("/schedule/{year:[1-9][0-9][0-9][0-9]}/w{week:[0-5]?[0-9]}/{day:[1-7]}/", schedWeekC.GetByYearWeek).Name("schedule-week-day-compat")

	onDemandC := controllers.NewOnDemandController(session, c)
	getRouter.HandleFunc("/ontap/", onDemandC.Get)

	podcastsC := controllers.NewPodcastController(session, c)
	getRouter.HandleFunc("/podcasts/", podcastsC.GetAllPodcasts)
	getRouter.HandleFunc("/podcasts/page/{page:[0-9]+}", podcastsC.GetAllPodcasts)
	getRouter.HandleFunc("/podcasts/{id:[0-9]+}/", podcastsC.Get)
	getRouter.HandleFunc("/podcasts/{id:[0-9]+}/player/", podcastsC.GetEmbed)
	// Redirect old podcast URLs
	getRouter.HandleFunc("/uryplayer/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ontap/", 301)
	})
	getRouter.HandleFunc("/listen/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/about/", 301)
	})
	getRouter.HandleFunc("/uryplayer/podcasts/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/podcasts/", 301)
	})
	getRouter.HandleFunc("/uryplayer/podcasts/{id:[0-9]+}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		http.Redirect(w, r, fmt.Sprintf("/podcasts/%d/", id), 301)
	})
	getRouter.HandleFunc("/uryplayer/podcasts/{id:[0-9]+}/player/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		http.Redirect(w, r, fmt.Sprintf("/podcasts/%d/player/", id), 301)
	})

	pc := controllers.NewPeopleController(session, c)
	getRouter.HandleFunc("/people/{id:[0-9]+}/", pc.Get)

	teamC := controllers.NewTeamController(session, c)
	getRouter.HandleFunc("/teams/", teamC.GetAll)
	getRouter.HandleFunc("/teams/{alias}/", teamC.Get)

	getinvolvedC := controllers.NewGetInvolvedController(session, c)
	getRouter.HandleFunc("/getinvolved/", getinvolvedC.Get)

	signupC := controllers.NewSignUpController(session, c)
	postRouter.HandleFunc("/signup/", signupC.Post)

	staticC := controllers.NewStaticController(c)
	getRouter.HandleFunc("/about/", staticC.GetAbout)
	getRouter.HandleFunc("/contact/", staticC.GetContact)
	getRouter.HandleFunc("/competitions/", staticC.GetCompetitions)
	if c.PageContext.CIN {
		getRouter.HandleFunc("/cin/", staticC.GetCIN)
	}
	// End routes

	s.UseHandler(router)

	return &s, nil

}
