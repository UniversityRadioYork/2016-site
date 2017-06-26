package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
)

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

// addCell adds a cell with rowspan s and item i to the row r.
func (r *WeekScheduleRow) addCell(s uint, i *structs.ScheduleItem) {
	r.Cells = append(r.Cells, WeekScheduleCell{RowSpan: s, Item: i})
}

// showStraddlesDay checks whether a show's start and finish cross over the boundary of a URY day.
func showStraddlesDay(start, finish time.Time) bool {
	nextDayStart := utils.StartOfDayOn(start.AddDate(0, 0, 1))
	return finish.After(nextDayStart)
}

// calculateScheduleBoundaries works out the earliest and latest hours in the schedule that need to display.
// It returns these as a pair of start and finish bound, both in terms of offsets from URY start time.
func calculateScheduleBoundaries(items []structs.ScheduleItem) (sOffset, fOffset int, err error) {
	if len(items) == 0 {
		err = errors.New("calculateScheduleBoundaries: no schedule")
		return
	}

	// These are the boundaries for culling, and are expanded upwards when we find shows that start earlier or finish later than the last-set boundary.
	// Initially they are set to one past their worst case to make the updating logic easier.
	// Since we assert we have a schedule, these values _will_ change.
	sOffset = 24
	fOffset = -1

	for _, s := range items {
		start := s.GetStart()
		finish := s.GetFinish()

		// Any show that isn't a sustainer affects the culling boundaries.
		if s.IsSustainer() {
			continue
		}


		if showStraddlesDay(start, finish) {
			// A show that straddles the day crosses over from the end of a day to the start of the day.
			// This means that we saturate the culling boundaries.
			// As an optimisation we don't need to consider any other show.
			sOffset = 0
			fOffset = 23
			return
		}

		// Otherwise, if its start/finish as offsets from start time are outside the current boundaries, update them.
		so := 0
		so, err = utils.HourToStartOffset(start.Hour())
		if err != nil {
			return
		}
		if so < sOffset {
			sOffset = so
		}
		
		fo := 0
		fo, err = utils.HourToStartOffset(finish.Hour())
		if err != nil {
			return
		}
		if fOffset < fo {
			fOffset = fo
		}
	}

	return
}

// calculateScheduleRows takes a schedule and determines which rows should be displayed.
func calculateScheduleRows(items []structs.ScheduleItem) ([]WeekScheduleRow, error) {
	// Internally, we use a 24-hour array to store our decisions.
	rows := make([]struct {
		MinuteMarks map[int]bool
		Cull        bool
	}, 24)

	// Now decide which rows to cull by calculating boundaries, then marking the rows outside of the boundaries.
	sOffset, fOffset, err := calculateScheduleBoundaries(items)
	if err != nil {
		return nil, err
	}
	if 23 < sOffset || sOffset < 0 || 23 < fOffset || fOffset < 0 || fOffset < sOffset {
		return nil, fmt.Errorf("calculateScheduleRows: row boundaries %d to %d are invalid", sOffset, fOffset)
	}

	// Go through each hour, culling ones before the boundaries, and adding on-the-hour minute marks to the others.
	// Boundaries are inclusive, so cull only things outside of them.
	for i := 0; i < 24; i++ {
		ri, err := utils.StartOffsetToHour(i)
		if err != nil {
			return nil, err
		}
		if i < sOffset || fOffset < i {
			rows[ri].Cull = true
		} else {
			rows[ri].MinuteMarks = map[int]bool{0: true}
		}
	}
	// Calculate the minute marks from non-on-the-hour show starts now.
	for _, item := range items {
		h := item.GetStart().Hour()
		if !rows[h].Cull {
			rows[item.GetStart().Hour()].MinuteMarks[item.GetStart().Minute()] = true
		}
	}

	// Now translate the above into a row table.
	wsrs := []WeekScheduleRow{}
	for i := 0; i < 24; i++ {
		ri, err := utils.StartOffsetToHour(i)
		if err != nil {
			return nil, err
		}

		if rows[ri].Cull {
			continue
		}

		minutes := make([]int, len(rows[ri].MinuteMarks))
		j := 0
		for k := range rows[ri].MinuteMarks {
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

	return wsrs, nil
}

// populateRows fills schedule rows with timeslots.
// It takes the list of schedule start times on the days the schedule spans,
// the slice of rows to populate, and the schedule items to add.
func populateRows(days []time.Time, rows []WeekScheduleRow, items []structs.ScheduleItem) {
	currentItem := 0

	for d, day := range days {
		// We use this to find out when we've gone over midnight
		lastHour := -1
		// And this to find out where the current show started
		thisShowIndex := -1

		// Now, go through all the rows for this day.
		// We have to be careful to make sure we tick over day if we go past midnight.
		for i := range rows {
			if rows[i].Hour < lastHour {
				day = day.AddDate(0, 0, 1)
			}
			lastHour = rows[i].Hour

			rowTime := time.Date(day.Year(), day.Month(), day.Day(), rows[i].Hour, rows[i].Minute, 0, 0, time.Local)

			// Seek forwards if the current show has finished.
			for !items[currentItem].GetFinish().After(rowTime) {
				currentItem++
				thisShowIndex = -1
			}

			// If this is not the first time we've seen this slot,
			// update the rowspan in the first instance's cell and
			// put in a placeholder.
			if thisShowIndex != -1 {
				rows[thisShowIndex].Cells[d].RowSpan++
				rows[i].addCell(0, nil)
			} else {
				thisShowIndex = i
				rows[i].addCell(1, &(items[currentItem]))
			}
		}
	}
}

// WeekSchedule is the type of week schedules.
type WeekSchedule struct {
	// Dates enumerates the dates this week schedule covers.
	Dates []time.Time
	// Table is the actual week table.
	// If there is no schedule for the given week, this will be nil.
	Table []WeekScheduleRow
}

// hasShows asks whether a schedule slice contains any non-sustainer shows.
// It assumes the slice has been filled with sustainer.
func hasShows(schedule []structs.ScheduleItem) bool {
	// This shouldn't happen, but if it does, this is the right thing to
	// do.
	if len(schedule) == 0 {
		return false
	}

	// We know that, if a slice is filled but has no non-sustainer, then
	// the slice will contain only one sustainer item.  So, eliminate the
	// other cases.
	if 1 < len(schedule) || !schedule[0].IsSustainer() {
		return true
	}

	return false
}

// tabulateWeekSchedule creates a schedule table from the given schedule slice.
func tabulateWeekSchedule(start, finish time.Time, schedule []structs.ScheduleItem) (*WeekSchedule, error) {
	days := []time.Time{}
	for d := start; d.Before(finish); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	if !hasShows(schedule) {
		return &WeekSchedule{
			Dates: days,
			Table: nil,
		}, nil
	}

	table, err := calculateScheduleRows(schedule)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	populateRows(days, table, schedule)

	return &WeekSchedule{
		Dates: days,
		Table: table,
	}, nil
}

//
// Controller
//

// ScheduleWeekController is the controller for looking up week schedules.
type ScheduleWeekController struct {
	Controller

	timeslotURLBuilder func(*myradio.Timeslot) (*url.URL, error)
}

// NewScheduleWeekController returns a new ScheduleWeekController with the MyRadio session s,
// router r, and configuration context c.
func NewScheduleWeekController(s *myradio.Session, r *mux.Router, c *structs.Config) *ScheduleWeekController {
	// We pass in the router so we can generate URL reversal functions.
	// Eventually we might want to clean this up, either by passing in
	// something more loosely coupled or handling this at a higher level.
	troute := r.Get("timeslot")
	tbuilder := func(t *myradio.Timeslot) (*url.URL, error) {
		return troute.URLPath("id", strconv.FormatUint(t.TimeslotID, 10))
	}

	return &ScheduleWeekController{
		Controller:         Controller{session: s, config: c},
		timeslotURLBuilder: tbuilder,
	}
}

// makeTimeslotItem creates a TimeslotItem for a given MyRadio timeslot.
func (sc *ScheduleWeekController) makeTimeslotItem(t *myradio.Timeslot) (*structs.TimeslotItem, error) {
	ts, err := structs.NewTimeslotItem(t, sc.timeslotURLBuilder)
	if err == nil && ts == nil {
		return nil, errors.New("NewTimeslotItem created nil timeslot item")
	}
	return ts, err
}

// makeWeekSchedule gets the week schedule for a given ISO year and week.
func (sc *ScheduleWeekController) makeWeekSchedule(yr, wk int) (*WeekSchedule, error) {
	startDate, err := utils.IsoWeekToDate(yr, wk, time.Monday)
	if err != nil {
		return nil, err
	}
	finishDate := startDate.AddDate(0, 0, 7)

	sm := models.NewScheduleWeekModel(sc.session)	
	timeslots, err := sm.Get(yr, wk)
	if err != nil {
		return nil, err
	}

	// Flatten the timeslots into one stream
	flat := []myradio.Timeslot{}
	for d := 1; d <= 7; d++ {
		flat = append(flat, timeslots[d]...)
	}

	// Now start filling from day start to day finish.
	weekStart := utils.StartOfDayOn(startDate)
	weekFinish := utils.StartOfDayOn(finishDate)
	filled, err := structs.FillTimeslotSlice(weekStart, weekFinish, flat, sc.makeTimeslotItem)
	if err != nil {
		return nil, err
	}

	return tabulateWeekSchedule(weekStart, weekFinish, filled)
}

// GetByYearWeek handles the HTTP GET request r for week schedules by year/week date reference, writing to w.
//
// It takes two request variables--year and week--, which correspond to an ISO 8601 year-week date.
func (sc *ScheduleWeekController) GetByYearWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	year, week, err := weekFromVars(vars)
	if err != nil {
		log.Println(err)
		return
	}

	yr, wk, _, err := utils.ParseIsoWeek(year, week, "1")
	if err != nil {
		log.Println(err)
		return
	}

	data, err := sc.makeWeekSchedule(yr, wk)
	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "schedule_week.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
