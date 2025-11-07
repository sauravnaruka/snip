import argparse

from lib.hybrid_search import (
    normalize_scores,
    weighted_search_command,
    rrf_search_command
)

from lib.search_utils import (
    DEFAULT_SEARCH_LIMIT,
    DEFAULT_ALPHA,
    DEFAULT_RRF_K
)

def main() -> None:
    parser = argparse.ArgumentParser(description="Hybrid Search CLI")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    #Command: normalize
    normalize_parser = subparsers.add_parser("normalize", help="Normalize the subparser")
    normalize_parser.add_argument("scores", nargs="+", type=float)

    #Command: weighted-search
    weighted_search_parser = subparsers.add_parser("weighted-search", help="Hybrid search based on keyword and semantic meaning")
    weighted_search_parser.add_argument("query", type=str, help="Search query")
    weighted_search_parser.add_argument("--alpha", type=float, default=DEFAULT_ALPHA, help="alpha to control weighting between the keyword and semantic search")
    weighted_search_parser.add_argument("--limit", type=int, default=DEFAULT_SEARCH_LIMIT, help="the number of search results to filter out")

    #Command: rrf_search
    rrf_search_parser = subparsers.add_parser("rrf-search", help="Perform Reciprocal Rank Fusion search")
    rrf_search_parser.add_argument("query", type=str)
    rrf_search_parser.add_argument("--k", type=int, default=DEFAULT_RRF_K, help="controls how much more weight we give to higher-ranked results vs. lower-ranked ones.")
    rrf_search_parser.add_argument("--limit", type=int, default=DEFAULT_SEARCH_LIMIT, help="the number of search results to filter out")
    rrf_search_parser.add_argument("--enhance", type=str, choices=["spell", "rewrite", "expand"], help="Query enhancement method")
    rrf_search_parser.add_argument("--rerank-method", type=str, choices=["individual", "batch", "cross_encoder"], help="search result re-ranking method")

    args = parser.parse_args()

    match args.command:
        case "normalize":
            results = normalize_scores(args.scores)
            for result in results:
                print(f"* {result:.4f}")
        case "weighted-search":
            result = weighted_search_command(args.query, args.alpha, args.limit)
            
            print(f"Weighted Hybrid Search Results for '{result['query']}' (alpha={result['alpha']}):")
            print(f"  Alpha {result['alpha']}: {int(result['alpha'] * 100)}% Keyword, {int((1 - result['alpha']) * 100)}% Semantic")
            for i, res in enumerate(result["results"], 1):
                print(f"{i}. {res['title']}")
                print(f"   Hybrid Score: {res.get('score', 0):.3f}")
                metadata = res.get("metadata", {})
                if "bm25_score" in metadata and "semantic_score" in metadata:
                    print(f"   BM25: {metadata['bm25_score']:.3f}, Semantic: {metadata['semantic_score']:.3f}")
                print(f"   {res['document'][:100]}...")
                print()
        case "rrf-search":
            result = rrf_search_command(args.query, args.k, args.enhance, args.rerank_method, args.limit)
            
            if result["enhanced_query"]:
                print(f"Enhanced query ({result['enhance_method']}): '{result['original_query']}' -> '{result['enhanced_query']}'\n")

            if result["reranked"]:
                print(f"Reranking top {len(result['results'])} results using {result['rerank_method']} method...\n")
            
            print(f"Reciprocal Rank Fusion Results for '{result['query']}' (k={result['k']}):")
            
            for i, res in enumerate(result["results"], 1):
                print(f"{i}. {res['title']}")
                if "individual_score" in res:
                    print(f"   Rerank Score: {res.get('individual_score', 0):.3f}/10")
                if "batch_rank" in res:
                    print(f"   Rerank Rank: {res.get('batch_rank', 0)}")
                if "cross_encoder_score" in res:
                    print(f"   Cross Encoder Score: {res.get('cross_encoder_score', 0)}")
               
                print(f"   RRF Score: {res.get('score', 0):.3f}")

                metadata = res.get("metadata", {})
                ranks = []
                if metadata.get("bm25_rank"):
                    ranks.append(f"BM25 Rank: {metadata['bm25_rank']}")
                if metadata.get("semantic_rank"):
                    ranks.append(f"Semantic Rank: {metadata['semantic_rank']}")
                if ranks:
                    print(f"   {', '.join(ranks)}")

                print(f"   {res['document'][:100]}...")
                print()
        case _:
            parser.print_help()


if __name__ == "__main__":
    main()