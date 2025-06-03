package services

import (
	"regexp"
	"strings"
)

// parseNextLinkRegex
// Parses linkHeader sent back from GitHub API to get the next link for paginated results
//
// Parameters:
//   - linkHeader: ex: <https://api.github.com/repositories/#/pulls?state=all&page=2>; rel="next",
//     <https://api.github.com/repositories/#/pulls?state=all&page=4>; rel="last"
//
// Returns the link for "next", ex: "https://api.github.com/repositories/#/pulls?state=all&page=2"
func parseNextLinkRegex(linkHeader string) string {
	// Regex to match <URL>; rel="next"
	re := regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
	matches := re.FindStringSubmatch(linkHeader)

	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// numLines
// Gets the number of lines in a given s
func numLines(s string) int {
	// Sourced from
	// https://stackoverflow.com/questions/27466545/how-do-i-get-the-number-of-lines-in-a-string
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}
