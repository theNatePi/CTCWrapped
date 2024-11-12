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
