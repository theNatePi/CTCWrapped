import os
import requests
import base64
from dotenv import load_dotenv
from pathlib import Path
from termcolor import colored

class statsGenerator:
    def __init__(self, base_url: str) -> None:
        load_dotenv(Path(".env"))
        self.gh_token = os.getenv("GITHUB_TOKEN")
        self.headers = {"Authorization": f"token {self.gh_token}"}

        # Default rate limit for API with auth
        self.rate_limit_remaining = 5000
        self._request_counter = 0
        self._current_request_type = ""

        self.base_url = base_url

        self.num_pulls = 0
        self.top_pulls = []
        self.num_commits = 0
        self.top_commits = []
        self.num_merges = 0
        self.top_files_changed = []
        self.largest_files = []
        self.total_lines_of_code = 0

    def _output_message(self, title: str, message: str, message_type: str) -> None:
        color = 'white'
        if message_type == "success":
            color = 'green'
        elif message_type == "progress":
            color = 'light_blue'
        elif message_type == "error":
            color = 'red'

        open_bracket = colored('[', attrs = ['dark'])
        close_bracket = colored(']', attrs = ['dark'])
        print(f"{open_bracket}"
              f"{colored(self._current_request_type, 'light_green')}{close_bracket} "
              f"{open_bracket}{colored(self.rate_limit_remaining, 'light_magenta', attrs = ['dark'])}"
              f"{close_bracket} {colored(title, color)}: "
              f"{colored(message, attrs = ['dark'])}")

    def _set_rate_limit(self, response: requests.Response) -> None:
        self.rate_limit_remaining = int(dict(response.headers)["X-RateLimit-Remaining"])

    def _make_request(self, link):
        self._output_message("Making request", link, "progress")
        self._request_counter += 1

        if self.rate_limit_remaining <= 100:
            self._output_message("Rate Limit Reached",
                                 "Please wait for a while before making more requests.",
                                 "error")
            return 403, None, {}

        response = requests.get(link, headers = self.headers)

        if response.status_code == 200:
            # Parse the JSON response and print it
            data = response.json()
            self._set_rate_limit(response)
            return response.status_code, response, data
        else:
            self._output_message("Request Failed",
                                 f"Status code: {response.status_code}",
                                 "error")
            return response.status_code, response, {}

    def _get_all_events(self, link) -> list:
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

    def _get_file_changes(self, commit_link, files_dict):
        response_code, response, commit = self._make_request(commit_link)
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

    def _get_file_contents(self, file_url):
        response_code, response, file = self._make_request(file_url)
        if response_code != 200:
            return

        content = file["content"]
        decoded_content = base64.b64decode(content).splitlines()
        return len(decoded_content)

    def _get_file_lines(self, contents_link, lines_dict):
        response_code, response, contents = self._make_request(contents_link)
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

    def get_top_five(self, files_dict):
        files = []
        for file, value in files_dict.items():
            files.append((file, value))

        return sorted(files, key = lambda x: x[1], reverse = True)[:5]

    def get_pulls(self):
        self._request_counter = 0
        self._current_request_type = "Pulls"
        self._output_message("Getting all pulls...",
                             "100 pulls per request",
                             "success")
        pulls = self._get_all_events(f"{self.base_url}/pulls?state=all&page=1&per_page=100")
        self.num_pulls = len(pulls)

        pulls_by_person = {}
        for pull in pulls:
            pulls_by_person[pull["user"]["login"]] = \
                pulls_by_person.get(pull["user"]["login"], 0) + 1
        self.top_pulls = self.get_top_five(pulls_by_person)
        self._output_message("Pulls retrieved",
                             f"{self._request_counter} requests made",
                             "success")

    def get_commits(self):
        self._request_counter = 0
        self._current_request_type = "Commits"
        self._output_message("Getting all commits...",
                             "100 commits per request",
                             "success")
        commits = self._get_all_events(f"{self.base_url}/commits?page=1&per_page=100")
        self.num_commits = len(commits)

        commits_per_person = {}
        failed_commits = 0
        for commit in commits:
            try:
                commits_per_person[commit["author"]["login"]] = commits_per_person.get(
                    commit["author"]["login"], 0) + 1
            except TypeError:
                failed_commits += 1
        self.top_commits = self.get_top_five(commits_per_person)
        self._output_message("Commits retrieved",
                             f"{self._request_counter} requests made",
                             "success")

        self._request_counter = 0
        self._current_request_type = "Merges and File Changes"
        self._output_message("Getting all merges and file changes...",
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

        self.num_merges = merges
        self.top_files_changed = self.get_top_five(files_changed)
        self._output_message("Merges and file changes retrieved",
                             f"{self._request_counter} requests made",
                             "success")

    def get_files(self):
        self._request_counter = 0
        self._current_request_type = "Files"
        self._output_message("Getting all files...",
                             "1 request per file in repository",
                             "success")
        line_counts = {}
        self._get_file_lines(f"{self.base_url}/contents", line_counts)
        self.largest_files = self.get_top_five(line_counts)
        self.total_lines_of_code = sum(line_counts.values())
        self._output_message("Contents retrieved, lines of code retrieved",
                             f"{self._request_counter} requests made",
                             "success")

    def output_results(self):
        def _output_top_five(title: str, top_five: list) -> None:
            print("\n" + colored(title, "light_magenta", attrs = ["bold"]))
            for item_index in range(len(top_five)):
                index = str(item_index + 1)
                item = top_five[item_index][0]
                number = str(top_five[item_index][1])
                print("  " + colored((index + ". " + item + " " + number), attrs = ["dark"]))

        print(colored(("#" + "-"*20 + "#"), "green", attrs=["bold"]))
        print(colored("Num pulls: ", "light_blue", attrs=["bold"]) + colored(self.num_pulls, "light_blue", ))
        print(colored("Num commits: ", "light_blue", attrs = ["bold"]) + colored(self.num_commits,"light_blue", ))
        print(colored("Num merges: ", "light_blue", attrs = ["bold"]) + colored(self.num_merges, "light_blue", ))
        print(colored("Total lines of code: ", "light_blue", attrs = ["bold"]) + colored(self.total_lines_of_code, "light_blue", ))

        _output_top_five("Most pulls", self.top_pulls)
        _output_top_five("Most commits", self.top_commits)
        _output_top_five("Most files changed", self.top_files_changed)
        _output_top_five("Largest files", self.largest_files)

        print()
        print(colored(f"Requests remaining: {self.rate_limit_remaining}", "red"))

        print(colored(("#" + "-" * 20 + "#"), "green", attrs = ["bold"]))
        print()


if __name__ == "__main__":
    user = input("Enter GitHub Username: ")
    repo = input("Enter GitHub Repository: ")
    base_url = f"https://api.github.com/repos/{user}/{repo}"

    stats_generator = statsGenerator(base_url)

    stats_generator.get_pulls()
    stats_generator.output_results()
    stats_generator.get_commits()
    stats_generator.get_files()

    print("Rate limit remaining:", stats_generator.rate_limit_remaining)
