package utils

import (
	"fmt"
	"strings"
)

// PrefixUrl prefixes a URL url with prefix p.
// If url begins with a single slash, it is treated as relative and prefixed with prefix.
// If p does not begin with a slash and end in zero slashes, it is changed to do so.
func PrefixUrl(url, p string) string {
	relative := strings.HasPrefix(url, "/") && !strings.HasPrefix(url, "//")
	if !relative {
		return url
	}

	pt := strings.Trim(p, "/")
	// Prevent accidentally outputting //url if prefix is empty
	if len(pt) == 0 {
		return url
	}

	// Because url is relative, it must already start with a slash
	return fmt.Sprintf("/%s%s", pt, url)
}
