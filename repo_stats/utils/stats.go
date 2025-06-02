package utils

import (
	"fmt"
	"strconv"
)

type Stats struct {
	RepoUser          string
	RepoName          string
	numPRs            int
	allPRs            []interface{}
	allCommits        []interface{}
	numCommits        int
	totalLinesOfCode  int
	prAttribution     map[string]int
	commitAttribution map[string]int
	fileURls          map[string]string
	fileChanges       map[string]int
	fileSizes         map[string]int
	ignoreExtensions  []string
	ignoreFiles       []string
}

func NewStats(repoUser string, repoName string,
	ignoreExtensions []string, ignoreFiles []string) *Stats {
	return &Stats{RepoUser: repoUser, RepoName: repoName,
		numPRs: 0, numCommits: 0, totalLinesOfCode: 0,
		allPRs: make([]interface{}, 0), allCommits: make([]interface{}, 0),
		prAttribution: make(map[string]int), commitAttribution: make(map[string]int),
		fileURls: make(map[string]string), fileChanges: make(map[string]int), fileSizes: make(map[string]int),
		ignoreExtensions: ignoreExtensions, ignoreFiles: ignoreFiles}
}

func (x *Stats) SetPRs(PRs []interface{}) {
	x.numPRs = len(PRs)
	x.allPRs = PRs
	for _, PR := range PRs {
		_user := PR.(map[string]interface{})["user"]
		_login := _user.(map[string]interface{})["login"]
		attribution := _login.(string)
		if attribution != "dependabot[bot]" {
			x.prAttribution[attribution]++
		}
	}
}

func (x *Stats) SetCommits(commits []interface{}) {
	x.numCommits = len(commits)
	x.allCommits = commits
	for _, commit := range commits {
		_commit := commit.(map[string]interface{})["commit"]
		_committer := _commit.(map[string]interface{})["committer"]
		_name := _committer.(map[string]interface{})["name"]
		if _name != "GitHub" {
			attribution := _name.(string)
			x.commitAttribution[attribution]++
		}
	}
}

func (x *Stats) SetFileUrls(fileURLs map[string]string) {
	x.fileURls = fileURLs
}

func (x *Stats) SetFileSizes(fileSizes map[string]int) {
	x.fileSizes = fileSizes
	for _, fileSize := range x.fileSizes {
		x.totalLinesOfCode += fileSize
	}
}

func (x *Stats) SetFileChanges(fileChanges map[string]int) {
	x.fileChanges = fileChanges
}

func (x *Stats) TopPRs(n int) map[string]int {
	result := topnMapStrInt(x.prAttribution, n)
	return result
}

func (x *Stats) TopCommits(n int) map[string]int {
	result := topnMapStrInt(x.commitAttribution, n)
	return result
}

func (x *Stats) TopFileSizes(n int) map[string]int {
	result := x.filterFiles(x.fileSizes)
	result = topnMapStrInt(result, n)
	return result
}

func (x *Stats) TopFileChanges(n int) map[string]int {
	result := x.filterFiles(x.fileChanges)
	result = topnMapStrInt(result, n)
	return result
}

func (x *Stats) TotalLinesOfCode() int {
	return x.totalLinesOfCode
}

func (x *Stats) OutputResults() {
	fmt.Println("\n")

	OutputWithTitle("Stats For:", Title,
		x.RepoUser+"/"+x.RepoName, Subtle)
	fmt.Println()

	OutputFrom([]string{"Lines of code:", strconv.Itoa(x.totalLinesOfCode)},
		[]Color{TitleNoBold, Subtle})
	OutputFrom([]string{"Total commits:", strconv.Itoa(x.numCommits)},
		[]Color{TitleNoBold, Subtle})
	OutputFrom([]string{"Total PRs:", strconv.Itoa(x.numPRs)},
		[]Color{TitleNoBold, Subtle})
	fmt.Println()

	Output("Top PRs:", TitleNoBold)
	printTop(x.TopPRs(5))
	Output("Top Commits:", TitleNoBold)
	printTop(x.TopCommits(5))
	Output("Top File Sizes (lines of code):", TitleNoBold)
	printTop(x.TopFileSizes(5))
	Output("Top File Changes:", TitleNoBold)
	printTop(x.TopFileChanges(5))
}
