import json
import os
import re
from time import sleep

from dotenv import load_dotenv
from google import genai

load_dotenv()
api_key = os.getenv("GEMINI_API_KEY")
client = genai.Client(api_key=api_key)
model = "gemini-2.0-flash"

def rerank(
    query: str, documents: list[dict], method: str = "batch", limit: int = 5
) -> list[dict]:
    match method:
        case "individual":
            return llm_rerank_individual(query, documents, limit)
        case "batch":
            return llm_rerank_batch(query, documents, limit)
        case _:
            return documents[:limit]


def llm_rerank_individual(
    query: str, documents: list[dict], limit: int = 5
) -> list[dict]:
    scored_docs = []

    for doc in documents:
        prompt = f"""Rate how well this movie matches the search query.

Query: "{query}"
Movie: {doc.get("title", "")} - {doc.get("document", "")}

Consider:
- Direct relevance to query
- User intent (what they're looking for)
- Content appropriateness

Rate 0-10 (10 = perfect match).
Give me ONLY the number in your response, no other text or explanation.

Score:"""

        response = client.models.generate_content(model=model, contents=prompt)
        score_text = (response.text or "").strip()
        score = int(score_text)
        scored_docs.append({**doc, "individual_score": score})
        sleep(3)

    scored_docs.sort(key=lambda x: x["individual_score"], reverse=True)
    return scored_docs[:limit]



def llm_rerank_batch(query: str, documents: list[dict], limit: int = 5):
    if not documents:
        return []

    doc_list_str = "\n".join(
        f"{doc['id']}: {doc['title']} - {doc['document']}"
        for doc in documents
    )

    prompt = f"""Rank these movies by relevance to the search query.

Query: "{query}"

Movies:
{doc_list_str}

Return ONLY the IDs in order of relevance (best match first). Return a valid JSON list, nothing else. For example:

[75, 12, 34, 2, 1]
"""
    response = client.models.generate_content(model=model, contents=prompt)
    print(f"LLM response.text: {response.text}")

    text = response.candidates[0].content.parts[0].text.strip()

    print(f"LLM text: {text}")
    try:
        result_ids = json.loads(text)
    except json.JSONDecodeError:
        print("Warning: JSON decoding failed, attempting regex fallback")
        result_ids = [int(x) for x in re.findall(r"\d+", text)]
    print(f"Parsed result IDs: {result_ids}")
    result = []
    for doc_id in result_ids:
        doc = next((d for d in documents if d["id"] == doc_id), None)
        if doc:
            result.append(doc)
    
    return result[:limit]