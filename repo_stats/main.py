from utils.stats_generator import StatsGenerator

if __name__ == "__main__":
    user = input("Enter GitHub Username: ")
    repo = input("Enter GitHub Repository: ")

    stats = StatsGenerator(repo_name=repo, user_name=user, verbose=True)

    stats.generate_all_stats()
    stats.stats.output_results()