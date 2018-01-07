package utils_test

import (
	"testing"
	"time"

	"github.com/UniversityRadioYork/2016-site/utils"
)

func TestFormatWeekRelativeTo_WeekStrides(t *testing.T) {
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
		start := now.AddDate(0, 0, c.weekOffset*7)
		for i := 0; i < 7; i++ {
			rstart := start.AddDate(0, 0, i)
			for j := 0; j < 7; j++ {
				rnow := now.AddDate(0, 0, j)
				if result := utils.FormatWeekRelativeTo(rstart, rnow); result != c.result {
					t.Errorf("mismatch: start week offset = %d, day offset = %d; reference day offset = %d; expected '%s', got '%s'", c.weekOffset, i, j, c.result, result)
				}
			}
		}
	}
}
