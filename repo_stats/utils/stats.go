package utils

import (
	"fmt"
	"strconv"
)

// Stats
// Represents all statistics pulled from a repo
type Stats struct {
	// The owner of the repo
	RepoUser string
	// The name of the repo
	RepoName string
	numPRs   int
	// An array of all PRs in the repo, initially empty
	allPRs     []interface{}
	numCommits int
	// An array of all commits in teh repo, initially empty
	allCommits       []interface{}
	totalLinesOfCode int
	// A map of GitHub username to number of PRs authored
	prAttribution map[string]int
	// A map of GitHub username to number of commits authored
	commitAttribution map[string]int
	// A map of file path to file api url
	fileURLs map[string]string
	// A map of file path to number of line changes (insertion + deletion) total
	fileChanges map[string]int
	// A map of file path to number of lines in the file
	fileSizes map[string]int
	// An array of file extensions (.png, .svg, .jpg, etc) to ignore
	ignoreExtensions []string
	// An array of file names (yarn.lock, package-lock.json, etc) to ignore
	ignoreFiles []string
	// An array of directories to ignore
	ignoreDirs []string
}

// NewStats
// Returns a new Stats struct with proper attributes
//
// Parameters:
//   - repoUser: username for GitHub repository
//   - repoName: the name of the GitHub repository
//   - ignoreExtensions: an array of file extensions to ignore (including "." before each)
//   - ignoreFiles: an array of file names to ignore
//
// Returns pointer to new Stats struct
func NewStats(repoUser string, repoName string,
	ignoreExtensions []string, ignoreFiles []string,
	ignoreDirs []string) *Stats {
	return &Stats{RepoUser: repoUser, RepoName: repoName,
		numPRs: 0, numCommits: 0, totalLinesOfCode: 0,
		allPRs: make([]interface{}, 0), allCommits: make([]interface{}, 0),
		prAttribution: make(map[string]int), commitAttribution: make(map[string]int),
		fileURLs: make(map[string]string), fileChanges: make(map[string]int), fileSizes: make(map[string]int),
		ignoreExtensions: ignoreExtensions, ignoreFiles: ignoreFiles, ignoreDirs: ignoreDirs}
}

// SetPRs
// Sets x.allPRs, x.numPRs, and x.prAttribution(s)
//
// Parameters:
//   - PRs: array of PRs in format returned from GitHub API
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

// SetCommits
// Sets x.numCommits, x.allCommits, and x.commitAttribution(s)
//
// Parameters:
//   - commits: array of commits in format returned from GitHub API
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

// SetFileUrls
// Sets the local fileURLs to fileURLs
func (x *Stats) SetFileUrls(fileURLs map[string]string) {
	x.fileURLs = fileURLs
}

// SetFileSizes
// Sets the local fileSizes to fileSizes
// Updates total lines of code for valid files
func (x *Stats) SetFileSizes(fileSizes map[string]int) {
	filteredFiles := x.filterFiles(fileSizes)

	x.fileSizes = filteredFiles
	for _, fileSize := range filteredFiles {
		x.totalLinesOfCode += fileSize
	}
}

// SetFileChanges
// Sets the local fileChanges to fileChanges, only includes valid files
func (x *Stats) SetFileChanges(fileChanges map[string]int) {
	filteredFiles := x.filterFiles(fileChanges)
	x.fileChanges = filteredFiles
}

// TopPRs
// Gets the top n PRs (in order)
func (x *Stats) TopPRs(n int) map[string]int {
	result := topnMapStrInt(x.prAttribution, n)
	return result
}

// TopCommits
// Gets the top n commits (in order)
func (x *Stats) TopCommits(n int) map[string]int {
	result := topnMapStrInt(x.commitAttribution, n)
	return result
}

// TopFileSizes
// Gets the top n files by size (in order)
func (x *Stats) TopFileSizes(n int) map[string]int {
	result := x.filterFiles(x.fileSizes)
	result = topnMapStrInt(result, n)
	return result
}

// TopFileChanges
// Gets the top n file changes (in order)
func (x *Stats) TopFileChanges(n int) map[string]int {
	result := x.filterFiles(x.fileChanges)
	result = topnMapStrInt(result, n)
	return result
}

// TotalLinesOfCode
// Gets total lines of code
func (x *Stats) TotalLinesOfCode() int {
	return x.totalLinesOfCode
}

// OutputResults
// Outputs results of a statistics collection
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
