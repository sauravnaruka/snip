import mimetypes
import os
from google import genai
from google.genai import types as genai_types
from dotenv import load_dotenv

from .search_utils import load_image

load_dotenv()
api_key = os.getenv("GEMINI_API_KEY")
client = genai.Client(api_key=api_key)
model = "gemini-2.0-flash"

def describe_image(query, image_path):
    mime, _ = mimetypes.guess_type(image_path)
    mime = mime or "image/jpeg"

    image_data = load_image(image_path)

    system_prompt = f"""
Given the included image and text query, rewrite the text query to improve search results from a movie database. Make sure to:
- Synthesize visual and textual information
- Focus on movie-specific details (actors, scenes, style, etc.)
- Return only the rewritten query, without any additional commentary
"""
    
    parts = [
        system_prompt,
        genai_types.Part.from_bytes(data=image_data, mime_type=mime),
        query.strip(),
    ]

    response = client.models.generate_content(model=model, contents=parts)
    if response.text is None:
        raise RuntimeError("No text in Gemini response")

    print(f"Rewritten query: {response.text.strip()}")
    if response.usage_metadata is not None:
        print(f"Total tokens:    {response.usage_metadata.total_token_count}")