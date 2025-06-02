package main

import (
	"log"
	"os"
	"repo_stats/services"
	"repo_stats/utils"
)

func main() {
	var err error

	envFile, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
		return
	}
	envData, err := utils.ReadEnv(envFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	repoUser := utils.GetInput("Repository Owner", utils.Title)
	repoName := utils.GetInput("Repository Name", utils.Title)

	// Make stuff
	api := services.NewGHAPI(repoUser, repoName, envData["GITHUB_TOKEN"])
	stats := utils.NewStats(repoUser, repoName, []string{".png", ".svg", ".jpg", ".lock"}, []string{"package-lock.json", "yarn.lock"})

	//	Get PRs
	prs, err := api.GetPRs()
	if err != nil {
		log.Fatal(err)
		return
	}
	stats.SetPRs(prs)

	// Get Commits
	commits, err := api.GetCommits()
	if err != nil {
		log.Fatal(err)
	}
	stats.SetCommits(commits)

	// Get data from the commits
	fileURLs, fileSizes, fileChanges, err := api.ExtractFileData(commits)
	if err != nil {
		log.Fatal(err)
		return
	}

	stats.SetFileUrls(fileURLs)
	stats.SetFileSizes(fileSizes)
	stats.SetFileChanges(fileChanges)

	stats.OutputResults()

	err = utils.OutputFrom([]string{"[Rate Limit]", api.GetRateLimitRemainingString()},
		[]utils.Color{utils.Subtle, utils.Highlight})
	if err != nil {
		log.Fatal(err)
		return
	}
}
