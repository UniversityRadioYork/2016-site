package structs

// TODO(CaptainHayashi): this probably doesn't belong in structs.

import (
	"fmt"
	"github.com/UniversityRadioYork/myradio-go"
	"time"
)

// ScheduleItem is an interface containing information about one item in a URY schedule.
type ScheduleItem interface {
	// GetName gets the display name of the schedule item.
	// The website page context is passed to resolve, for example, the sustainer name.
	GetName(context *PageContext) string

	// GetStart gets the start time of the schedule item.
	GetStart() time.Time

	// GetFinish gets the finish time of the schedule item.
	GetFinish() time.Time

	// GetBlock gets the block name of the schedule item.
	GetBlock() string

	// IsSustainer gets whether this schedule item is the URY sustainer.
	IsSustainer() bool
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

// NewSustainerItem creates a new SustainerItem.
func NewSustainerItem(start, finish time.Time) *SustainerItem {
	return &SustainerItem{Start: start, Finish: finish}
}

// TimeslotItem adapts myradio.Timeslot into a ScheduleItem.
type TimeslotItem struct {
	Timeslot *myradio.Timeslot
}

// NewTimeslotItem converts a myradio.Timeslot into a TimeslotItem.
func NewTimeslotItem(t *myradio.Timeslot) *TimeslotItem {
	return &TimeslotItem{Timeslot: t}
}

/*
 * Implementation of ScheduleItem for SustainerItem
 */

// GetName gets the display name of a SustainerItem.
func (s *SustainerItem) GetName(context *PageContext) string {
	return context.SustainerName
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

// IsSustainer gets whether a SustainerItem is sustainer (it is).
func (s *SustainerItem) IsSustainer() bool {
	return true
}

/*
 * Implementation of ScheduleItem for TimeslotItem
 */

// GetName gets the display name of a TimeslotItem.
func (t *TimeslotItem) GetName(context *PageContext) string {
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

// IsSustainer gets whether a TimeslotItem is sustainer (it isn't).
func (t *TimeslotItem) IsSustainer() bool {
	return false
}

/*
 * Schedule filling
 * TODO(CaptainHayashi): This DEFINITELY doesn't belong in structs
 */

// FillTimeslotSlice converts a slice of Timeslots to a slice of ScheduleItems.
// It does so by filling in any gaps between the start time and the first show, the final show and the finish time, and any two shows.
// It will return an error if any two shows overlap.
// It presumes the timeslot slice is already sorted in chronological order.
func FillTimeslotSlice(start, finish time.Time, slots []myradio.Timeslot) ([]ScheduleItem, error) {
	nslots := len(slots)

	// The maximum possible number of items is 2(nslots) + 1:
	// nslots slots, (nslots - 1) sustainers in between, and 2 sustainers at the ends.
	items := make([]ScheduleItem, (2*len(slots))+1)

	// Now deal with the easy case--no slots.
	if nslots == 0 {
		items[0] = NewSustainerItem(start, finish)
		return items, nil
	}

	// Otherwise, we now have to do some actual filling.
	i := 0

	// First, work out if we need to fill before the first show.
	firstShow := slots[0]
	if start.Before(firstShow.StartTime) {
		items[i] = NewSustainerItem(start, firstShow.StartTime)
		i++
	}

	// Now, if possible, start filling between.
	// This will add all but the last show.
	for j := range slots {
		if j < nslots-1 {
			first := &slots[j]
			second := &slots[j+1]

			firstFinish := first.StartTime.Add(first.Duration)
			if firstFinish.After(second.StartTime) {
				return nil, fmt.Errorf(
					"Timeslot '%s', ID %d, finishing at %v overlaps with timeslot '%s', ID %d, starting at %v'",
					first.Title,
					first.TimeslotID,
					firstFinish,
					second.Title,
					second.TimeslotID,
					second.StartTime,
				)
			}

			items[i] = NewTimeslotItem(first)
			i++
			if firstFinish.Before(second.StartTime) {
				items[i] = NewSustainerItem(firstFinish, second.StartTime)
				i++
			}
			// Don't add second -- it'll either be the next first, or we'll add it at the end.
		}
	}

	lastShow := &slots[nslots-1]
	items[i] = NewTimeslotItem(lastShow)
	i++
	lastFinish := lastShow.StartTime.Add(lastShow.Duration)
	if lastFinish.Before(finish) {
		items[i] = NewSustainerItem(lastFinish, finish)
		i++
	}

	return items[:i], nil
}
