import json
import os
import numpy as np
from .semantic_search import (
    SemanticSearch,
    semantic_chunking
)

from .search_utils import (
    CHUNK_EMBEDDINGS_PATH,
    CHUNK_METADATA_PATH,
    DEFAULT_SEMANTIC_CHUNK_SIZE,
    DEFAULT_CHUNK_OVERLAP,
    load_movies
)

class ChunkedSemanticSearch(SemanticSearch):
    def __init__(self, model_name = "all-MiniLM-L6-v2") -> None:
        super().__init__(model_name)
        self.chunk_embeddings = None
        self.chunk_metadata = None

    def build_chunk_embeddings(self, documents):
        self.documents = documents
        self.document_map = {}

        for doc in documents:
            self.document_map[doc['id']] = doc

        all_chunks = []
        chunk_metadata = []

        for doc_idx, doc in enumerate(documents):
            if not doc['description'] or not doc['description'].strip():
                continue
             
            chunks = semantic_chunking(
                doc['description'], 
                max_chunk_size=DEFAULT_SEMANTIC_CHUNK_SIZE,
                overlap=DEFAULT_CHUNK_OVERLAP
            )

            for chunk_idx, chunk in enumerate(chunks):
                 all_chunks.append(chunk)
                 chunk_metadata.append({
                    'movie_idx': doc_idx,
                    'chunk_idx': chunk_idx,
                    'total_chunks': len(chunks)
                })
                 
        self.chunk_embeddings = self.model.encode(all_chunks, show_progress_bar=True)
        self.chunk_metadata = chunk_metadata

        os.makedirs(os.path.dirname(CHUNK_EMBEDDINGS_PATH), exist_ok=True)
        np.save(CHUNK_EMBEDDINGS_PATH, self.chunk_embeddings)

        with open(CHUNK_METADATA_PATH, 'w') as f:
            json.dump(
                {"chunks": chunk_metadata, "total_chunks": len(all_chunks)}, f, indent=2
            )
        
        return self.chunk_embeddings
    
    def load_or_create_chunk_embeddings(self, documents: list[dict]) -> np.ndarray:
        self.documents = documents
        self.document_map = {}

        for doc in documents:
            self.document_map[doc['id']] = doc

        if os.path.exists(CHUNK_EMBEDDINGS_PATH) and os.path.exists(CHUNK_METADATA_PATH):
            self.chunk_embeddings = np.load(CHUNK_EMBEDDINGS_PATH)

            with open(CHUNK_METADATA_PATH, 'r') as f:
                data = json.load(f)
                self.chunk_metadata = data['chunks']
            
            return self.chunk_embeddings

        return self.build_chunk_embeddings(documents)
            
def embed_chunks():
    documents = load_movies()

    chunk_semantic_search = ChunkedSemanticSearch()
    embeddings = chunk_semantic_search.load_or_create_chunk_embeddings(documents)

    print(f"Generated {len(embeddings)} chunked embeddings")