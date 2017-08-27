package models

import (
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

// Get gets the week schedule with ISO-8601 year year and week number week.
//
// On success, it returns the day-split map of timeslots in the week schedule, and nil.
// Otherwise, it returns undefined data and the error causing failure.
func (m *ScheduleModel) Get(year, week int) (map[int][]myradio.Timeslot, error) {
	// TODO(CaptainHayashi): Jukebox filling
	timeslots, err := m.session.GetWeekSchedule(year, week)
	if err != nil {
		return nil, err
	}

	return timeslots, err
}
