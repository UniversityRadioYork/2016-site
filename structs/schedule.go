package structs

// TODO(CaptainHayashi): this probably doesn't belong in structs.

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/UniversityRadioYork/myradio-go"
)

// ScheduleItem contains information about one item in a URY schedule.
type ScheduleItem struct {
	// Name is the display name of the schedule item.
	Name string

	// Desc is the description of the schedule item.
	Desc string

	// Start is the start time of the schedule item.
	Start time.Time

	// Finish is the finish time of the schedule item.
	Finish time.Time

	// Block is the block name of the schedule item.
	Block string

	// PageURL is the root-relative URL to this schedule item's page,
	// or "" if there is no URL.
	PageURL string
}

// IsSustainer checks whether this schedule item is the URY sustainer.
func (s *ScheduleItem) IsSustainer() bool {
	return s.Block == "sustainer"
}

// NewSustainerItem creates a new sustainer schedule item lasting from start to finish.
// It takes a sustainer config, c, to work out the sustainer name.
func NewSustainerItem(c SustainerConfig, start, finish time.Time) *ScheduleItem {
	return &ScheduleItem{
		Name:    c.Name,
		Desc:    c.Desc,
		Start:   start,
		Finish:  finish,
		Block:   "sustainer",
		PageURL: "",
	}
}

// NewTimeslotItem converts a myradio.Timeslot into a TimeslotItem.
func NewTimeslotItem(t *myradio.Timeslot, u func(*myradio.Timeslot) (*url.URL, error)) (*ScheduleItem, error) {
	if t == nil {
		return nil, errors.New("NewTimeslotItem: given nil timeslot")
	}

	url, err := u(t)
	if err != nil {
		return nil, err
	}
	return &ScheduleItem{
		Name:    t.Title,
		Desc:    t.Description,
		Start:   t.StartTime,
		Finish:  t.StartTime.Add(t.Duration),
		Block:   "regular", // TODO(MattWindsor91): get this from elsewhere
		PageURL: url.Path,
	}, nil
}

// FillTimeslotSlice converts a slice of Timeslots to a slice of ScheduleItems.
// It does so by filling in any gaps between the start time and the first show, the final show and the finish time, and any two shows.
// It expects a constructor function for lifting Timeslots to TimeslotItems.
// It will return an error if any two shows overlap.
// It presumes the timeslot slice is already sorted in chronological order.
func FillTimeslotSlice(c SustainerConfig, start, finish time.Time, slots []myradio.Timeslot, tbuilder func(*myradio.Timeslot) (*ScheduleItem, error)) ([]*ScheduleItem, error) {
	nslots := len(slots)

	// The maximum possible number of items is 2(nslots) + 1:
	// nslots slots, (nslots - 1) sustainers in between, and 2 sustainers at the ends.
	items := make([]*ScheduleItem, (2*len(slots))+1)

	// Now deal with the easy case--no slots.
	if nslots == 0 {
		items[0] = NewSustainerItem(c, start, finish)
		return items, nil
	}

	// Otherwise, we now have to do some actual filling.
	i := 0

	// First, work out if we need to fill before the first show.
	firstShow := slots[0]
	if start.Before(firstShow.StartTime) {
		items[i] = NewSustainerItem(c, start, firstShow.StartTime)
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
				items[i] = NewSustainerItem(c, firstFinish, second.StartTime)
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
		items[i] = NewSustainerItem(c, lastFinish, finish)
		i++
	}

	return items[:i], nil
}
