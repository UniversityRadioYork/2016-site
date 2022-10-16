package controllers

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// weekFromVars extracts the year, and week strings from vars.
func weekFromVars(vars map[string]string) (string, string, error) {
	y, ok := vars["year"]
	if !ok {
		return "", "", utils.NewHTTPError(http.StatusBadRequest, "no year provided")
	}
	w, ok := vars["week"]
	if !ok {
		return "", "", utils.NewHTTPError(http.StatusBadRequest, "no week provided")
	}

	return y, w, nil
}

//
// Controller
//

// ScheduleWeekController is the controller for looking up week schedules.
type ScheduleWeekController struct {
	Controller

	timeslotURLBuilder     func(*myradio.Timeslot) (*url.URL, error)
	weekScheduleURLBuilder func(*time.Time) (*url.URL, error)
}

// NewScheduleWeekController returns a new ScheduleWeekController with the MyRadio session s,
// router r, and configuration context c.
// It assumes r already has routes installed for timeslots.
func NewScheduleWeekController(s *myradio.Session, r *mux.Router, c *structs.Config) *ScheduleWeekController {
	/* We pass in the router so we can generate URL reversal functions.
	   However, at the time we get the router, it hasn't set up the week
	   schedule routes yet, so we make those URL functions look up the relevant
	   routes whenever they're called.

	   TODO(MattWindsor91):
	       this is probably slow, but I didn't want to optimise prematurely. */
	wbuilder := func(t *time.Time) (*url.URL, error) {
		// TODO(MattWindsor91): use URL instead of string?

		wroute := r.Get("schedule-week")
		year, week := t.ISOWeek()

		// The router can't handle years outside this range by design.
		// Don't try to reverse them!
		if year < 1000 || 9999 < year {
			return nil, nil
		}

		return wroute.URLPath(
			"year", strconv.Itoa(year),
			"week", strconv.Itoa(week))
	}

	/* The router should have timeslot routes installed already, so we can
	   get the route eagerly. */
	troute := r.Get("timeslot")
	tbuilder := func(t *myradio.Timeslot) (*url.URL, error) {
		return troute.URLPath("id", strconv.FormatUint(t.TimeslotID, 10))
	}

	return &ScheduleWeekController{
		Controller:             Controller{session: s, config: c},
		timeslotURLBuilder:     tbuilder,
		weekScheduleURLBuilder: wbuilder,
	}
}

// GetThisWeek handles the HTTP GET request r for this week's week schedule, writing to w.
//
// It takes no request variables.
func (sc *ScheduleWeekController) GetThisWeek(w http.ResponseWriter, r *http.Request) {
	/* Today's ISO week is the same as Monday's ISO week, so we need not do
	   anything fancy like working out the Monday of this week.
	   This seems obvious, but some bits of Go disagree with us on the first
	   day of the week, so caution is always a good thing. */
	today := time.Now()
	year, week := today.ISOWeek()

	sc.makeAndRenderWeek(w, r, year, week)
}

// GetByYearWeek handles the HTTP GET request r for week schedules by year/week date reference, writing to w.
//
// It takes two request variables--year and week--, which correspond to an ISO 8601 year-week date.
func (sc *ScheduleWeekController) GetByYearWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ystr, wstr, err := weekFromVars(vars)
	if err != nil {
		sc.handleError(w, r, err, "weekFromVars")
		return
	}

	year, week, _, err := utils.ParseIsoWeek(ystr, wstr, "1")
	if err != nil {
		sc.handleError(w, r, err, "ParseIsoWeek")
		return
	}

	sc.makeAndRenderWeek(w, r, year, week)
}

// makeAndRenderWeek makes and renders a week schedule for year and week, writing to w.
func (sc *ScheduleWeekController) makeAndRenderWeek(w http.ResponseWriter, r *http.Request, year, week int) {
	m := models.NewScheduleModel(sc.session)
	ws, err := m.WeekSchedule(year, week, sc.config.Schedule.Sustainer, sc.timeslotURLBuilder)
	if err != nil {
		sc.handleError(w, r, err, "ScheduleModel.WeekSchedule")
		return
	}

	purl, curl, nurl, err := sc.getRelatedScheduleURLs(ws)
	if err != nil {
		sc.handleError(w, r, err, "scheduleController.getRelatedScheduleURLs")
		return
	}

	currentAndNext, err := m.GetCurrentAndNext()
	if err != nil {
		log.Printf("Error from ScheduleModel.GetCurrentAndNext: %v", err)
		// Safe to continue since we can check if CaN is absent in the template
	}

	subtypes, err := sc.session.GetAllShowSubtypes() // TODO(markspolakovs): strictly this should be on ScheduleModel
	if err != nil {
		sc.handleError(w, r, err, "Session.GetAllShowSubtypes")
		return
	}

	data := struct {
		Schedule                  *models.WeekSchedule
		PrevURL, CurrURL, NextURL *url.URL
		CurrentAndNext            *myradio.CurrentAndNext
		StartHour                 int
		Subtypes                  []myradio.ShowSeasonSubtype
	}{
		Schedule:       ws,
		PrevURL:        purl,
		CurrURL:        curl,
		NextURL:        nurl,
		CurrentAndNext: currentAndNext,
		StartHour:      utils.StartHour,
		Subtypes:       subtypes,
	}

	utils.RenderTemplate(w, sc.config.PageContext, data, "schedule_week.tmpl", "elements/current_and_next.tmpl")
}

// getRelatedScheduleURLs gets the URLs for the previous, current, and next schedules relative to ws.
// Any schedule that doesn't exist returns "" as an URL.
// It can fail with err if it can't generate the URLs.
func (sc *ScheduleWeekController) getRelatedScheduleURLs(ws *models.WeekSchedule) (purl, curl, nurl *url.URL, err error) {
	if len(ws.Dates) == 0 {
		err = errors.New("week schedule has no assigned dates")
		return
	}

	cdat := ws.Dates[0]
	if curl, err = sc.weekScheduleURLBuilder(&cdat); err != nil {
		return
	}

	pdat := cdat.AddDate(0, 0, -7)
	if purl, err = sc.weekScheduleURLBuilder(&pdat); err != nil {
		return
	}

	ndat := cdat.AddDate(0, 0, 7)
	nurl, err = sc.weekScheduleURLBuilder(&ndat)
	return
}
