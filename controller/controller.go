package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/UniversityRadioYork/2016-site/config"
	view "github.com/UniversityRadioYork/2016-site/views"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Controller struct {
	Session *myradio.Session
	Config  *config.Config

	notFoundView view.View
}

func (c *Controller) Handler() http.Handler {
	n := negroni.Classic()

	router := mux.NewRouter().StrictSlash(true)

	getRouter := router.Methods("GET").Subrouter()
	postRouter := router.Methods("POST").Subrouter()
	headRouter := router.Methods("HEAD").Subrouter()

	c.notFoundView = view.View{
		MainTmpl: "404.tmpl",
		Context:  c.Config.PageContext,
	}
	router.NotFoundHandler = c.HandleShorts()

	getRouter.HandleFunc("/", c.HandleIndex())
	postRouter.HandleFunc("/", c.HandleIndexPost())

	getRouter.HandleFunc("/search/", c.HandleSearch())

	getRouter.HandleFunc("/schedule/shows/", redirect("/schedule/thisweek/"))
	getRouter.HandleFunc("/schedule/shows/{id:[0-9]+}/", c.HandleShow()).Name("show")
	getRouter.HandleFunc("/schedule/shows/timeslots/{id:[0-9]+}/", c.HandleTimeslot()).Name("timeslot")
	getRouter.HandleFunc("/schedule/shows/seasons/{id:[0-9]+}/", c.HandleSeason).Name("season")

	getRouter.HandleFunc("/schedule/shows/{id:[0-9]+}/podcast_rss", c.HandlePodcastRss).Name("podcast_rss")
	headRouter.HandleFunc("/schedule/shows/{id:[0-9]+}/podcast_rss", c.HandlePodcastRssHead).Name("podcast_rss_head")

	getRouter.HandleFunc("/uyco/", c.HandleUYCO()).Name("uyco")

	getRouter.HandleFunc("/schedule/", redirect("/schedule/thisweek/"))

	handleScheduleWeek := c.HandleScheduleWeek(router)
	getRouter.HandleFunc("/schedule/thisweek/", handleScheduleWeek).Name("schedule-thisweek")
	getRouter.HandleFunc("/schedule/{year:[1-9][0-9][0-9][0-9]}/w{week:[0-5]?[0-9]}/", handleScheduleWeek).Name("schedule-week")
	// This route exists so that day schedule links from the previous website aren't broken.
	getRouter.HandleFunc("/schedule/{year:[1-9][0-9][0-9][0-9]}/w{week:[0-5]?[0-9]}/{day:[1-7]}/", handleScheduleWeek).Name("schedule-week-day-compat")

	getRouter.HandleFunc("/ontap/", c.HandleOnDemand())

	handlePodcasts := c.HandlePodcasts()
	getRouter.HandleFunc("/podcasts/", handlePodcasts)
	getRouter.HandleFunc("/podcasts/page/{page:[0-9]+}", handlePodcasts)
	handlePodcast := c.HandlePodcast()
	getRouter.HandleFunc("/podcasts/{id:[0-9]+}/", handlePodcast)
	getRouter.HandleFunc("/podcasts/{id:[0-9]+}/player/", handlePodcast)

	// Redirect old podcast URLs
	getRouter.HandleFunc("/uryplayer/", redirect("/ontap/"))
	getRouter.HandleFunc("/listen/", redirect("/about/"))
	getRouter.HandleFunc("/uryplayer/podcasts/", redirect("/podcasts/"))
	getRouter.HandleFunc("/uryplayer/podcasts/{id:[0-9]+}/", redirectID("/podcasts/%d/"))
	getRouter.HandleFunc("/uryplayer/podcasts/{id:[0-9]+}/player/", redirectID("/podcasts/%d/player/"))

	getRouter.HandleFunc("/people/{id:[0-9]+}/", c.HandlePeople())

	getRouter.HandleFunc("/teams/", c.HandleTeams())
	getRouter.HandleFunc("/teams/{alias}/", c.HandleTeam())

	getRouter.HandleFunc("/getinvolved/", c.HandleGetInvolved())

	postRouter.HandleFunc("/signup/", c.HandleSignupPost())

	getRouter.HandleFunc("/about/", c.HandleStatic("about.tmpl"))
	getRouter.HandleFunc("/contact/", c.HandleStatic("contact.tmpl"))
	getRouter.HandleFunc("/competitions/", c.HandleStatic("competitions.tmpl"))

	if c.Config.PageContext.CIN {
		getRouter.HandleFunc("/cin/", c.HandleStatic("cin.tmpl"))
	}

	n.UseHandler(router)
	return n
}

func redirect(to string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, http.StatusMovedPermanently)
	}
}

func redirectID(to string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		next := redirect(fmt.Sprintf(to, id))
		next(w, r)
	}
}
