package structs

import (
	"github.com/UniversityRadioYork/myradio-go"
	"time"
)

// ScheduleItem is an interface containing information about one item in a URY schedule.
type ScheduleItem interface {
	// GetName gets the display name of the schedule item.
	// The website config is passed to resolve, for example, the sustainer name.
	GetName(config *Config) string

	// GetStart gets the start time of the schedule item.
	GetStart() time.Time

	// GetFinish gets the finish time of the schedule item.
	GetFinish() time.Time

	// GetBlock gets the block name of the schedule item.
	GetBlock() string
}

// SustainerItem is a struct containing information about a sustainer (filler) item in a URY schedule.
type SustainerItem struct {
	// The start time of this item.
	// This is equal to Timeslot.StartTime if Timeslot is non-nil.
	Start time.Time

	// The end time of this item.
	// This is equal to Timeslot.StartTime + Timeslot.Duration if Timeslot is non-nil.
	Finish time.Time
}

// TimeslotItem adapts myradio.Timeslot into a ScheduleItem.
type TimeslotItem struct {
	Timeslot *myradio.Timeslot
}

// NewTimeslotItem converts a myradio.Timeslot into a TimeslotItem
func NewTimeslotItem(t *myradio.Timeslot) *TimeslotItem {
	return &TimeslotItem{Timeslot: t}
}

/*
 * Implementation of ScheduleItem for SustainerItem
 */

// GetName gets the display name of a SustainerItem.
func (s *SustainerItem) GetName(config *Config) string {
	return config.PageContext.SustainerName
}

// GetStart gets the start time of a SustainerItem.
func (s *SustainerItem) GetStart() time.Time {
	return s.Start
}

// GetFinish gets the finish time of a SustainerItem.
func (s *SustainerItem) GetFinish() time.Time {
	return s.Finish
}

// GetBlock gets the block name of a SustainerItem.
func (s *SustainerItem) GetBlock() string {
	return "sustainer"
}

/*
 * Implementation of ScheduleItem for TimeslotItem
 */

// GetName gets the display name of a TimeslotItem.
func (t *TimeslotItem) GetName(config *Config) string {
	return t.Timeslot.Title
}

// GetStart gets the start time of a TimeslotItem.
func (t *TimeslotItem) GetStart() time.Time {
	return t.Timeslot.StartTime
}

// GetFinish gets the finish time of a TimeslotItem.
func (t *TimeslotItem) GetFinish() time.Time {
	return t.Timeslot.StartTime.Add(t.Timeslot.Duration)
}

// GetBlock gets the block name of a TimeslotItem.
func (t *TimeslotItem) GetBlock() string {
	// TODO(CaptainHayashi): calculate schedule block here.
	return "normal"
}
