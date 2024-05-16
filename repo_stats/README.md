## Generating Stats from Repos

#### Steps

1. Clone the repository
2. In `/CTCWrapped`
   1. [Optional] Create VENV `python -m venv .venv`  (then activate it)
   2. Install requirements `pip install -r requirements.txt`
3. Navigate to `/repo_stats`
4. Create a `.env` file
   1. Create a **Classic** GitHub user token with `public_repo` permissions [(more info)](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)
   2. Inside the `.env` file, input `GITHUB_TOKEN = "YOUR TOKEN"`
5. Run generate_stats.py with `python generate_stats.py`
   1. **Note:** due to limitations with the GitHub API, all changes on the repo must be merged into `main` before stats can be gathered

