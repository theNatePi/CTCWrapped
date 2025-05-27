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
    parser.add_argument("-v", "--verbose", nargs = 1, help = "Verbose output, yes / no")
    args = parser.parse_args()

    if not args.user:
        user = input("Enter GitHub Username: ")
    else:
        user = args.user[0]

    if not args.repo:
        repo = input("Enter GitHub Repository: ")
    else:
        repo = args.repo[0]

    if not args.verbose:
        print("Loading...")
        verbose = "no"
    else:
        verbose = args.verbose[0]

    verbose = True if verbose == "yes" else False

    stats = StatsGenerator(repo_name=repo, user_name=user, verbose=verbose)

    stats.generate_all_stats()
    stats.stats.output_results()
