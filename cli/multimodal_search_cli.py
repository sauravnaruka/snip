import argparse

from lib.search_utils import DOCUMENT_PREVIEW_LENGTH
from lib.multimodal_search import search_with_image, verify_image_embedding


def main() -> None:
    parser = argparse.ArgumentParser(description="Multimodal Search CLI")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    verify_image_embedding_parser = subparsers.add_parser("verify_image_embedding", help="verify image embedding")
    verify_image_embedding_parser.add_argument("path", type=str, help="Image path")

    mage_search_parser = subparsers.add_parser("image_search", help="vsearch movie with image")
    mage_search_parser.add_argument("path", type=str, help="Image path")

    args = parser.parse_args()

    match args.command:
        case "verify_image_embedding":
            verify_image_embedding(args.path)
        case "image_search":
            results = search_with_image(args.path)

            for i, result in enumerate(results, start=1):
                print(f"{i}. {result['title']} (similarity: {result['score']:.3f})")
                print(f"    {result['description'][:DOCUMENT_PREVIEW_LENGTH]}")

if __name__ == "__main__":
    main()