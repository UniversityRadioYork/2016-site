package models

import (
	"errors"
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

// straddlesDay checks whether a show's start and finish cross over the boundary of a URY day.
func straddlesDay(s *ScheduleItem) bool {
	dayBoundary := utils.StartHour
	adjustedStart := s.Start.Add(time.Hour * time.Duration(-dayBoundary))
	adjustedEnd := s.Finish.Add(time.Hour * time.Duration(-dayBoundary))
	straddle := adjustedEnd.Day() != adjustedStart.Day() && s.Finish.Sub(s.Start) > time.Hour
	return straddle
}

// calcScheduleBoundaries gets the offsets of the earliest and latest visible schedule hours.
// It returns these as top and bot respectively.
func calcScheduleBoundaries(items []*ScheduleItem, scheduleStart time.Time) (top, bot utils.StartOffset, err error) {
	if len(items) == 0 {
		err = errors.New("calculateScheduleBoundaries: no schedule")
		return
	}

	// These are the boundaries for culling, and are expanded upwards when we find shows that start earlier or finish later than the last-set boundary.
	// Initially they are set to one past their worst case to make the updating logic easier.
	// Since we assert we have a schedule, these values _will_ change.
	// (Top must be before 00:00 or the populator gets screwed up)
	top = utils.StartOffset(23 - utils.StartHour)
	bot = utils.StartOffset(-1)

	for _, s := range items {
		// Any show that isn't a sustainer affects the culling boundaries.
		if s.IsSustainer() {
			continue
		}

		if straddlesDay(s) {
			if scheduleStart.After(s.Start) {
				//This is the first item on the schedule and straddles the week, so we only set the top of the schedule
				//top = utils.StartOffset(0)
				//Temporarily disabled as this slot doesn't show up on the schedule
				continue
			} else if s.Finish.After(scheduleStart.AddDate(0, 0, 7)) {
				//This is the last item on the schedule and straddles the week, so we only set the bottom of the schedule
				bot = utils.StartOffset(23)
				continue
			} else {
				// An item that straddles the day crosses over from the end of a day to the start of the day.
				// This means that we saturate the culling boundaries.
				// As an optimisation we don't need to consider any other show.
				return utils.StartOffset(0), utils.StartOffset(23), nil
			}
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
		// Disabled this for the 2020 schedule rework - marks.polakovs@
		//if s.Finish.Minute() == 0 && s.Finish.Second() == 0 && s.Finish.Nanosecond() == 0 {
		//	cbot--
		//}

		if bot < cbot {
			bot = cbot
		}
	}

	return
}

type WeekSchedule struct {
	Days         []time.Time
	DayItems     map[time.Time][]*ScheduleItem
	EarliestHour utils.StartOffset
	LatestHour   utils.StartOffset
}
