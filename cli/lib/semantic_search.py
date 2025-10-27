"""
Semantic search module using sentence transformers for embedding-based search.
"""

from sentence_transformers import SentenceTransformer

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
        print(f"Model loaded: {self.model}")
        print(f"Max sequence length: {self.model.max_seq_length}")

def verify_model():
    """
    Verify that the semantic search model loads correctly and print its information.
    """
    search = SemanticSearch()
    print(f"Model loaded: {search.model}")
    print(f"Max sequence length: {search.model.max_seq_length}")