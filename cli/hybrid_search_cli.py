import argparse

from lib.hybrid_search import (
    normalize_score
)

def main() -> None:
    parser = argparse.ArgumentParser(description="Hybrid Search CLI")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    #Command: normalize
    normalize_parser = subparsers.add_parser("normalize", help="Normalize the subparser")
    normalize_parser.add_argument("scores", nargs="+", type=float)

    args = parser.parse_args()

    match args.command:
        case "normalize":
            results = normalize_score(args.scores)
            for result in results:
                print(f"* {result:.4f}")
        case _:
            parser.print_help()


if __name__ == "__main__":
    main()