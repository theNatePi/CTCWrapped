## Generating Stats from Repos

#### Example output
When run, the repo stats will display something like the following:

<img width="334" alt="Screenshot 2025-01-30 at 1 51 46 PM" src="https://github.com/user-attachments/assets/ab1d91d0-73df-45a3-adf2-dced84ea5b59" />

#### Steps

1. Clone the repository
2. In `/CTCWrapped`
   1. [Optional] Create VENV `python -m venv .venv`  (then activate it)
   2. Install requirements `pip install -r requirements.txt`
3. Navigate to `/repo_stats`
4. Create a `.env` file
   1. Create a **Classic** GitHub user token with `public_repo` permissions [(more info)](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)
   2. Inside the `.env` file, input `GITHUB_TOKEN = "YOUR TOKEN"`
5. Run main.py with `python main.py`
   1. **Note:** due to limitations with the GitHub API, all changes on the repo must be merged into `main` on the respective repo before stats can be gathered
   2. **Note:** if using the `public_repo` token permissions, the repository must be public
  
#### Advanced Usage
You can run the script from the command line as well: <br>
`python main.py [-u USERNAME] [-r REPO] [-v VERBOSE]` <br>
Where `VERBOSE` is one of `yes` or `no`

##### Verbose Output
The verbose output will be in the following format <br>
`[<SEARCH CATEGORY>] [<REAMAINING API CALLS>] Making request: <URL>`

See `python main.py -h` for more.
