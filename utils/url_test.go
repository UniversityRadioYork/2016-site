package utils_test

import (
	"testing"

	utils "github.com/UniversityRadioYork/2016-site/utils"
)

// TestParseURL tests whether ParseURL works correctly.
func TestParseURL(t *testing.T) {
	cases := []struct {
		Expected string
		Prefix   string
		URL      string
	}{
		{
			Expected: "/testsite/schedule",
			Prefix:   "testsite",
			URL:      "/schedule",
		},
		{
			Expected: "//schedule",
			Prefix:   "testsite",
			URL:      "//schedule",
		},
		{
			Expected: "schedule",
			Prefix:   "testsite",
			URL:      "schedule",
		},
		{
			Expected: "/schedule",
			Prefix:   "",
			URL:      "/schedule",
		},
		{
			Expected: "//schedule",
			Prefix:   "",
			URL:      "//schedule",
		},
		{
			Expected: "schedule",
			Prefix:   "",
			URL:      "schedule",
		},
		{
			Expected: "/testsite/schedule",
			Prefix:   "/testsite",
			URL:      "/schedule",
		},
		{
			Expected: "/testsite/schedule",
			Prefix:   "testsite/",
			URL:      "/schedule",
		},
	}

	for i, c := range cases {
		got := utils.PrefixURL(c.URL, c.Prefix)
		if c.Expected != got {
			t.Errorf("case %d: expected: %s; got: %s", i+1, c.Expected, got)
		}
	}
}
