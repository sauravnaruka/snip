import argparse

from lib.describe_image import describe_image


def main() -> None:
    parser = argparse.ArgumentParser(description="Multimodal Query Rewriting CLI")
    parser.add_argument("--query", type=str, required=True, help="Text query to rewrite based on the image")
    parser.add_argument("--image", type=str, required=True, help="Path to an image file")

    args = parser.parse_args()

    describe_image(args.query, args.image)

if __name__ == "__main__":
    main()