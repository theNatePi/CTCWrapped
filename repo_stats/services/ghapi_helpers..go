package services

import (
	"regexp"
	"strings"
)

func parseNextLinkRegex(linkHeader string) string {
	// Regex to match <URL>; rel="next"
	re := regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
	matches := re.FindStringSubmatch(linkHeader)

	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func numLines(s string) int {
	// https://stackoverflow.com/questions/27466545/how-do-i-get-the-number-of-lines-in-a-string
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}
