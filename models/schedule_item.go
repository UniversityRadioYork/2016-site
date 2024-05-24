package models

import (
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/UniversityRadioYork/2016-site/structs"
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

	// ShowArtURL is the URL to the show art image of this show.
	ShowArtURL string
}

// IsSustainer checks whether this schedule item is the URY sustainer.
func (s *ScheduleItem) IsSustainer() bool {
	return s.Block == "sustainer"
}

// NewSustainerItem creates a new sustainer schedule item lasting from start to finish.
// It takes a sustainer config, c, to work out the sustainer name.
func NewSustainerItem(c structs.SustainerConfig, start, finish time.Time) *ScheduleItem {
	return &ScheduleItem{
		Name:       c.Name,
		Desc:       c.Desc,
		Start:      start,
		Finish:     finish,
		Block:      "sustainer",
		PageURL:    "",
		ShowArtURL: c.Image,
	}
}

// NewTimeslotItem converts a myradio.Timeslot into a TimeslotItem.
// It accepts a separate finish time to account for any truncating that occurs when resolving overlaps.
func NewTimeslotItem(t *myradio.Timeslot, finish time.Time, u func(*myradio.Timeslot) (*url.URL, error)) (*ScheduleItem, error) {
	if t == nil {
		return nil, errors.New("NewTimeslotItem: given nil timeslot")
	}

	url, err := u(t)
	if err != nil {
		return nil, err
	}
	return &ScheduleItem{
		Name:       t.Title,
		Desc:       t.Description,
		Start:      t.StartTime,
		Finish:     finish,
		Block:      t.Subtype.Class,
		PageURL:    url.Path,
		ShowArtURL: t.Photo,
	}, nil
}

// scheduleBuilder is an internal type holding information about a schedule slice under construction.
type scheduleBuilder struct {
	// config is the sustainer config to use when creating sustainer slots.
	config structs.SustainerConfig
	// slice is the schedule slice being constructed.
	items []*ScheduleItem
	// nitems is the number of items currently inside the schedule.
	nitems int
	// tbuilder is the function used to create schedule items from timeslots.
	tbuilder func(*myradio.Timeslot, time.Time) (*ScheduleItem, error)
	// err stores any error caused while building the schedule.
	err error
}

// newScheduleBuilder creates an empty schedule builder for nslots shows, given config c and builder tbuilder.
func newScheduleBuilder(c structs.SustainerConfig, tbuilder func(*myradio.Timeslot, time.Time) (*ScheduleItem, error), nslots int) *scheduleBuilder {
	return &scheduleBuilder{
		config: c,
		// nslots slots, (nslots - 1) sustainers in between, and 2 sustainers at the ends.
		items:    make([]*ScheduleItem, ((2 * nslots) + 1)),
		nitems:   0,
		tbuilder: tbuilder,
		err:      nil,
	}
}

// add adds an item to the scheduleBuilder s.
func (s *scheduleBuilder) add(i *ScheduleItem) {
	s.items[s.nitems] = i
	s.nitems++
}

// fill adds a sustainer timeslot between start and finish into the scheduleBuilder s if one needs to be there.
func (s *scheduleBuilder) fill(start, finish time.Time) {
	if start.Before(finish) {
		s.add(NewSustainerItem(s.config, start, finish))
	}
}

// addTimeslot converts a timeslot t to a schedule item, then adds it to the scheduleBuilder s.
// It takes an overlap-adjusted finish, and does not add an item if this adjustment causes t to disappear.
func (s *scheduleBuilder) addTimeslot(t *myradio.Timeslot, finish time.Time) {
	if s.err != nil || !t.StartTime.Before(finish) {
		return
	}

	var ts *ScheduleItem
	if ts, s.err = s.tbuilder(t, finish); s.err != nil {
		return
	}

	s.add(ts)
}

// schedule gets the schedule from a scheduleBuilder, or an err if schedule building failed.
func (s *scheduleBuilder) schedule() ([]*ScheduleItem, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.items[:s.nitems], nil
}

// truncateOverlap clips finish to nextStart if the two overlap and nextShow exists.
// If so, we log an overlap warning, whose content depends on show and nextShow.
// If nextShow is nil, we've overlapped with the end of the schedule, which doesn't need truncating.
func truncateOverlap(finish, nextStart time.Time, show, nextShow *myradio.Timeslot) time.Time {
	if nextShow == nil || !finish.After(nextStart) {
		return finish
	}

	// If the show starts after the next ends then there is no overlap
	if show.StartTime.After(nextStart.Add(nextShow.Duration)) {
		return finish
	}

	log.Println("Truncating" + show.Title)

	log.Printf(
		"Timeslot '%s', ID %d, finishing at %v overlaps with timeslot '%s', ID %d, starting at %v'",
		show.Title,
		show.TimeslotID,
		finish,
		nextShow.Title,
		nextShow.TimeslotID,
		nextStart,
	)

	return nextStart
}

// MakeScheduleSlice converts a slice of Timeslots to a slice of ScheduleItems.
// It does so by filling in any gaps between the start time and the first show, the final show and the finish time, and any two shows.
// Any overlaps are resolved by truncating the timeslot finish time, and dropping it if this makes the timeslot disappear.
// It expects a constructor function for lifting Timeslots (and overlap-adjusted finish times) to TimeslotItems.
// It will return an error if any two shows overlap.
// It presumes the timeslot slice is already sorted in chronological order.
func MakeScheduleSlice(c structs.SustainerConfig, start, finish time.Time, slots []myradio.Timeslot, tbuilder func(*myradio.Timeslot, time.Time) (*ScheduleItem, error)) ([]*ScheduleItem, error) {
	nslots := len(slots)
	if nslots == 0 {
		return []*ScheduleItem{NewSustainerItem(c, start, finish)}, nil
	}

	s := newScheduleBuilder(c, tbuilder, nslots)

	// Now, if possible, start filling between.
	var show, nextShow *myradio.Timeslot
	for i := range slots {
		show = &slots[i]
		rawShowFinish := show.StartTime.Add(show.Duration)

		var nextStart time.Time
		// Is the next start another show, or the end of the schedule?
		if i < nslots-1 {
			nextShow = &slots[i+1]
			nextStart = nextShow.StartTime
		} else {
			nextShow = nil
			nextStart = finish
		}

		showFinish := truncateOverlap(rawShowFinish, nextStart, show, nextShow)
		s.addTimeslot(show, showFinish)
		s.fill(showFinish, nextStart)
	}

	return s.schedule()
}
