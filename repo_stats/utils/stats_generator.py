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

        if self.api_client.get_rate_limit_remaining() <= 100:
            reset_time = self.api_client.get_rate_limit_reset_time()
            self.io_handler.output("Rate Limit Reached",
                                  f"Waiting {(reset_time - datetime.now()).total_seconds()}", 
                                   "for rate limit reset", "error")
            time.sleep((reset_time - datetime.now()).total_seconds())
            return self._make_request(link)

        response = requests.get(link, headers = self.headers)

        if response.status_code == 200:
            # Parse the JSON response and print it
            data = response.json()
            self._set_rate_limit(response)
            return RequestResult(response.status_code, response, data)
        else:
            self.io_handler.error("Request Failed",
                                  "Status code: {response.status_code}")
            return RequestResult(response.status_code, response, {})
