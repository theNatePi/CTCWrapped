package utils

import "sort"

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

func (x Stats) SetPRs(PRs []interface{}) {
	x.numPRs = len(PRs)
	x.allPRs = PRs
	for _, PR := range PRs {
		_user := PR.(map[string]interface{})["user"]
		_login := _user.(map[string]interface{})["login"]
		attribution := _login.(string)
		x.prAttribution[attribution]++
	}
}

func (x Stats) SetCommits(commits []interface{}) {
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

func (x Stats) TopPRs(n int) map[string]int {
	keys := make([]string, 0, len(x.prAttribution))
	for key := range x.prAttribution {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return x.prAttribution[keys[i]] > x.prAttribution[keys[j]]
	})

	if n > len(keys) {
		n = len(keys)
	}

	result := make(map[string]int)
	for i := 0; i < n && i < len(keys); i++ {
		key := keys[i]
		result[key] = x.prAttribution[key]
	}
	return result
}

func (x Stats) TopCommits(n int) map[string]int {
	keys := make([]string, 0, len(x.commitAttribution))
	for key := range x.commitAttribution {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return x.commitAttribution[keys[i]] > x.commitAttribution[keys[j]]
	})

	if n > len(keys) {
		n = len(keys)
	}

	result := make(map[string]int)
	for i := 0; i < n && i < len(keys); i++ {
		key := keys[i]
		result[key] = x.commitAttribution[key]
	}
	return result
}
