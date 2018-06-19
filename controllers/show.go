package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
)

// ShowController is the controller for looking up shows.
type ShowController struct {
	Controller
}

// NewShowController returns a new ShowController with the MyRadio session s
// and configuration context c.
func NewShowController(s *myradio.Session, c *structs.Config) *ShowController {
	return &ShowController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for all shows, writing to w.
func (sc *ShowController) Get(w http.ResponseWriter, r *http.Request) {

	// Do the pagination!!

	// Call the DB for the things

	// Show the things

}

// GetShow handles the HTTP GET request r for an individual show, writing to w.
func (sc *ShowController) GetShow(w http.ResponseWriter, r *http.Request) {
	sm := models.NewShowModel(sc.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	show, seasons, creditsToUsers, err := sm.GetShow(id)

	// Needed so that credits are grouped by type

	var scheduledSeasons = make([]myradio.Season, 0)
	var timeslots = make([]myradio.Timeslot, 0)

	if err != nil {
		log.Println(err)
		utils.RenderTemplate(w, sc.config.PageContext, struct{}{}, "404.tmpl")
		return
	}

	for _, season := range seasons {
		_, timeslotsSingleSeason, _ := sm.GetSeason(season.SeasonID)
		if season.FirstTimeRaw != "Not Scheduled" {
			scheduledSeasons = append(scheduledSeasons, season)
			timeslots = append(timeslots, timeslotsSingleSeason...)
		}
	}
	var latestEndTime time.Time
	var currentTime = time.Now()
	var latestTimeslot myradio.Timeslot
	var latestMixcloud bool

	for _, timeslot := range timeslots {

		layout := "02/01/2006 15:04"
		startTimeRaw, _ := time.Parse(layout, timeslot.StartTimeRaw)
		var endTimeRaw = startTimeRaw.Add(timeslot.Duration)
		if endTimeRaw.After(latestEndTime) && endTimeRaw.Before(currentTime) {
			latestEndTime = endTimeRaw
			latestTimeslot = timeslot
		}
	}
	latestMixcloud = strings.HasPrefix(latestTimeslot.MixcloudStatus, "/URY1350/")

	data := struct {
		Show           myradio.ShowMeta
		Seasons        []myradio.Season
		Timeslots      []myradio.Timeslot
		LatestTimeslot myradio.Timeslot
		LatestMixcloud bool
		CreditsToUsers map[string][]myradio.User
	}{
		Show:           *show,
		Seasons:        scheduledSeasons,
		Timeslots:      timeslots,
		LatestTimeslot: latestTimeslot,
		LatestMixcloud: latestMixcloud,
		CreditsToUsers: creditsToUsers,
	}

	if err != nil {
		log.Println(err)
		utils.RenderTemplate(w, sc.config.PageContext, data, "404.tmpl")
		return
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "show.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}

// GetTimeslot handles the HTTP GET request r for an individual timeslot, writing to w.
func (sc *ShowController) GetTimeslot(w http.ResponseWriter, r *http.Request) {
	sm := models.NewShowModel(sc.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	timeslot, tracklist, creditsToUsers, err := sm.GetTimeslot(id)
	mixcloudavailable := false

	if strings.HasPrefix(timeslot.MixcloudStatus, "/URY1350/") {
		mixcloudavailable = true
	}
	data := struct {
		Timeslot          myradio.Timeslot
		Tracklist         []myradio.TracklistItem
		MixcloudAvailable bool
		CreditsToUsers    map[string][]myradio.User
	}{
		Timeslot:          timeslot,
		Tracklist:         tracklist,
		MixcloudAvailable: mixcloudavailable,
		CreditsToUsers:    creditsToUsers,
	}

	if err != nil {
		log.Println(err)
		utils.RenderTemplate(w, sc.config.PageContext, data, "404.tmpl")
		return
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "timeslot.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}

// GetSeason handles the HTTP GET request r for an individual season, writing to w.
func (sc *ShowController) GetSeason(w http.ResponseWriter, r *http.Request) {
	sm := models.NewShowModel(sc.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	season, _, err := sm.GetSeason(id)

	if err != nil {
		utils.RenderTemplate(w, sc.config.PageContext, struct{}{}, "404.tmpl")
		log.Println(err)
		return
	}

	//We don't want a dedicated season page, redirect to the show page.
	var showURL = fmt.Sprintf("/schedule/shows/%d/?seasonID=%d", season.ShowMeta.ShowID, season.SeasonID)

	http.Redirect(w, r, showURL, 301)

}
