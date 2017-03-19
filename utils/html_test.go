package utils_test

import (
	utils "github.com/UniversityRadioYork/2016-site/utils"
	"testing"
)

// TestStripHtml tests StripHtml on various HTML fragments.
func TestStripHtml(t *testing.T) {
	cases := []struct {
		Expected string
		Input    string
	}{
		{
			Expected: "",
			Input:    "",
		},
		{
			Expected: "raw test",
			Input:    "raw test",
		},
		{
			Expected: "1\n\n2",
			Input:    "<p>1</p><p>2</p>",
		},

		{
			Expected: "the quick brown fox jumps over the lazy dog",
			Input:    "<p>the quick <strong>brown</strong> fox jumps over the <a href=\"foo\">lazy</a> dog</p>",
		},
	}

	for _, c := range cases {
		got, err := utils.StripHtml(c.Input)
		if err != nil {
			t.Error(err)
		} else if c.Expected != got {
			t.Errorf("expected:\n%s\n\ngot:\n%s", c.Expected, got)
		}
	}
}
