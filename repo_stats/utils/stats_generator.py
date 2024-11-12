from collections import namedtuple
import requests
from datetime import datetime
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
