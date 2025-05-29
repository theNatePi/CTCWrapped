package utils

type Stats struct {
	repoUser          string
	repoName          string
	numPRs            int
	numCommits        int
	numMerges         int
	totalLinesOfCode  int
	prAttribution     map[string]int
	commitAttribution map[string]int
	fileChanges       map[string]int
	fileSizes         map[string]int
	ignoreExtensions  []string
	ignoreFiles       []string
}
