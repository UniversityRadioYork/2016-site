package controllers

import (
	//"github.com/UniversityRadioYork/2016-site/models"
	"errors"
	"fmt"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ScheduleWeekController is the controller for looking up week schedules.
type ScheduleWeekController struct {
	Controller
}

// NewScheduleWeekController returns a new ShowController with the MyRadio session s
// and configuration context c.
func NewScheduleWeekController(s *myradio.Session, c *structs.Config) *ScheduleWeekController {
	return &ScheduleWeekController{Controller{session: s, config: c}}
}

// weekFromVars extracts the year, and week strings from vars.
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

// weekdayFromVars extracts the year, week, and weekday strings from vars.
func weekdayFromVars(vars map[string]string) (string, string, string, error) {
	y, ok := vars["year"]
	if !ok {
		return "", "", "", errors.New("no year provided")
	}
	w, ok := vars["week"]
	if !ok {
		return "", "", "", errors.New("no week provided")
	}
	d, ok := vars["weekday"]
	if !ok {
		return "", "", "", errors.New("no weekday provided")
	}

	return y, w, d, nil
}

// parseIsoWeek parses an ISO weekday from year, week, and weekday strings.
// It performs bounds checking.
func parseIsoWeek(year, week, weekday string) (int, int, time.Weekday, error) {
	y, err := strconv.Atoi(year)
	if err != nil {
		return 0, 0, 0, err
	}
	if y < 0 {
		return 0, 0, 0, fmt.Errorf("Invalid year: %d", y)
	}

	w, err := strconv.Atoi(week)
	if err != nil {
		return 0, 0, 0, err
	}
	if w < 1 || 53 < w {
		return 0, 0, 0, fmt.Errorf("Invalid week: %d", w)
	}

	// Two-stage conversion: first to int, then to Weekday.
	// Go treats Sunday as day 0: we must correct this grave mistake.
	dI, err := strconv.Atoi(weekday)
	if err != nil {
		return 0, 0, 0, err
	}
	if dI < 1 || 7 < dI {
		return 0, 0, 0, fmt.Errorf("Invalid day: %d", dI)
	}

	var d time.Weekday
	if dI == 7 {
		d = time.Sunday
	} else {
		d = time.Weekday(dI)
	}

	return y, w, d, nil
}

// isoWeekToDate interprets year, week, and weekday strings as an ISO weekday.
// The time is set to local midnight.
func isoWeekToDate(year, week, weekday string) (time.Time, error) {
	// This is based on the calculation given at:
	// https://en.wikipedia.org/wiki/ISO_week_date#Calculating_a_date_given_the_year.2C_week_number_and_weekday

	y, w, d, err := parseIsoWeek(year, week, weekday)
	if err != nil {
		return time.Time{}, err
	}

	// We need to find the first week in the year.
	// This always contains the 4th of January, so find that, and get
	// ISOWeek on it.
	fj := time.Date(y, time.January, 4, 0, 0, 0, 0, time.Local)
	fjWeekday := fj.Weekday()

	// Sanity check to make sure time (and our intuition) is still working.
	fjYear, fjWeek := fj.ISOWeek()
	if fjYear != y {
		return time.Time{}, fmt.Errorf("ISO weekday year %d != calendar year %d!", fjYear, y)
	}
	if fjWeek != 1 {
		return time.Time{}, fmt.Errorf("ISO weekday week of 4 Jan (%d) not week 1!", fjWeek)
	}

	// The ISO 8601 ordinal date, which may belong to the next or previous
	// year.
	ord := (w * 7) + int(d) - (int(fjWeekday) + 3)

	// The ordinal date is just the number of days since 1 Jan y plus one,
	// so calculate the year from that.
	oj := time.Date(y, time.January, 1, 0, 0, 0, 0, time.Local)
	return oj.AddDate(0, 0, ord-1), nil
}

// Get handles the HTTP GET request r for all shows, writing to w.
//
// ScheduleWeek's Get takes three request variables--year, week, and weekday--,
// which correspond to an ISO 8601 weekday-format time.
func (sc *ScheduleWeekController) Get(w http.ResponseWriter, r *http.Request) {
	//sm := models.NewScheduleModel(sc.session)

	vars := mux.Vars(r)

	year, week, err := weekFromVars(vars)
	if err != nil {
		log.Println(err)
		return
	}

	startDate, err := isoWeekToDate(year, week, "1")
	if err != nil {
		log.Println(err)
		return
	}
	endDate := startDate.AddDate(0, 0, 7)

	//timeslots, err := sm.GetShow(id)
	//if err != nil {
	//	//@TODO: Do something proper here, render 404 or something
	//	log.Println(err)
	//	return
	//}

	data := struct {
		StartDate time.Time
		EndDate   time.Time
	}{
		StartDate: startDate,
		EndDate:   endDate,
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "schedule_week.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
