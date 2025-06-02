package services

import (
	"fmt"
	"net/http"
	"repo_stats/utils"
	"strconv"
	"time"
)

type GHAPI struct {
	RepoOwner          string
	RepoName           string
	RequestCategory    string
	Verbose            bool
	rateLimitRemaining int
	rateLimitReset     time.Time
	authToken          string
	stats              utils.Stats
}

func NewGHAPI(repoOwner, repoName string, authToken string) *GHAPI {
	api := &GHAPI{RepoOwner: repoOwner, RepoName: repoName, RequestCategory: "", Verbose: true,
		rateLimitRemaining: 5000, authToken: authToken}
	api.RequestCategory = "Rate Limit"
	// Make any call to set the rate limit given the response header
	api.makeRequest("https://api.github.com/rate_limit", "")
	api.RequestCategory = ""
	return api
}

// makeRequest
//
// Make a request to the GitHub API. Used to force update of rate-limit
func (x *GHAPI) makeRequest(url string, body string) (string, http.Header, error) {
	if x.RequestCategory != "" && x.Verbose {
		utils.OutputFrom([]string{"[" + x.GetRateLimitRemainingString() + "]",
			x.RequestCategory, url},
			[]utils.Color{utils.Highlight, utils.TitleNoBold, utils.Subtle})
	}

	if x.rateLimitRemaining == 0 {
		// If we have hit the rate limit, wait for the rest
		timeToWait := x.rateLimitReset.Sub(time.Now())
		utils.OutputFrom([]string{"Rate Limit Hit:", "Waiting For",
			strconv.FormatInt(int64(timeToWait), 10)},
			[]utils.Color{utils.Err, utils.Subtle, utils.Subtle})
		time.Sleep(timeToWait * time.Second)
	}

	respBody, header, err := utils.Get(url, body,
		map[string]string{
			"Accept":        "application/vnd.github+json",
			"Authorization": "Bearer " + x.authToken,
		})
	if err != nil {
		return "", nil, err
	}
	convertedLimit, err := strconv.Atoi(header.Get("X-RateLimit-Remaining"))
	if err == nil {
		x.rateLimitRemaining = convertedLimit
	}
	if x.rateLimitRemaining == 0 {
		convertedReset, err := strconv.ParseInt(header.Get("X-RateLimit-Reset"), 10, 64)
		if err != nil {
			return "", nil, err
		}
		x.rateLimitReset = time.Unix(convertedReset, 0)
	}
	return respBody, header, nil
}

func (x *GHAPI) GetRateLimitRemaining() int {
	return x.rateLimitRemaining
}

func (x *GHAPI) GetRateLimitRemainingString() string {
	return strconv.Itoa(x.rateLimitRemaining)
}

func (x *GHAPI) GetRateLimitReset() time.Time {
	return x.rateLimitReset
}

func (x *GHAPI) GetPRs() ([]interface{}, error) {
	x.RequestCategory = "Pull Requests"
	formattedUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=all",
		x.RepoOwner, x.RepoName)
	body, _, err := x.makeRequest(formattedUrl, "")
	if err != nil {
		return nil, err
	}
	parsedBody, err := utils.ParseBody(body)
	return parsedBody.([]interface{}), nil
}

func (x *GHAPI) GetCommits() ([]interface{}, error) {
	x.RequestCategory = "Commits"
	// Make initial call
	formattedUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits",
		x.RepoOwner, x.RepoName)
	body, headers, err := x.makeRequest(formattedUrl, "")
	if err != nil {
		return nil, err
	}
	parsedBody, err := utils.ParseBody(body)
	if err != nil {
		return nil, err
	}
	commits := parsedBody.([]interface{})

	//	While there are more pages, explore them
	nextLink := parseNextLinkRegex(headers.Get("Link"))
	for nextLink != "" {
		body, headers, err = x.makeRequest(nextLink, "")
		if err != nil {
			return nil, err
		}
		parsedBody, err = utils.ParseBody(body)
		if err != nil {
			return nil, err
		}
		commits = append(commits, parsedBody.([]interface{})...)
		nextLink = parseNextLinkRegex(headers.Get("Link"))
	}
	return commits, nil
}

func (x *GHAPI) getCommitData(commit interface{}) (interface{}, error) {
	x.RequestCategory = "Individual Commit"
	commitURL := commit.(map[string]interface{})["url"].(string)
	body, _, err := x.makeRequest(commitURL, "")
	if err != nil {
		return nil, err
	}
	parsedBody, err := utils.ParseBody(body)
	if err != nil {
		return nil, err
	}
	return parsedBody, nil
}

func (x *GHAPI) getFileSize(fileURL string) (int, error) {
	x.RequestCategory = "File"
	file, _, err := x.makeRequest(fileURL, "")
	if err != nil {
		return -1, err
	}

	lineCount := numLines(file)
	return lineCount, nil
}

func (x *GHAPI) ExtractFileData(commits []interface{}) (map[string]string, map[string]int, map[string]int, error) {
	fileURLMap := make(map[string]string)
	fileSizeMap := make(map[string]int)
	fileChangesMap := make(map[string]int)

	for _, commit := range commits {
		commitData, err := x.getCommitData(commit)
		if err != nil {
			return nil, nil, nil, err
		}
		_files := commitData.(map[string]interface{})["files"].([]interface{})

		for _, file := range _files {
			_filename := file.(map[string]interface{})["filename"].(string)
			_fileURL := file.(map[string]interface{})["raw_url"].(string)
			fileURLMap[_filename] = _fileURL

			_fileChanges := file.(map[string]interface{})["changes"].(float64)
			fileChangesMap[_filename] += int(_fileChanges)
		}
	}

	for file, url := range fileURLMap {
		fileSize, err := x.getFileSize(url)
		if err != nil {
			return nil, nil, nil, err
		}
		fileSizeMap[file] = fileSize
	}

	return fileURLMap, fileSizeMap, fileChangesMap, nil
}
