import json
from .hybrid_search import HybridSearch
from .search_utils import (
    load_golden_dataset,
    load_movies,
)
from .semantic_search import SemanticSearch
from .reranking import (
    client,
    model
)

def precision_at_k(
    retrieved_docs: list[str], relevant_docs: set[str], k: int = 5
) -> float:
    top_k = retrieved_docs[:k]
    relevant_count = 0
    for doc in top_k:
        if doc in relevant_docs:
            relevant_count += 1
    return relevant_count / k


def recall_at_k(
    retrieved_docs: list[str], relevant_docs: set[str], k: int = 5
) -> float:
    top_k = retrieved_docs[:k]
    relevant_count = 0
    for doc in top_k:
        if doc in relevant_docs:
            relevant_count += 1

    return relevant_count / len(relevant_docs)

def f1_score(precision: float, recall: float) -> float:
    if precision + recall == 0:
        return 0
    return 2 * (precision * recall) / (precision + recall)


def evaluate_command(limit: int = 5) -> dict:
    movies = load_movies()
    test_cases = load_golden_dataset()

    semantic_search = SemanticSearch()
    semantic_search.load_or_create_embeddings(movies)
    hybrid_search = HybridSearch(movies)

    total_precision = 0
    results_by_query = {}
    for test_case in test_cases:
        query = test_case["query"]
        relevant_docs = set(test_case["relevant_docs"])
        search_results = hybrid_search.rrf_search(query, k=60, limit=limit)

        retrieved_docs = []
        for result in search_results:
            title = result.get("title", "")
            if title:
                retrieved_docs.append(title)

        precision = precision_at_k(retrieved_docs, relevant_docs, limit)
        recall = recall_at_k(retrieved_docs, relevant_docs, limit)
        f1_score = f1_score(precision, recall)

        results_by_query[query] = {
            "precision": precision,
            "recall": recall,
            "f1_score": f1_score,
            "retrieved": retrieved_docs[:limit],
            "relevant": list(relevant_docs),
        }

        total_precision += precision

    return {
        "test_cases_count": len(test_cases),
        "limit": limit,
        "results": results_by_query,
    }


def llm_judge_results(query: str, results: list[dict]):
    movie_titles = []
    for res in results:
        movie_titles.append(res["title"])

    prompt = f"""Rate how relevant each result is to this query on a 0-3 scale:

Query: "{query}"

Results:
{chr(10).join(movie_titles)}

Scale:
- 3: Highly relevant
- 2: Relevant
- 1: Marginally relevant
- 0: Not relevant

Do NOT give any numbers out than 0, 1, 2, or 3.

Return ONLY the scores in the same order you were given the documents. Return a valid JSON list, nothing else. For example:

[2, 0, 3, 2, 0, 1]"""
    
    response = client.models.generate_content(model=model, contents=prompt)
    ranking_text = (response.text or "").strip()

    parsed_scores = json.loads(ranking_text)

    for i, (title, score) in enumerate(zip(movie_titles, parsed_scores), start=1):
        print(f"{i}. {title}: {score}/3")