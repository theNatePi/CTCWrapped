from termcolor import colored

class RepoStats:
    """Stores statistics for a GitHub repository."""
    def __init__(self):
        self.num_pulls = 0
        self.top_pulls = []
        self.num_commits = 0
        self.top_commits = []
        self.num_merges = 0
        self.top_files_changed = []
        self.largest_files = []
        self.total_lines_of_code = 0

    @staticmethod
    def get_top_n(items_dict: dict[str, int], n_amount: int) -> list[tuple[str, int]]:
        """Get the top n items by value from a dictionary."""
        items = [(k, v) for k, v in items_dict.items()]
        return sorted(items, key=lambda x: x[1], reverse=True)[:n_amount]

    def output_results(self):
        """Outputs formatted statistics about the repository to the console.
        
        Prints a summary including:
        - Total number of pulls, commits, merges and lines of code
        - Top 5 contributors by number of pull requests
        - Top 5 contributors by number of commits 
        - Top 5 most frequently changed files
        - Top 5 largest files by line count
        - Remaining GitHub API rate limit
        
        Uses colored output for better readability.
        """
        def _output_top_five(title: str, top_five: list) -> None:
            print("\n" + colored(title, "light_magenta", attrs = ["bold"]))
            for item_index in range(len(top_five)):
                index = str(item_index + 1)
                item = top_five[item_index][0]
                number = str(top_five[item_index][1])
                print("  " + colored((index + ". " + item + " " + number), attrs = ["dark"]))

        print(colored(("#" + "-"*20 + "#"), "green", attrs=["bold"]))
        print(colored("Num PRs: ", "light_blue", attrs=["bold"]) + colored(self.num_pulls, "light_blue", ))
        print(colored("Num commits: ", "light_blue", attrs = ["bold"]) + colored(self.num_commits,"light_blue", ))
        print(colored("Num merges: ", "light_blue", attrs = ["bold"]) + colored(self.num_merges, "light_blue", ))
        print(colored("Total lines of code: ", "light_blue", attrs = ["bold"]) + colored(self.total_lines_of_code, "light_blue", ))

        _output_top_five("Most PRs", self.top_pulls)
        _output_top_five("Most commits", self.top_commits)
        _output_top_five("Most files changed", self.top_files_changed)
        _output_top_five("Largest files", self.largest_files)

        print()
        # print(colored(f"Requests remaining: {self.rate_limit_remaining}", "red"))

        print(colored(("#" + "-" * 20 + "#"), "green", attrs = ["bold"]))
        print()
