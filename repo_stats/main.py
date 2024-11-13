import argparse
from utils.stats_generator import StatsGenerator

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog="repo_stats",
        description="Generate stats for a GitHub repository",
        epilog="Example: main.py -u <USER> -r <REPO>"
    )
    parser.add_argument("-u", "--user", nargs=1, help="The GitHub username of the repository owner")
    parser.add_argument("-r", "--repo", nargs=1, help="The name of the repository")
    args = parser.parse_args()

    if not args.user:
        user = input("Enter GitHub Username: ")
    else:
        user = args.user

    if not args.repo:
        repo = input("Enter GitHub Repository: ")
    else:
        repo = args.repo

    stats = StatsGenerator(repo_name=repo, user_name=user, verbose=True)

    stats.generate_all_stats()
    stats.stats.output_results()