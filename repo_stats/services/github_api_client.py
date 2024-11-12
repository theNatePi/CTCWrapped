import os
from dotenv import load_dotenv
from pathlib import Path
import requests
from datetime import datetime

class GitHubAPIClient:
    """A client for the GitHub API."""
    def __init__(self, base_url: str, *, 
                 token: str | None = None, 
                 env_path: str | None = None):
        self.base_url = base_url
        self.token = token

        if self.gh_token is None:
            load_dotenv(Path(env_path or ".env"))
            self.gh_token = os.getenv("GITHUB_TOKEN")
        
        if self.gh_token is None:
            raise ValueError("No GitHub token provided\n"
                             "  - Pass a token as an argument\n"
                             "  - Set the GITHUB_TOKEN environment variable\n")

        self.headers = {"Authorization": f"token {self.gh_token}"}
    
    def get_rate_limit_remaining(self) -> int:
        """Get the rate limit remaining for the GitHub API."""
        response = requests.get("https://api.github.com/rate_limit", headers=self.headers)
        return response.json()["resources"]["core"]["remaining"]
    
    def get_rate_limit_used(self) -> int:
        """Get the rate limit used for the GitHub API."""
        response = requests.get("https://api.github.com/rate_limit", headers=self.headers)
        return response.json()["resources"]["core"]["used"]
    
    def get_rate_limit_reset_time(self) -> datetime:
        """Get the time at which the rate limit will reset."""
        response = requests.get("https://api.github.com/rate_limit", headers=self.headers)
        return datetime.fromtimestamp(response.json()["resources"]["core"]["reset"])
