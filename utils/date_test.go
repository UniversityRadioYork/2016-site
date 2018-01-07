package utils_test

import (
	"testing"
	"time"

	"github.com/UniversityRadioYork/2016-site/utils"
)

// offset is the type of day and hour offsets used to test FormatWeekRelativeTo.
type offset struct {
	d int
	h time.Duration
}

// buildOffsetsTable builds a table of day and hour offsets to use to test FormatWeekRelativeTo.
// It contains the cartesian product of every day offset 0--6 and every hour offset 0--23.
func buildOffsetsTable(t *testing.T) []offset {
	t.Helper()

	offsets := []offset{}
	for d := 0; d < 7; d++ {
		for h := 0; h < 23; h++ {
			offsets = append(offsets, offset{d, time.Duration(h)})
		}
	}

	return offsets
}

func testFormatWeekRelativeToRunner(t *testing.T, now time.Time, weekOffset int, expected string) {
	t.Helper()

	start := now.AddDate(0, 0, weekOffset*7)

	ot := buildOffsetsTable(t)
	for _, so := range ot {
		rstart := start.AddDate(0, 0, so.d).Add(time.Hour * so.h)
		for _, no := range ot {
			rnow := now.AddDate(0, 0, no.d).Add(time.Hour * no.h)
			if result := utils.FormatWeekRelativeTo(rstart, rnow); result != expected {
				t.Fatalf("mismatch: start week+%d, day+%d, hour+%d; reference day+= %d, hour += %d; expected '%s', got '%s'", weekOffset, so.d, so.h, no.d, no.h, expected, result)
			}
		}
	}
}

// TestFormatWeekRelativeTo tests FormatWeekRelativeTo.
// It does so by comparing a single known time against a series of week-offset times, applying day and hour offsets to each to try to find edge cases.
func TestFormatWeekRelativeTo(t *testing.T) {
	cases := []struct {
		weekOffset int
		result     string
	}{
		{-2, "01 Jan 2018 to 07 Jan 2018"},
		{-1, "last week"},
		{0, "this week"},
		{1, "next week"},
		{2, "29 Jan 2018 to 04 Feb 2018"},
	}

	// This should always be a Monday, and line up with the strings above.
	now := time.Date(2018, time.January, 15, 0, 0, 0, 0, time.UTC)

	for _, c := range cases {
		testFormatWeekRelativeToRunner(t, now, c.weekOffset, c.result)
	}
}
