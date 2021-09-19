package utils_test

import (
	"testing"

	utils "github.com/UniversityRadioYork/2016-site/utils"
)

// TestStripHtml tests StripHtml on various HTML fragments.
func TestStripHTML(t *testing.T) {
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
			Expected: "1\n2",
			Input:    "<p>1</p>\n<p>2</p>",
		},

		{
			Expected: "the quick brown fox jumps over the lazy dog",
			Input:    "<p>the quick <strong>brown</strong> fox jumps over the <a href=\"foo\">lazy</a> dog</p>",
		},
	}

	for _, c := range cases {
		got, err := utils.StripHTML(c.Input)
		if err != nil {
			t.Error(err)
		} else if c.Expected != got {
			t.Errorf("expected:\n%s\n\ngot:\n%s", c.Expected, got)
		}
	}
}
