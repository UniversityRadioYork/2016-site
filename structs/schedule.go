package structs

// TODO(CaptainHayashi): this probably doesn't belong in structs.

import (
	"errors"
	"fmt"
	"github.com/UniversityRadioYork/myradio-go"
	"net/url"
	"time"
)

// ScheduleItem is an interface containing information about one item in a URY schedule.
type ScheduleItem interface {
	// GetName gets the display name of the schedule item.
	// The website page context is passed to resolve, for example, the sustainer name.
	GetName(context *PageContext) string

	// GetDesc gets the description of the schedule item.
	// The website page context is passed to resolve, for example, the
	// sustainer description.
	GetDesc(context *PageContext) string

	// GetStart gets the start time of the schedule item.
	GetStart() time.Time

	// GetFinish gets the finish time of the schedule item.
	GetFinish() time.Time

	// GetBlock gets the block name of the schedule item.
	GetBlock() string

	// IsSustainer gets whether this schedule item is the URY sustainer.
	IsSustainer() bool

	// HasPage is true when this schedule item has an associated page.
	HasPage() bool

	// GetPageURL gets the root-relative URL to this schedule item's page.
	GetPageURL() string
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
	PageURL  string
}

// NewTimeslotItem converts a myradio.Timeslot into a TimeslotItem.
func NewTimeslotItem(t *myradio.Timeslot, u func(*myradio.Timeslot) (*url.URL, error)) (*TimeslotItem, error) {
	if t == nil {
		return nil, errors.New("NewTimeslotItem: given nil timeslot")
	}

	url, err := u(t)
	if err != nil {
		return nil, err
	}
	return &TimeslotItem{Timeslot: t, PageURL: url.Path}, nil
}

/*
 * Implementation of ScheduleItem for SustainerItem
 */

// GetName gets the display name of a SustainerItem.
func (s *SustainerItem) GetName(context *PageContext) string {
	return context.Sustainer.Name
}

// GetDesc gets the description of a SustainerItem.
func (s *SustainerItem) GetDesc(context *PageContext) string {
	return context.Sustainer.Desc
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

// HasPage gets whether a SustainerItem has a show page (it doesn't).
func (s *SustainerItem) HasPage() bool {
	return false
}

// GetPageUrl gets the root-relative URL to a SustainerItem's (non-existent) show page.
func (s *SustainerItem) GetPageURL() string {
	return ""
}

/*
 * Implementation of ScheduleItem for TimeslotItem
 */

// GetName gets the display name of a TimeslotItem.
func (t *TimeslotItem) GetName(context *PageContext) string {
	return t.Timeslot.Title
}

// GetDesc gets the description of a TimeslotItem.
func (t *TimeslotItem) GetDesc(context *PageContext) string {
	return t.Timeslot.Description
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

// HasPage gets whether a TimeslotItem has a show page (it does).
func (s *TimeslotItem) HasPage() bool {
	return true
}

// GetPageUrl gets the root-relative URL to a TimeslotItem's show page.
func (s *TimeslotItem) GetPageURL() string {
	return s.PageURL
}

/*
 * Schedule filling
 * TODO(CaptainHayashi): This DEFINITELY doesn't belong in structs
 */

// FillTimeslotSlice converts a slice of Timeslots to a slice of ScheduleItems.
// It does so by filling in any gaps between the start time and the first show, the final show and the finish time, and any two shows.
// It expects a constructor function for lifting Timeslots to TimeslotItems.
// It will return an error if any two shows overlap.
// It presumes the timeslot slice is already sorted in chronological order.
func FillTimeslotSlice(start, finish time.Time, slots []myradio.Timeslot, tbuilder func(*myradio.Timeslot) (*TimeslotItem, error)) ([]ScheduleItem, error) {
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
	var err error
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

			items[i], err = tbuilder(first)
			if err != nil {
				return nil, err
			}
			i++
			if firstFinish.Before(second.StartTime) {
				items[i] = NewSustainerItem(firstFinish, second.StartTime)
				i++
			}
			// Don't add second -- it'll either be the next first, or we'll add it at the end.
		}
	}

	lastShow := &slots[nslots-1]
	items[i], err = tbuilder(lastShow)
	if err != nil {
		return nil, err
	}
	i++
	lastFinish := lastShow.StartTime.Add(lastShow.Duration)
	if lastFinish.Before(finish) {
		items[i] = NewSustainerItem(lastFinish, finish)
		i++
	}

	return items[:i], nil
}
