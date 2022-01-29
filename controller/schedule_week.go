package controller

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/UniversityRadioYork/2016-site/compat"
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/utils"
	view "github.com/UniversityRadioYork/2016-site/views"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
)

func (c *Controller) HandleScheduleWeek(r *mux.Router) http.HandlerFunc {
	wbuilder := func(t *time.Time) (*url.URL, error) {

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
	v := view.View{
		MainTmpl: "schedule_week.tmpl",
		AddTmpls: []string{"elements/current_and_next.tmpl"},
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		today := time.Now()
		year, week := today.ISOWeek()

		if _, set := vars["week"]; set {
			ystr, wstr, err := weekFromVars(vars)
			if err != nil {
				log.Println(err)
				return
			}
			year, week, _, err = utils.ParseIsoWeek(ystr, wstr, "1")
			if err != nil {
				log.Println(err)
				return
			}
		}

		m := models.NewScheduleModel(c.Session)
		// TODO: Remove compat patch
		ws, err := m.WeekSchedule(year, week, compat.OldConfig().Schedule.Sustainer, tbuilder)
		if err != nil {
			c.HandleNotFound(w, r)
			return
		}

		if len(ws.Dates) == 0 {
			log.Println("week schedule has no assigned dates")
			c.HandleNotFound(w, r)
			return
		}

		cdat := ws.Dates[0]
		curl, err := wbuilder(&cdat)
		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}

		pdat := cdat.AddDate(0, 0, -7)
		purl, err := wbuilder(&pdat)
		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}

		ndat := cdat.AddDate(0, 0, 7)
		nurl, err := wbuilder(&ndat)
		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}

		currentAndNext, err := m.GetCurrentAndNext()
		if err != nil {
			log.Println(err)
			// Safe to continue since we can check if CaN is absent in the template
		}

		subtypes, err := c.Session.GetAllShowSubtypes()
		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
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

		err = v.Render(w, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func weekFromVars(vars map[string]string) (string, string, error) {
	y, ok := vars["year"]
	if !ok {
		return "", "", errors.New("no year provided")
	}
	w, ok := vars["week"]
	if !ok {
		return "", "", errors.New("no week provided")
	}

	return y, w, nil
}
