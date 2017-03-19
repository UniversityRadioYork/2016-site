package utils

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
	"unicode"
)

// This is non-exhaustive, but captures most required nodes
const stripAllowedNodes = "html body p a strong em span section div"

// StripHtml strips HTML tags from a string, extracting all plain text.
func StripHtml(htmls string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmls))
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	anodes := strings.Split(stripAllowedNodes, " ")

	var f func(n *html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.DocumentNode:
			if n.FirstChild != nil {
				f(n.FirstChild)
			}
			break
		case html.ElementNode:
			// Decide whether to descend into the element's
			// children
			allowed := false
			for _, anode := range anodes {
				if n.Data == anode {
					allowed = true
					break
				}
			}

			if allowed && n.FirstChild != nil {
				f(n.FirstChild)
			}

			// Ensure whitespace between paragraphs (hack!)
			if n.Data == "p" {
				buffer.WriteString("\n\n")
			}

			break
		case html.TextNode:
			buffer.WriteString(n.Data)
			break
		default:
			break
		}

		if n.NextSibling != nil {
			f(n.NextSibling)
		}
	}

	f(doc)

	return strings.TrimFunc(buffer.String(), unicode.IsSpace), nil
}
