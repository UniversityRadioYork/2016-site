package utils_test

import (
	utils "github.com/UniversityRadioYork/2016-site/utils"
	"testing"
)

// TestParseUrl tests whether ParseUrl works correctly.
func TestParseUrl(t *testing.T) {
	cases := []struct {
		Expected string
		Prefix   string
		Url      string
	}{
		{
			Expected: "/testsite/schedule",
			Prefix:   "testsite",
			Url:      "/schedule",
		},
		{
			Expected: "//schedule",
			Prefix:   "testsite",
			Url:      "//schedule",
		},
		{
			Expected: "schedule",
			Prefix:   "testsite",
			Url:      "schedule",
		},
		{
			Expected: "/schedule",
			Prefix:   "",
			Url:      "/schedule",
		},
		{
			Expected: "//schedule",
			Prefix:   "",
			Url:      "//schedule",
		},
		{
			Expected: "schedule",
			Prefix:   "",
			Url:      "schedule",
		},
		{
			Expected: "/testsite/schedule",
			Prefix:   "/testsite",
			Url:      "/schedule",
		},
		{
			Expected: "/testsite/schedule",
			Prefix:   "testsite/",
			Url:      "/schedule",
		},
	}

	for i, c := range cases {
		got := utils.PrefixUrl(c.Url, c.Prefix)
		if c.Expected != got {
			t.Errorf("case %d: expected: %s; got: %s", i+1, c.Expected, got)
		}
	}
}
