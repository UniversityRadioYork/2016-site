package models

import (
	"net/url"
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

// WeekSchedule gets the week schedule with ISO-8601 year year and week number week.
// It also takes sustainer configuration, and a function to use for generating timeslot URLs.
// On success, it returns the fully tabulated WeekSchedule for the given week.
// Otherwise, it returns undefined data and the error causing failure.
func (m *ScheduleModel) WeekSchedule(year, week int, sustainerConfig structs.SustainerConfig, timeslotURLBuilder func(*myradio.Timeslot) (*url.URL, error)) (*WeekSchedule, error) {
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

	return tabulateWeekSchedule(weekStart, weekFinish, filled)
}
