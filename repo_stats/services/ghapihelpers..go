package services

import "regexp"

func parseNextLinkRegex(linkHeader string) string {
	// Regex to match <URL>; rel="next"
	re := regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
	matches := re.FindStringSubmatch(linkHeader)

	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}
