package controllers

import (
	"errors"
	"fmt"
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

//
// Date manipulation functions
// TODO(CaptainHayashi): move
//

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
func isoWeekToDate(year, week int, weekday time.Weekday) (time.Time, error) {
	// This is based on the calculation given at:
	// https://en.wikipedia.org/wiki/ISO_week_date#Calculating_a_date_given_the_year.2C_week_number_and_weekday

	// We need to find the first week in the year.
	// This always contains the 4th of January, so find that, and get
	// ISOWeek on it.
	fj := time.Date(year, time.January, 4, 0, 0, 0, 0, time.Local)

	// Correct Go's stupid Sunday is 0 decision, making the weekdays ISO 8601 compliant
	intWeekday := int(weekday)
	if intWeekday == 0 {
		intWeekday = 7
	}
	fjWeekday := int(fj.Weekday())
	if fjWeekday == 0 {
		fjWeekday = 7
	}

	// Sanity check to make sure time (and our intuition) is still working.
	fjYear, fjWeek := fj.ISOWeek()
	if fjYear != year {
		return time.Time{}, fmt.Errorf("ISO weekday year %d != calendar year %d!", fjYear, year)
	}
	if fjWeek != 1 {
		return time.Time{}, fmt.Errorf("ISO weekday week of 4 Jan (%d) not week 1!", fjWeek)
	}

	// The ISO 8601 ordinal date, which may belong to the next or previous
	// year.
	ord := (week * 7) + intWeekday - (fjWeekday + 3)

	// The ordinal date is just the number of days since 1 Jan y plus one,
	// so calculate the year from that.
	oj := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	return oj.AddDate(0, 0, ord-1), nil
}

// uryStartOfDayOn gets the URY start of day on a given date.
func uryStartOfDayOn(date time.Time) time.Time {
	y, m, d := date.Date()
	// TODO(CaptainHayashi): de-hardcode the hour?
	return time.Date(y, m, d, 6, 0, 0, 0, time.Local)
}

//
// Week schedule algorithm
// TODO(CaptainHayashi): move?
//

// WeekScheduleCell represents one cell in the week schedule.
type WeekScheduleCell struct {
	// Number of rows this cell spans.
	// If 0, this is a continuation from a cell further up.
	RowSpan uint

	// Pointer to the timeslot in this cell, if any.
	// Will be nil if 'RowSpan' is 0.
	Item *structs.ScheduleItem
}

// WeekScheduleRow represents one row in the week schedule.
type WeekScheduleRow struct {
	// The hour of the row (0..23).
	Hour int
	// The minute of the show (0..59).
	Minute int
	// The cells inside this row.
	Cells []WeekScheduleCell
}

// calculateScheduleRows takes a schedule and determines which rows should be displayed.
func calculateScheduleRows(items []structs.ScheduleItem) []WeekScheduleRow {
	// Internally, we use a 24-hour array to store our decisions.
	rows := make([]struct {
		MinuteMarks     map[int]bool
		HasNonSustainer bool
		Cull            bool
	}, 24)

	// Populate with all of the hours of the URY day.
	for i := 0; i < 24; i++ {
		rows[i].MinuteMarks = map[int]bool{0: true}
	}

	for _, s := range items {
		start := s.GetStart()
		rows[start.Hour()].MinuteMarks[start.Minute()] = true
		if !s.IsSustainer() {
			rows[start.Hour()].HasNonSustainer = true
		}
	}

	// Now decide which rows to cull by working forwards and backwards
	for i := 0; i < 24; i++ {
		ri := (i + 6) % 24

		if rows[ri].HasNonSustainer {
			break
		}
		rows[ri].Cull = true
	}
	for i := 0; i < 24; i++ {
		ri := ((24 + 6) - i) % 24

		if rows[ri].HasNonSustainer {
			break
		}
		rows[ri].Cull = true
	}

	// Now translate the above into a row table
	wsrs := []WeekScheduleRow{}
	for i := 0; i < 24; i++ {
		ri := (i + 6) % 24
		if rows[ri].Cull {
			continue
		}

		minutes := make([]int, len(rows[ri].MinuteMarks))
		j := 0
		for k, _ := range rows[ri].MinuteMarks {
			minutes[j] = k
			j++
		}
		sort.Ints(minutes)

		hwsrs := make([]WeekScheduleRow, len(minutes))
		for j, m := range minutes {
			hwsrs[j] = WeekScheduleRow{Hour: ri, Minute: m, Cells: []WeekScheduleCell{}}
		}

		wsrs = append(wsrs, hwsrs...)
	}

	return wsrs
}

// populateRows fills schedule rows with timeslots.
// It takes local midnight on the start and end days of the schedule to fill.
func populateRows(startMidnight, endMidnight time.Time, rows []WeekScheduleRow, items []structs.ScheduleItem) {
	// How many days does this timetable actually span?
	scheduleSpan := endMidnight.Sub(startMidnight)
	scheduleDays := int(scheduleSpan / time.Hour / 24)

	currentItem := 0

	// Handle each day individually
	for d := 0; d < scheduleDays; d++ {
		dayMidnight := startMidnight.AddDate(0, 0, d)

		// We use this to find out when we've gone over midnight
		lastHour := -1
		// And this to find out where the current show started
		thisShowIndex := -1

		// Now, go through all the rows for this day.
		// We have to be careful to make sure we tick over dayMidnight if we go past midnight.
		for i := range rows {
			if rows[i].Hour < lastHour {
				dayMidnight = dayMidnight.AddDate(0, 0, 1)
			}
			lastHour = rows[i].Hour

			rowTime := time.Date(dayMidnight.Year(), dayMidnight.Month(), dayMidnight.Day(), rows[i].Hour, rows[i].Minute, 0, 0, time.Local)

			// Seek forwards if the current show has finished.
			for !items[currentItem].GetFinish().After(rowTime) {
				currentItem++
				thisShowIndex = -1
			}

			// If this is not the first time we've seen this slot, update its rowspan
			// and put in a placeholder.
			if thisShowIndex != -1 {
				rows[thisShowIndex].Cells[d].RowSpan++
				rows[i].Cells = append(rows[i].Cells, WeekScheduleCell{RowSpan: 0, Item: nil})
			} else {
				thisShowIndex = i
				rows[i].Cells = append(rows[i].Cells, WeekScheduleCell{RowSpan: 1, Item: &(items[currentItem])})
			}
		}
	}
}

//
// Controller
//

// ScheduleWeekController is the controller for looking up week schedules.
type ScheduleWeekController struct {
	Controller
}

// NewScheduleWeekController returns a new ShowController with the MyRadio session s
// and configuration context c.
func NewScheduleWeekController(s *myradio.Session, c *structs.Config) *ScheduleWeekController {
	return &ScheduleWeekController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for all shows, writing to w.
//
// ScheduleWeek's Get takes two request variables--year and week--,
// which correspond to an ISO 8601 year-week date.
func (sc *ScheduleWeekController) GetByYearWeek(w http.ResponseWriter, r *http.Request) {
	sm := models.NewScheduleWeekModel(sc.session)

	vars := mux.Vars(r)

	year, week, err := weekFromVars(vars)
	if err != nil {
		log.Println(err)
		return
	}

	yr, wk, dy, err := parseIsoWeek(year, week, "1")
	if err != nil {
		log.Println(err)
		return
	}

	startDate, err := isoWeekToDate(yr, wk, dy)
	if err != nil {
		log.Println(err)
		return
	}
	finishDate := startDate.AddDate(0, 0, 7)

	log.Printf("getting year %d week %d\n", yr, wk)
	timeslots, err := sm.Get(yr, wk)
	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	// Flatten the timeslots into one stream
	flat := []myradio.Timeslot{}
	for d := 1; d <= 7; d++ {
		flat = append(flat, timeslots[d]...)
	}

	// Now start filling from URY start to URY start.
	startUry := uryStartOfDayOn(startDate)
	finishUry := uryStartOfDayOn(finishDate)
	filled, err := structs.FillTimeslotSlice(startUry, finishUry, flat)
	if err != nil {
		log.Println(err)
		return
	}

	table := calculateScheduleRows(filled)
	populateRows(startDate, finishDate, table, filled)

	data := struct {
		StartDate  time.Time
		FinishDate time.Time
		Table      []WeekScheduleRow
	}{
		StartDate:  startUry,
		FinishDate: finishUry,
		Table:      table,
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "schedule_week.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
