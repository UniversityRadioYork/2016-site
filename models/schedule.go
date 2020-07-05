package models

import (
	"net/url"
	"sort"
	"time"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// ScheduleModel is the model for the schedule controllers.
type ScheduleModel struct {
	Model
}

// NewScheduleModel returns a new ScheduleModel on the MyRadio session s.
func NewScheduleModel(s *myradio.Session) *ScheduleModel {
	// @TODO: Pass in the config options
	return &ScheduleModel{Model{session: s}}
}

// OLD_WeekSchedule gets the week schedule with ISO-8601 year year and week number week.
// It also takes sustainer configuration, and a function to use for generating timeslot URLs.
// On success, it returns the fully tabulated OLD_WeekSchedule for the given week.
// Otherwise, it returns undefined data and the error causing failure.
func (m *ScheduleModel) WeekSchedule(
	year, week int,
	sustainerConfig structs.SustainerConfig,
	timeslotURLBuilder func(*myradio.Timeslot) (*url.URL, error),
) (*WeekSchedule, error) {
	startDate, err := utils.IsoWeekToDate(year, week, time.Monday)
	if err != nil {
		return nil, err
	}
	finishDate := startDate.AddDate(0, 0, 7)

	timeslots, err := m.session.GetWeekSchedule(year, week)
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

	makeTimeslotItem := func(t *myradio.Timeslot, finish time.Time) (*ScheduleItem, error) {
		return NewTimeslotItem(t, finish, timeslotURLBuilder)
	}

	filled, err := MakeScheduleSlice(sustainerConfig, weekStart, weekFinish, flat, makeTimeslotItem)
	if err != nil {
		return nil, err
	}

	top, bot, err := calcScheduleBoundaries(filled, weekStart)

	days, dayItems := m.divideItemsIntoDays(filled, weekStart, weekFinish, top, bot)

	return &WeekSchedule{
		Days:         days,
		DayItems:     dayItems,
		EarliestHour: top,
		LatestHour:   bot,
	}, nil
}

func (m *ScheduleModel) divideItemsIntoDays(
	items []*ScheduleItem,
	start, finish time.Time,
	dayTop, dayBot utils.StartOffset,
) ([]time.Time, map[time.Time][]*ScheduleItem) {

	// These numbers _should_ have come from calcScheduleBoundaries, so they _should_ be legit
	lastHour, err := dayBot.ToHour()
	if err != nil {
		panic(err)
	}
	firstHour, err := dayTop.ToHour()
	if err != nil {
		panic(err)
	}

	days := make([]time.Time, 0)
	dayItems := make(map[time.Time][]*ScheduleItem)

	for d := start; d.Before(finish); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
		dayItems[d] = make([]*ScheduleItem, 0)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Start.Before(items[j].Start)
	})

	// We now know the items are in chronological order.
	today := start
	for i, item := range items {
		nextDay := today.AddDate(0, 0, 1)
		if item.Start.After(nextDay) {
			// The first item "today" is actually tomorrow - it means we've finished a day
			today = nextDay
			nextDay = today.AddDate(0, 0, 1)
		}
		// Check if this item straddles the day
		if item.Finish.After(nextDay) {
			// If so, add it to both today's and tomorrow's schedules
			// Make a copy, so we can fiddle some fields
			tomorrowsVersion := *item

			// If it's a sustainer, extend it out - but only as much as needed.
			if item.IsSustainer() {
				item.Finish = time.Date(
					item.Finish.Year(),
					item.Finish.Month(),
					item.Finish.Day(),
					lastHour,
					0,
					0,
					0,
					item.Finish.Location(),
				)
				tomorrowsVersion.Start = time.Date(
					item.Start.Year(),
					item.Start.Month(),
					item.Start.Day(),
					firstHour,
					0,
					0,
					0,
					tomorrowsVersion.Start.Location(),
				)
				// Only add in tomorrow's version if there isn't already a show then
				if i != len(items)-2 {
					if tomorrowsVersion.Start != items[i+1].Start {
						_, ok := dayItems[nextDay]
						if ok {
							dayItems[nextDay] = append(dayItems[nextDay], &tomorrowsVersion)
						}
					}
				} else {
					_, ok := dayItems[nextDay]
					if ok {
						dayItems[nextDay] = append(dayItems[nextDay], &tomorrowsVersion)
					}
				}
			} else {
				// If it's a real show, make it properly straddle.
				item.ContinuesOnNextDay = true
				item.Finish = nextDay
				tomorrowsVersion.Start = nextDay
				tomorrowsVersion.ContinuedFromPrevDay = true
				// Sanity check
				_, ok := dayItems[nextDay]
				if ok {
					dayItems[nextDay] = append(dayItems[nextDay], &tomorrowsVersion)
				}
			}
		}

		// Fun edge case
		if item.IsSustainer() {
			if item.Start.Before(item.Finish) /* Don't question it. */ {
				dayItems[today] = append(dayItems[today], item)
			}
		} else {
			dayItems[today] = append(dayItems[today], item)
		}
	}

	return days, dayItems
}

// GetCurrentAndNext retrieves the data for the current and next show
func (m *ScheduleModel) GetCurrentAndNext() (*myradio.CurrentAndNext, error) {
	currentAndNext, err := m.session.GetCurrentAndNext()
	if err != nil {
		return nil, err
	}
	return currentAndNext, nil
}
