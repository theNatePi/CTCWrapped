from collections import namedtuple
from datetime import datetime
import requests
import base64
import time
from repo_stats.services.github_api_client import GitHubAPIClient
from repo_stats.utils.stats import RepoStats
from repo_stats.utils.input_output import InputOutputHandler


RequestResult = namedtuple("RequestResult", ["status_code", "response", "data"])

class StatsGenerator:
    """Generate statistics for a GitHub repository."""
    def __init__(self, repo_name: str, user_name: str, 
                 api_client: GitHubAPIClient = None, stats: RepoStats = None, 
                 input_output_handler: InputOutputHandler = None, verbose: bool = False):
        self.repo_name = repo_name
        self.user_name = user_name
        self.api_client = api_client or GitHubAPIClient(
            base_url=f"https://api.github.com/repos/{user_name}/{repo_name}"
        )
        self.stats = stats or RepoStats()
        self.io_handler = input_output_handler or InputOutputHandler(
            verbose=verbose
        )
        self.verbose = verbose
    
    def _make_request(self, link: str) -> RequestResult:
        """Make a request to a given link and return the response code, response object and data."""
        self.io_handler.output("Making request", link, "progress")

        response = self.api_client.make_request(link, self.io_handler)

        if response.status_code == 200:
            # Parse the JSON response and print it
            data = response.json()
            return RequestResult(response.status_code, response, data)
        else:
            self.io_handler.error("Request Failed",
                                  "Status code: {response.status_code}")
            return RequestResult(response.status_code, response, {})
    
    def _get_all_events(self, link: str) -> list:
        """Get all events from a link where there may be multiple pages."""
        events = []
        response_code, response, data = self._make_request(link)
        if response_code != 200:
            return []

        events.extend(data)
        try:
            events.extend(self._get_all_events(response.links["next"]["url"]))
            return events
        except KeyError:
            # No more pages
            return events
    
    def _get_file_changes(self, commit_link: str, files_dict: dict) -> None:
        """Get the file changes from a commit link."""
        response_code, _, commit = self._make_request(commit_link)
        if response_code != 200:
            return

        files = commit["files"]
        for file in files:
            filename = file["filename"]
            changes_stats = int(file["changes"])

            if filename not in files_dict:
                files_dict[filename] = changes_stats
            else:
                files_dict[filename] += changes_stats
    
    def _get_file_contents(self, file_url: str) -> int:
        """Get the contents of a file from a file link."""
        response_code, _, file = self._make_request(file_url)
        if response_code != 200:
            return

        content = file["content"]
        decoded_content = base64.b64decode(content).splitlines()
        return len(decoded_content)
    
    def _get_file_lines(self, contents_link, lines_dict):
        """Get the lines of code for all files in a directory from a contents link."""
        response_code, _, contents = self._make_request(contents_link)
        if response_code != 200:
            return

        for file in contents:
            filename = file["name"]
            type = file["type"]

            if type == "dir":
                self._get_file_lines(file["url"], lines_dict)
                continue

            if type == "file" and \
                    (filename.split(".")[-1] in ["js", "jsx", "css", "html", "ts", "tsx", "md"]):
                lines = self._get_file_contents(file["url"])
                lines_dict[filename] = lines
    
    def get_pulls(self):
        """
        Retrieves all pull requests from the repository and calculates statistics.
        
        Gets all pull requests using the GitHub API, counts the total number of pulls,
        and tracks the number of pull requests per user. Updates instance variables:
        - self.num_pulls: Total number of pull requests
        - self.top_pulls: Top 5 users by number of pull requests created
        """
        self.api_client.request_type = "Pulls"
        self.io_handler.output("Getting all pulls...",
                               "100 pulls per request",
                               "success")
        pulls = self._get_all_events(f"{self.api_client.base_url}/pulls?state=all&page=1&per_page=100")
        self.stats.num_pulls = len(pulls)

        pulls_by_person = {}
        for pull in pulls:
            pulls_by_person[pull["user"]["login"]] = \
                pulls_by_person.get(pull["user"]["login"], 0) + 1
        self.stats.top_pulls = self.stats.get_top_n(pulls_by_person, 5)
        self.io_handler.output("Pulls retrieved",
                              f"{self.api_client.request_counter} requests made",
                               "success")

    def get_commits(self):
        """
        Retrieves all commits from the repository and calculates statistics.
        
        Gets all commits using the GitHub API, counts the total number of commits,
        tracks commits per user, and analyzes merge commits and file changes.
        Updates instance variables:
        - self.num_commits: Total number of commits
        - self.top_commits: Top 5 users by number of commits authored
        - self.num_merges: Number of merge commits (commits with multiple parents)
        - self.top_files_changed: Top 5 files by number of times modified in commits
        """
        self.api_client.request_type = "Commits"
        self.io_handler.output("Getting all commits...",
                             "100 commits per request",
                             "success")
        commits = self._get_all_events(f"{self.api_client.base_url}/commits?page=1&per_page=100")
        self.stats.num_commits = len(commits)

        commits_per_person = {}
        failed_commits = 0
        for commit in commits:
            try:
                commits_per_person[commit["author"]["login"]] = commits_per_person.get(
                    commit["author"]["login"], 0) + 1
            except TypeError:
                failed_commits += 1
        self.top_commits = self.get_top_n(commits_per_person, 5)
        self.io_handler.output("Commits retrieved",
                              f"{self.api_client.request_counter} requests made",
                               "success")

        self.api_client.request_type = "Merges and File Changes"
        self.io_handler.output("Getting all merges and file changes...",
                               "1 request per file in each commit",
                              "success")
        merges = 0
        files_changed = {}
        for commit in commits:
            parents = commit["parents"]
            if len(parents) > 1:
                merges += 1
            commit_url = commit["url"]
            self._get_file_changes(commit_url, files_changed)

        self.stats.num_merges = merges
        self.stats.top_files_changed = self.stats.get_top_n(files_changed, 5)
        self.io_handler.output("Merges and file changes retrieved",
                              f"{self.api_client.request_counter} requests made",
                               "success")

    def get_files(self):
        """Gets all files in the repository using the GitHub API and counts lines of code.
        
        Makes requests to get contents of all files and directories recursively.
        Updates instance variables:
        - self.largest_files: Top 5 files by number of lines of code
        - self.total_lines_of_code: Total lines of code across all files
        """
        self.api_client.request_type = "Files"
        self.io_handler.output("Getting all files...",
                               "1 request per file in repository",
                             "success")
        line_counts = {}
        self._get_file_lines(f"{self.api_client.base_url}/contents", line_counts)
        self.stats.largest_files = self.stats.get_top_n(line_counts, 5)
        self.stats.total_lines_of_code = sum(line_counts.values())
        self.io_handler.output("Contents retrieved, lines of code retrieved",
                              f"{self.api_client.request_counter} requests made",
                               "success")
