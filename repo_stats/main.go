package main

import (
	"fmt"
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

	// Make stuff
	api := services.NewGHAPI("ctc-uci", "lpa", envData["GITHUB_TOKEN"])
	stats := utils.NewStats("ctc-uci", "lpa", []string{}, []string{})

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

	fmt.Println(stats.TopPRs(5))
	fmt.Println(stats.TopCommits(5))

	// Get data from the commits
	err = utils.OutputFrom([]string{"[Rate Limit]", api.GetRateLimitRemainingString()},
		[]utils.Color{utils.Subtle, utils.Highlight})
	if err != nil {
		log.Fatal(err)
		return
	}
}
