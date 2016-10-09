package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

// ScheduleWeekModel is the model for the ScheduleWeek controller.
type ScheduleWeekModel struct {
	Model
}

// NewScheduleWeekModel returns a new ScheduleWeekModel on the MyRadio session s.
func NewScheduleWeekModel(s *myradio.Session) *ScheduleWeekModel {
	// @TODO: Pass in the config options
	return &ScheduleWeekModel{Model{session: s}}
}

// Get gets the week schedule with ISO-8601 year year and week number week.
//
// On success, it returns the day-split map of timeslots in the week schedule, and nil.
// Otherwise, it returns undefined data and the error causing failure.
func (m *ScheduleWeekModel) Get(year, week int) (map[int][]myradio.Timeslot, error) {
	// TODO(CaptainHayashi): Jukebox filling
	timeslots, err := m.session.GetWeekSchedule(year, week)
	if err != nil {
		return nil, err
	}

	return timeslots, err
}
