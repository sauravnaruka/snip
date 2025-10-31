"""
Semantic search module using sentence transformers for embedding-based search.
"""

import re
from sentence_transformers import SentenceTransformer
import numpy as np
import os

from .search_utils import (
    DEFAULT_CHUNK_SIZE,
    DEFAULT_CHUNK_OVERLAP,
    DEFAULT_SEMANTIC_CHUNK_SIZE,
    MOVIE_EMBEDDINGS_PATH,
    load_movies
)

class SemanticSearch:
    """
    A semantic search engine that uses sentence transformers to find
    semantically similar content.
    """
    
    def __init__(self, model_name="all-MiniLM-L6-v2"):
        self.model = SentenceTransformer(model_name)
        self.embeddings = None
        self.documents = None
        self.document_map = {}


    def generate_embedding(self, text):
        if not text or not text.strip():
            raise ValueError("cannot generate embedding for empty text")
        
        return self.model.encode([text])[0]


    def build_embeddings(self, documents):
        """
        Build embeddings for a list of movie documents.
        
        Args:
            documents: List of dictionaries, each representing a movie with 'id', 'title', and 'description'
        
        Returns:
            The generated embeddings as a numpy array
        """
                
        self.documents = documents
        self.document_map = {}
        movie_strings = []

        for doc in documents:
            self.document_map[doc['id']] = doc
            movie_strings.append(f"{doc['title']}: {doc['description']}")
        
        self.embeddings = self.model.encode(movie_strings, show_progress_bar=True)

        os.makedirs(os.path.dirname(MOVIE_EMBEDDINGS_PATH), exist_ok=True)
        np.save(MOVIE_EMBEDDINGS_PATH, self.embeddings)
        return self.embeddings


    def load_or_create_embeddings(self, documents):
        self.documents = documents
        self.document_map = {}

        for doc in documents:
            self.document_map[doc['id']] = doc

        if os.path.exists(MOVIE_EMBEDDINGS_PATH):
            self.embeddings = np.load(MOVIE_EMBEDDINGS_PATH)

            if len(self.embeddings) == len(documents):
                return self.embeddings
        
        return self.build_embeddings(documents)
    

    def search(self, query, limit):
        if self.embeddings is None or self.embeddings.size == 0:
            raise ValueError("No embeddings loaded. Call `load_or_create_embeddings` first.")

        if self.documents is None or len(self.documents) == 0:
            raise ValueError(
                "No documents loaded. Call `load_or_create_embeddings` first."
            )
        
        query_embedding = self.generate_embedding(query)
        
        similarities = []
        for i, doc_embedding in enumerate(self.embeddings):
            similarity = cosine_similarity(query_embedding, doc_embedding)
            similarities.append((similarity, self.documents[i]))
            
        similarities.sort(key=lambda x: x[0], reverse=True)

        results = []
        for score, doc in similarities[:limit]:
            results.append({
                'score': score,
                'title': doc['title'],
                'description': doc['description']
            })
        
        return results

def cosine_similarity(vec1, vec2):
    dot_product = np.dot(vec1, vec2)
    norm1 = np.linalg.norm(vec1)
    norm2 = np.linalg.norm(vec2)

    if norm1 == 0 or norm2 == 0:
        return 0.0

    return dot_product / (norm1 * norm2)

def verify_model():
    """
    Verify that the semantic search model loads correctly and print its information.
    """
    search = SemanticSearch()
    print(f"Model loaded: {search.model}")
    print(f"Max sequence length: {search.model.max_seq_length}")



def embed_text(text):
    """
    Create a SemanticSearch instance, generate an embedding, and print it.
    """
    search = SemanticSearch()
    embedding = search.generate_embedding(text)
    
    print(f"Text: {text}")
    print(f"First 3 dimensions: {embedding[:3]}")
    print(f"Dimensions: {embedding.shape[0]}")


def verify_embeddings():
    search = SemanticSearch()
    movies = load_movies()
    embeddings = search.load_or_create_embeddings(movies)
    print(f"Number of docs:   {len(movies)}")
    print(f"Embeddings shape: {embeddings.shape[0]} vectors in {embeddings.shape[1]} dimensions")


def embed_query_text(query):
    search = SemanticSearch()
    embedding = search.generate_embedding(query)
    print(f"Query: {query}")
    print(f"First 5 dimensions: {embedding[:5]}")
    print(f"Shape: {embedding.shape}")


def semantic_search(query, limit):
    search = SemanticSearch()
    movies = load_movies()
    search.load_or_create_embeddings(movies)

    results = search.search(query, limit)
    print(f"Query: {query}")
    print(f"Top {len(results)} results:")
    print()

    for i, result in enumerate(results, start=1):
        print(f"{i}. {result['title']} (score: {result['score']:.4f})")
        print(f"   {result['description']}")
        print()


def chunk_text(text: str, chunk_size: int = DEFAULT_CHUNK_SIZE, overlap: int = DEFAULT_CHUNK_OVERLAP) -> None:
    chunks = fixed_size_chunking(text, chunk_size, overlap)
    print(f"Chunking {len(text)} characters")
    for i, chunk in enumerate(chunks):
        print(f"{i + 1}. {chunk}")


def fixed_size_chunking(text: str, chunk_size: int = DEFAULT_CHUNK_SIZE, overlap: int = DEFAULT_CHUNK_OVERLAP) -> list[str]:
    words = text.split()
    chunks = []
    i = 0
    while i < len(words):
        chunk_words = words[i : i + chunk_size]
        chunks.append(" ".join(chunk_words))
        i += chunk_size - overlap
        print(f"i={i}")

    return chunks


def chunk_semantic_text(text: str, max_chunk_size: int, overlap: int) -> None:
    print(f"Semantically chunking {len(text)} characters")
    chunks = semantic_chunking(text, max_chunk_size, overlap)

    for i, chunk in enumerate(chunks, start=1):
        print(f"{i}. {chunk}")

def semantic_chunking(text: str, max_chunk_size: int = DEFAULT_SEMANTIC_CHUNK_SIZE, overlap: int = DEFAULT_CHUNK_OVERLAP) -> list[str]:
    text = text.strip()
    if not text:
        return []
    
    # Split text into sentences
    sentences = re.split(r"(?<=[.!?])\s+", text)

    if len(sentences) == 1 and not text.endswith(('.', '!', '?')):
        # Treat the whole text as one sentence
        return [text]

    chunks = []
    
    i = 0
    while i < len(sentences):
        # Get up to max_chunk_size sentences
        chunk_sentences = sentences[i : i + max_chunk_size]

        cleaned_sentences = []
        for chunk_sentence in chunk_sentences:
            cleaned_sentences.append(chunk_sentence.strip())
        if not cleaned_sentences:
            continue

        chunks.append(" ".join(cleaned_sentences))

        # Move forward by (max_chunk_size - overlap)
        i += max_chunk_size - overlap

        if overlap > 0 and i + overlap >= len(sentences):
            break
    
    return chunks