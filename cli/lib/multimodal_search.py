from PIL import Image
from sentence_transformers import SentenceTransformer

from .search_utils import load_movies
from .semantic_search import cosine_similarity

class MultimodalSearch:
    def __init__(self, documents, model_name="clip-ViT-B-32"):
        self.model = SentenceTransformer(model_name)
        self.documents = documents

        self.texts = []
        for doc in documents:
            self.texts.append(f"{doc['title']}: {doc['description']}")

        self.text_embeddings = self.model.encode(self.texts, show_progress_bar=True)

    def embed_image(self, image_path):
        image = Image.open(image_path)

        embeddings = self.model.encode([image])
        
        return embeddings[0]
    
    def search_with_image(self, image_path):
        image_embedding = self.embed_image(image_path)

        results = []
        for i, text_embedding in enumerate(self.text_embeddings):
           score = cosine_similarity(image_embedding, text_embedding)
           doc = self.documents[i]
           results.append({
                'id': doc['id'],
                'score': score,
                'title': doc['title'],
                'description': doc['description']
           })
        
        results.sort(key=lambda x: x['score'], reverse=True)
        return results[:5]

def verify_image_embedding(image_path):
    multimodal_search = MultimodalSearch()
    embedding = multimodal_search.embed_image(image_path)

    print(f"Embedding shape: {embedding.shape[0]} dimensions")


def search_with_image(image_path):
    movies = load_movies()

    multimodal_search = MultimodalSearch(movies)
    results = multimodal_search.search_with_image(image_path)

    return results
    
