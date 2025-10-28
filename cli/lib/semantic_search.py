"""
Semantic search module using sentence transformers for embedding-based search.
"""

from sentence_transformers import SentenceTransformer
import numpy as np
import os
import json

class SemanticSearch:
    """
    A semantic search engine that uses sentence transformers to find
    semantically similar content.
    """
    
    def __init__(self):
        """
        Initialize the SemanticSearch with the all-MiniLM-L6-v2 model.
        This model will be downloaded automatically on first use.
        """
        print("Loading semantic search model...")
        self.model = SentenceTransformer('all-MiniLM-L6-v2')
        self.embeddings = None
        self.documents = None
        self.document_map = {}
        print(f"Model loaded: {self.model}")
        print(f"Max sequence length: {self.model.max_seq_length}")


    def generate_embedding(self, text):
        """
        Generate an embedding for the given text using the model.
        """
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

        for doc in documents:
            self.document_map[doc['id']] = doc

        movie_strings = [f"{doc['title']}: {doc['description']}" for doc in documents]
        
        self.embeddings = self.model.encode(movie_strings, show_progress_bar=True)

        os.makedirs('cache', exist_ok=True)
        np.save('cache/movie_embeddings.npy', self.embeddings)

        return self.embeddings


    def load_or_create_embeddings(self, documents):
        self.documents = documents

        for doc in documents:
            self.document_map[doc['id']] = doc

        if os.path.exists('cache/movie_embeddings.npy'):
            self.embeddings = np.load('cache/movie_embeddings.npy')

            if len(self.embeddings) == len(documents):
                return self.embeddings
        
        return self.build_embeddings(documents)
    

    def search(self, query, limit):
        if self.embeddings is None:
            raise ValueError("No embeddings loaded. Call `load_or_create_embeddings` first.")

        query_embedding = self.generate_embedding(query)
        
        results = []
        for i, doc_embedding in enumerate(self.embeddings):
            similarity = cosine_similarity(query_embedding, doc_embedding)
            document = self.documents[i]
            results.append((similarity, document))
            
        results.sort(key=lambda x: x[0], reverse=True)

        top_results = []
        for score, doc in results[:limit]:
            top_results.append({
                'score': score,
                'id': doc['id'],
                'title': doc['title'],
                'description': doc['description']
            })
        
        return top_results


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

    with open('./data/movies.json', 'r') as f:
        data = json.load(f)
        movies = data['movies']

    embeddings = search.load_or_create_embeddings(movies)

    print(f"Number of docs:   {len(movies)}")
    print(f"Embeddings shape: {embeddings.shape[0]} vectors in {embeddings.shape[1]} dimensions")


def embed_query_text(query):
    search = SemanticSearch()

    embedding = search.generate_embedding(query)

    print(f"Query: {query}")
    print(f"First 5 dimensions: {embedding[:5]}")
    print(f"Shape: {embedding.shape}")


def cosine_similarity(vec1, vec2):
    dot_product = np.dot(vec1, vec2)
    norm1 = np.linalg.norm(vec1)
    norm2 = np.linalg.norm(vec2)

    if norm1 == 0 or norm2 == 0:
        return 0.0

    return dot_product / (norm1 * norm2)


def search_movie(query, limit):
    search = SemanticSearch()

    with open('./data/movies.json', 'r') as f:
        data = json.load(f)
        movies = data['movies']

    embeddings = search.load_or_create_embeddings(movies)

    results = search.search(query, limit)

    for i, result in enumerate(results, start=1):
        print(f"{i}. {result['title']} (score: {result['score']:.4f})")
        print(f"   {result['description']}")
        print()