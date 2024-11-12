import os
import requests
import time
from dotenv import load_dotenv
from pathlib import Path
from datetime import datetime
from repo_stats.utils.input_output import InputOutputHandler, PrintWrapper

class GitHubAPIClient:
    """A client for the GitHub API."""
    def __init__(self, base_url: str, *, 
                 token: str | None = None, 
                 env_path: str | None = None):
        self._rate_limit_refresh_interval_default = 20
        self._rate_limit_refresh_interval = self._rate_limit_refresh_interval_default

        self.base_url = base_url
        self.token = token
        self.rate_limit_refresh = time.time()
        self.rate_limit_remaining = self.get_rate_limit_remaining()

        if self.gh_token is None:
            load_dotenv(Path(env_path or ".env"))
            self.gh_token = os.getenv("GITHUB_TOKEN")
        
        if self.gh_token is None:
            raise ValueError("No GitHub token provided\n"
                             "  - Pass a token as an argument\n"
                             "  - Set the GITHUB_TOKEN environment variable\n")

        self.headers = {"Authorization": f"token {self.gh_token}"}

    def _set_rate_limit(self, rate_limit: int) -> None:
        """Update the current rate limit remaining from a response object."""
        self.rate_limit_remaining = rate_limit
        if rate_limit <= 200:
            self._rate_limit_refresh_interval = 1
        else:
            self._rate_limit_refresh_interval = self._rate_limit_refresh_interval_default

        self.rate_limit_refresh = time.time() + self._rate_limit_refresh_interval

    def get_rate_limit_remaining(self) -> int:
        """Get the rate limit remaining for the GitHub API."""
        if time.time() > self.rate_limit_refresh:
            response = requests.get("https://api.github.com/rate_limit", headers=self.headers)
            self._set_rate_limit(response.json()["resources"]["core"]["remaining"])

        return self.rate_limit_remaining
    
    def get_rate_limit_used(self) -> int:
        """Get the rate limit used for the GitHub API."""
        response = requests.get("https://api.github.com/rate_limit", headers=self.headers)
        return response.json()["resources"]["core"]["used"]
    
    def get_rate_limit_reset_time(self) -> datetime:
        """Get the time at which the rate limit will reset."""
        response = requests.get("https://api.github.com/rate_limit", headers=self.headers)
        return datetime.fromtimestamp(response.json()["resources"]["core"]["reset"])
    
    def make_request(self, link: str, 
                     err_ouput_handler: InputOutputHandler = PrintWrapper()
                    ) -> requests.Response:
        """Make a request to a given link and return the response code, response object and data."""
        if self.get_rate_limit_remaining() <= 100:
            reset_time = self.get_rate_limit_reset_time()
            err_ouput_handler.output("Rate Limit Reached",
                                    f"Waiting {(reset_time - datetime.now()).total_seconds()}", 
                                     "for rate limit reset", "error")
            time.sleep((reset_time - datetime.now()).total_seconds())
            return self.make_request(link, err_ouput_handler)

        response = requests.get(link, headers=self.headers)
        self._set_rate_limit(int(dict(response.headers)["X-RateLimit-Remaining"]))
        return response
