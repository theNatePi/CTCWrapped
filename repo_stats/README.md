## Repo Stats Generator
Generates stats about repositories, including:
- Number of lines of code
- Number of PRs
- Number of commits

As well as:
- Top (5) users with the most PRs
- Top (5) users with the most commits
- Top (5) files by number of lines
- Top (5) files by total changes

## Installation
1. Install Go https://go.dev/doc/install
2. Clone repo, navigate into `/CTCWrapped/repo_stats`
3. Run `go build`
4. Run `./repo_stats`

## Setup
Generate a GitHub Personal Access Token
1. Go to: https://github.com/settings/tokens
2. Generate new token (top right)
    - Generate a classic token
3. Click the "repo" scope
4. Generate token, add it into an `.env` file
    - Match the format of `.env.example`

## Generate Stats
1. Run `./repo_stats`
2. Provide GitHub username for the repository owner
3. Provide the repository name
4. Wait for all calls to complete, this may take a while

## Example Output
<img width="375" alt="image" src="https://github.com/user-attachments/assets/c811b50d-7e49-41ed-a7c4-92ecdd26f75d" />
