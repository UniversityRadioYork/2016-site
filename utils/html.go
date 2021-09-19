package utils

import (
	"github.com/microcosm-cc/bluemonday"
)

// StripHTML strips HTML tags from a string, extracting all plain text.
func StripHTML(htmls string) string {
  sanitizer := bluemonday.StrictPolicy()

  return sanitizer.Sanitize(htmls)

}


