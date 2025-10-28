#!/usr/bin/env python3

import argparse
from lib.semantic_search import verify_model,embed_text, verify_embeddings, embed_query_text, search_movie

def main():
    parser = argparse.ArgumentParser(description="Semantic Search CLI")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    # Command: verify
    subparsers.add_parser("verify", help="Verify the semantic search model loads correctly")

    # Command: embed_text <text>
    embed_parser = subparsers.add_parser("embed_text", help="Generate embedding for given text")
    embed_parser.add_argument("text", type=str, help="Input text to generate embedding for")

    subparsers.add_parser("verify_embeddings", help="Generate embeddings for the entire movie dataset; to save them to disk")
    
    embed_query = subparsers.add_parser("embedquery", help="Generate embeddings for the user query.")
    embed_query.add_argument("text", type=str, help="Input text to generate embedding for")

    # Command: search
    search_query = subparsers.add_parser("search", help="Search movie database to get similar movies")
    search_query.add_argument("query", type=str, help="Input text to search movie for")
    search_query.add_argument("--limit", type=int, default=5, help="Limit the number of search results (default: 5)")

    args = parser.parse_args()

    match args.command:
        case "verify":
            verify_model()
        case "embed_text":
            embed_text(args.text)
        case "verify_embeddings":
            verify_embeddings()
        case "embedquery":
            embed_query_text(args.text)
        case "search":
            search_movie(args.query, args.limit)
        case _:
            parser.print_help()

if __name__ == "__main__":
    main()