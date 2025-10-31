#!/usr/bin/env python3

import argparse

from lib.chunked_semantic_search import embed_chunks
from lib.semantic_search import (
    DEFAULT_CHUNK_SIZE,
    DEFAULT_CHUNK_OVERLAP,
    verify_model,embed_text, 
    verify_embeddings, 
    embed_query_text, 
    semantic_search, 
    chunk_text,
    chunk_semantic_text,
)

def main() -> None:
    parser = argparse.ArgumentParser(description="Semantic Search CLI")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    # Command: verify
    subparsers.add_parser("verify", help="Verify that the embedding model is loaded")

    # Command: embed_text <text>
    single_embed_parser = subparsers.add_parser("embed_text", help="Generate embedding for given text")
    single_embed_parser.add_argument("text", type=str, help="Text to embed")

    subparsers.add_parser("verify_embeddings", help="Generate embeddings for the entire movie dataset; to save them to disk")
    
    embed_query_parser = subparsers.add_parser("embedquery", help="Generate embeddings for the user query.")
    embed_query_parser.add_argument("query", type=str, help="Query to embed")

    # Command: search
    search_parser = subparsers.add_parser("search", help="Search movie database to get similar movies")
    search_parser.add_argument("query", type=str, help="Input text to search movie for")
    search_parser.add_argument("--limit", type=int, default=5, help="Limit the number of search results (default: 5)")

    # Command: chunk
    chunk_parser = subparsers.add_parser("chunk", help="add fixed size chunking to split long text into smaller pieces for embedding")
    chunk_parser.add_argument("text", type=str, help="Input text to create chunk")
    chunk_parser.add_argument("--chunk-size", type=int, default=DEFAULT_CHUNK_SIZE, help="Define the size of the chunk")
    chunk_parser.add_argument("--overlap", type=int, default=DEFAULT_CHUNK_OVERLAP, help="Define the overlap for the chunk")

    # Command: semantic_chunk
    semantic_chunk_parser = subparsers.add_parser("semantic_chunk", help="add semantic chunking to split long text into smaller pieces for embedding")
    semantic_chunk_parser.add_argument("text", type=str, help="Input text to create chunk")
    semantic_chunk_parser.add_argument("--max-chunk-size", type=int, default=4, help="Define the size of the chunk")
    semantic_chunk_parser.add_argument("--overlap", type=int, default=0, help="Define the overlap for the chunk")

    #Command: embed_chunks
    subparsers.add_parser("embed_chunks", help="Generate embeddings for chunked documents")

    args = parser.parse_args()

    match args.command:
        case "verify":
            verify_model()
        case "embed_text":
            embed_text(args.text)
        case "verify_embeddings":
            verify_embeddings()
        case "embedquery":
            embed_query_text(args.query)
        case "search":
            semantic_search(args.query, args.limit)
        case "chunk":
            chunk_text(args.text, args.chunk_size, args.overlap)
        case "semantic_chunk":
            chunk_semantic_text(args.text, args.max_chunk_size, args.overlap)
        case "embed_chunks":
            embed_chunks()
        case _:
            parser.print_help()

if __name__ == "__main__":
    main()