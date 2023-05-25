package utils

import (
	"regexp"
  "fmt"
)

// Extract first em tag from html string
func ExtractFirstEmTagContent(htmlStr string) (string, error) {
	re := regexp.MustCompile(`(?i)<em>(.*?)</em>`)
	match := re.FindStringSubmatch(htmlStr)
	if len(match) == 2 {
		return match[1], nil
	}
	return "", fmt.Errorf("no <em> tag found")
}
