package models

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/UniversityRadioYork/2016-site/utils"
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
	Item *ScheduleItem
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
func (r *WeekScheduleRow) addCell(s uint, i *ScheduleItem) {
	r.Cells = append(r.Cells, WeekScheduleCell{RowSpan: s, Item: i})
}

// straddlesDay checks whether a show's start and finish cross over the boundary of a URY day.
func straddlesDay(s *ScheduleItem) bool {
	nextDayStart := utils.StartOfDayOn(s.Start.AddDate(0, 0, 1))
	return s.Finish.After(nextDayStart)
}

// calcScheduleBoundaries gets the offsets of the earliest and latest visible schedule hours.
// It returns these as top and bot respectively.
func calcScheduleBoundaries(items []*ScheduleItem) (top, bot utils.StartOffset, err error) {
	if len(items) == 0 {
		err = errors.New("calculateScheduleBoundaries: no schedule")
		return
	}

	// These are the boundaries for culling, and are expanded upwards when we find shows that start earlier or finish later than the last-set boundary.
	// Initially they are set to one past their worst case to make the updating logic easier.
	// Since we assert we have a schedule, these values _will_ change.
	top = utils.StartOffset(24)
	bot = utils.StartOffset(-1)

	for _, s := range items {
		// Any show that isn't a sustainer affects the culling boundaries.
		if s.IsSustainer() {
			continue
		}

		if straddlesDay(s) {
			// An item that straddles the day crosses over from the end of a day to the start of the day.
			// This means that we saturate the culling boundaries.
			// As an optimisation we don't need to consider any other show.
			return utils.StartOffset(0), utils.StartOffset(23), nil
		}

		// Otherwise, if its start/finish as offsets from start time are outside the current boundaries, update them.
		var ctop utils.StartOffset
		if ctop, err = utils.HourToStartOffset(s.Start.Hour()); err != nil {
			return
		}
		if ctop < top {
			top = ctop
		}

		var cbot utils.StartOffset
		if cbot, err = utils.HourToStartOffset(s.Finish.Hour()); err != nil {
			return
		}
		// cbot is the offset of the hour in which the item finishes.
		// This is _one past_ the last row the item occupies if the item ends cleanly at :00:00.
		if s.Finish.Minute() == 0 && s.Finish.Second() == 0 && s.Finish.Nanosecond() == 0 {
			cbot--
		}

		if bot < cbot {
			bot = cbot
		}
	}

	return
}

// rowDecision is an internal type recording information about which rows to display in the week schedule.
// It records, for one hour, the minute rows (00, 30, etc) that are switched 'on' for that row.
type rowDecision map[int]struct{}

// visible checks if the hour represented by row decision r is to be shown on the schedule.
func (r rowDecision) visible() bool {
	// Each visible row has its on-the-hour row set.
	_, visible := r[0]
	return visible
}

// mark adds a mark for the given minute to row decision r.
func (r rowDecision) mark(minute int) {
	r[minute] = struct{}{}
}

// toRow converts row decision r to a slice of schedule rows for the given hour.
func (r rowDecision) toRows(hour int) []WeekScheduleRow {
	minutes := make([]int, len(r))
	j := 0
	for k := range r {
		minutes[j] = k
		j++
	}
	sort.Ints(minutes)

	rows := make([]WeekScheduleRow, len(minutes))
	for j, m := range minutes {
		rows[j] = WeekScheduleRow{Hour: hour, Minute: m, Cells: []WeekScheduleCell{}}
	}
	return rows
}

// initRowDecisions creates 24 rowDecisions, from schedule start to schedule end.
// Each is marked as visble or invisible depending on the offsets top and bot.
func initRowDecisions(top, bot utils.StartOffset) ([]rowDecision, error) {
	// Make sure the offsets are valid.
	if !top.Valid() || !bot.Valid() {
		return nil, fmt.Errorf("initRowDecisions: row boundaries %d to %d are invalid", int(top), int(bot))
	}

	rows := make([]rowDecision, 24)

	// Go through each hour, culling ones before the boundaries, and adding on-the-hour minute marks to the others.
	// Boundaries are inclusive, so cull only things outside of them.
	for i := utils.StartOffset(0); i < utils.StartOffset(24); i++ {
		h, err := i.ToHour()
		if err != nil {
			return nil, err
		}

		rows[h] = rowDecision{}
		if top <= i && i <= bot {
			// This has the effect of making the row visible.
			rows[h].mark(0)
		}
	}

	return rows, nil
}

// addItemsToRowDecisions populates the row decision list rows with minute marks from schedule items not starting on the hour.
func addItemsToRowDecisions(rows []rowDecision, items []*ScheduleItem) {
	for _, item := range items {
		h := item.Start.Hour()
		if rows[h].visible() {
			rows[h].mark(item.Start.Minute())
		}
	}
}

// rowDecisionsToRows generates rows based on the per-hourly row decisions in rdecs.
func rowDecisionsToRows(rdecs []rowDecision) ([]WeekScheduleRow, error) {
	rows := []WeekScheduleRow{}

	for i := utils.StartOffset(0); i < utils.StartOffset(24); i++ {
		h, err := i.ToHour()
		if err != nil {
			return nil, err
		}

		if rdecs[h].visible() {
			rows = append(rows, rdecs[h].toRows(h)...)
		}
	}

	return rows, nil
}

// initScheduleRows takes a schedule and determines which rows should be displayed.
func initScheduleRows(items []*ScheduleItem) ([]WeekScheduleRow, error) {
	top, bot, err := calcScheduleBoundaries(items)
	if err != nil {
		return nil, err
	}

	rdecs, err := initRowDecisions(top, bot)
	if err != nil {
		return nil, err
	}
	addItemsToRowDecisions(rdecs, items)

	return rowDecisionsToRows(rdecs)
}

// populateRows fills schedule rows with timeslots.
// It takes the list of schedule start times on the days the schedule spans,
// the slice of rows to populate, and the schedule items to add.
func populateRows(days []time.Time, rows []WeekScheduleRow, items []*ScheduleItem) {
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
			for !items[currentItem].Finish.After(rowTime) {
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
				rows[i].addCell(1, items[currentItem])
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
func hasShows(schedule []*ScheduleItem) bool {
	// This shouldn't happen, but if it does, this is the right thing to do.
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
func tabulateWeekSchedule(start, finish time.Time, schedule []*ScheduleItem) (*WeekSchedule, error) {
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

	rows, err := initScheduleRows(schedule)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	populateRows(days, rows, schedule)

	return &WeekSchedule{
		Dates: days,
		Table: rows,
	}, nil
}
