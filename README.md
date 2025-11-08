# Go Application — Development Setup

This project is a Go-based application.  
Follow these steps to set up and run it in development mode.

---

## 🧱 Prerequisites

Make sure you have the following installed:

- [Go](https://go.dev/dl/) **version 1.23+**
- [Git](https://git-scm.com/)
- (Optional) [Air](https://github.com/cosmtrek/air) for live reload during development

Check your Go version:

```bash
go version
```

## Getting Started

1. Initialize the Go module (if not already)

```bash
go mod init github.com/sauravnaruka/snip
```

2. Install dependencies

```bash
go mod tidy
```

If you see an error like
`no required module provides package github.com/joho/godotenv,`
run:

```bash
go get github.com/joho/godotenv
```

3. Environment Variables
   Create a `.env` file in the project root:
   Add following fields to `.env` file

- `DB_FILE_PATH`: Location of data file
- `STOP_WORD_FILE_PATH`: Location of stop word file
- `BM25_K1`: Tunable parameter that controls the diminishing returns for BM25 algo. Default is `1.5`
- `BM25_B`: Tunable parameter that controls how much we care about document length.

4. Data file
   Create data file at .<Project-root>/data/movies.json
   File path can be anything also also.

- Add the path to `.gitignore`
- Update path in `.env` file `DB_FILE_PATH`

## Running the Application

Option 1 — Run directly

```bash
go run .
```

Option 2 — Build and run

```bash
go build -o snip
./snip
```

> [!IMPORTANT] <br> **Important:** Before running the program, make sure to build inverted index using the command
> `go run . build` or `snip build`

## 🧰 Python CLI Utilities

Some auxiliary tools for data preparation and experimentation are implemented in Python under the [`cli/`](./cli) folder.  
These scripts are isolated from the Go application and use a local virtual environment managed by **uv**.

If you plan to use or modify the Python utilities, see the detailed setup instructions in [`cli/README.md`](./cli/README.md).

## Search

Search works in following steps:
**Prerequisite**: [Inverted Index](#inverted-index) needs to be build before the search operation, using the `build` command

1. User query [text is preprocessed](#Text-Preprocessing) and converted into a list of tokens
2. List of tokens are searched based on [Inverted Index](#inverted-index). In Inverted index we maintain token -> Movie mapping. Therefore we fetch movies for each token. Once we get 5 results, we break from the search

### Text Preprocessing

To improve keyword-based search accuracy, we need to normalize both the query and the target text through a few preprocessing steps. For instance, the words run, Run, and running should all be recognized as the same keyword during search.

Following are text processing Steps

#### 1. Case Insensitivity

Search should be case-insensitive. Both the query and the target text are converted to the same case (usually lowercase) so that run and Run are treated as the same word.

#### 2. Punctuation Handling

Punctuation marks should not affect search results. For example, the tokens run, run., and run! should all be treated as the same word after removing punctuation.

#### 3. Tokenization

The input text needs to be split into individual tokens (words) for further processing and matching.For example:

- input `"Running fast, he reached the goal."`
- Output `["Running", "fast", "he", "reached", "the", "goal"]`

#### 4. Remove Stop Words

words that don't have much semantic meaning are called `stop words`. We will remove them from our search. Example of stop words are `the`, `a`, `is`, `of`, `in`

#### 5. Stemming

Each token should be reduced to its stem (base) form so that words derived from the same root are matched together. For example `running` should convert to `run`

### Keyword Search

#### TF-IDF

Term Frequency & Inverse Dcoument Frequency are often used together to create TF-IDF.

```
TF-IDF = TF * IDF
```

- TF (Term Frequency): How often a term appears in a document
- IDF (Inverse Document Frequency): How rare a term is across all document

- Frequent words gets higger TF score
- Rare words gets high IDF score
- Best matches have both high TF & high IDF

##### Inverted Index

An inverted index make search fast, it's like a SQL database index but for text. Instead of searching for all the documents each time a user searches, we build and index for fast lookup.

- Forward index maps location to value. eg: `doc1: [matrix, reload]`
- Inverted index maps value (token) to location. eg: `matrix: [doc1, doc5]`

We build inverted index in `internal/search` module.

1. We use struct `InvertedIndex` to store

- `Index` which maps tokens to set of document IDs
- `DocMap` which maps document IDs to the movie object

2. When `BuildInvertedIndex` called from cmd module it iterate through all the movies, create token from title and description, then call `addDocument`

3. `addDcoument` create a mapping of token and corresponding document id

##### Term Frequency

Term frequency (TF) measures how often a word appears in a document.

Early search engine relied on Term Frequency only. Where they use to scan a text and find frequency of token. A document containg a lot of token user searched considered more relevant.

However this led to keyword over use called keyword stuffing.

In out program, we also building `TermFrequencies`, It's a dictionary of document ID with a dictionary of token and it's frequency in the document.

##### Inverse Document Frequency (IDF)

Inverse Document Frequency (IDF) prioritize the words which are rare over the words which are generic.

For example the text "A movie about an actor who becomes a coding instructor".
In a movie database, the word movie or actor will be used multiple times while word like coding or instructor might be less frequent. With Inverse Dcoument Frequency we can prioritize words like 'coding' or 'instructor' over the generic word for the dataset like 'actor' or 'movie'

**Document frequency (DF)** measures how many documents in the dataset contain a term. The more documents a term appears in, the bigger its value, we don't want that. So we take **inverse**, It's because we want rare terms to have higher scores.

```
math.log((doc_count + 1) / (term_doc_count + 1))
```

#### BM25

BM25 ([Okapi BM25](https://en.wikipedia.org/wiki/Okapi_BM25)) is an improvement over [TF-IDF](#TF-IDF).

The formula for BM25 Search is:

```
BM25 = bm25_tf * bm25_idf

Where,
- bm25_idf = log((doc_count - term_doc_count + 0.5) / (term_doc_count + 0.5) + 1)
- bm25_tf = (TF * (k1 + 1)) / (TF + k1 * length_norm)
- length_norm = 1 - b + (b * (doc_length / avg_doc_length))
```

##### BM25 Improvements over IDF

BM25 gives More stable IDF calculations.

The origina [IDF](<#Inverse_Document_Frequency_(IDF)>) calculation is as follows:

```
math.log((doc_count + 1) / (term_doc_count + 1))
```

The problems with this formula:

- Division by zero. We solved this by adding 1 in above formula.
- In case of rare term we get high IDF.
- Very common terms can get negative scores.

**BM25-IDF folmula is**

```
BM25-IDF = log((doc_count - term_doc_count + 0.5) / (term_doc_count + 0.5) + 1)
```

Where,

- `doc_count - term_doc_count + 0.5` is count of documents without the term + smoothing (laplace smoothing) to prevent division by `0`
- `term_doc_count + 0.5` is count of documents with the term + smoothing (laplace smoothing) to prevent division by `0`

##### BM25 Improvements over Term Frequency (TF)

BM25 does term frequency saturation. In standard (TF)[#Term_Frequency] we simply caclulate the number of times a term apeared in token.

The issue is that a single term appearing a lot of times can get exponentially more weight. For example

- Document A: Title: Run Run Run, Description: A thrilling chase where the hero runs to save the day. Run faster!
- Document B: Title: Runaway Bride, Description: A romantic comedy about a bride who escapes her wedding.

Now when user search `run bride`, the document A gets higher prefrence because the number of times run apeared in document A.

BM25 uses diminishing returns after certain points. Additional occurrences matter less.

```
BM25-TF = (TF * (k1 + 1)) / (TF + k1)
```

Where,

- `TF` is the term frequency
- `K1` is tunable parameter that controls the diminsihing returns. A common value is `1.5`.

| Term Frequency | Basic TF | BM25 TF (k₁ = 1.5) |
| -------------- | -------- | ------------------ |
| 1              | 1        | 1.0                |
| 2              | 2        | 1.4                |
| 5              | 5        | 1.9                |
| 10             | 10       | 2.2                |
| 20             | 20       | 2.3                |

##### BM25 Normaliz document Length

The (TF)[#Term_Frequency] calculation benifit longer text as longer text has more words thus higher TF score.

BM25 improves the TF calculation by normalizing document length, ensuring longer documents don't get unfair advantanges over shorter, more focused once.

```
length_norm = 1 - b + (b * (doc_length / avg_doc_length))

BM25-TF = (TF * (k1 + 1)) / (TF + k1 * length_norm)
```

- If `b=0` then length norm is always 1. Thus no effect
- If `b=1` then full normalization is applied
- Usually the value of `b=0.75`

**Key Insight**

- Longer documents are penalized
- Shorter documents are boosted

### Semantic Search

Semantic search looks for similar meaning or the intent of the query rather then trying to match the sub-string or tokens

[Embeddings](https://www.cloudflare.com/en-in/learning/ai/what-are-embeddings/): are numerical representations of text that capture the meaning of words.

For semantic search we will use embeddings as [vectors](<https://en.wikipedia.org/wiki/Vector_(mathematics_and_physics)>). Vectors are a list of numbers that represent two things. First it is numerical representation of text, after all it is a embedding. Second, it represent a point in space, specifically, a direction and a magnitude from the origin `(0,0)`

We use vector as now the distance between the vectors represents how similar the meaning of the words are.

The process of converting text to vector is essentially a machine learning problem where a lot of data is used and a lot of computations are performed on the data to learn patterns about how different words and phrases relate to each other.

#### Vector Dimensions

Dimensions are just list of numbers or entries making one dimension. To capture semantic value we use vector some times containg 300+ dimensions.

Then we calculate distance between two vectors which tells us how similar or different are these things.

#### Vector Operations

We can perform actions on vector to optimize it further

##### Vector Addition

Vector addition is useful for combining/mixing concepts. For example, I want result that's like this and like that

##### Vector Subtraction

Vector subtraction is useful for removing concepts. For example, I want a result that's like this but not that.

##### Vector Dot Product

We perform dot product to see how similar two embeddings (vector) are. The dot product measure how much two vector "point in the same direction"

It's calculated by multiplying corresponding elements and then summing the result.

Python has popular library called NumPy to handle vector math.

> [!tip]
>
> - The More similar the vectors are, the higher the dot product will be.
> - If Opposite are the vectors, the dot product will be negative

Formula for the Dot product is

```
dot_product(A,B) = A_x \times B_x + A_y \times B_y
```

Code for same is

```py
def dot(vec1, vec2):
    if len(vec1) != len(vec2):
        raise ValueError("vectors must be the same length")
    total = 0.0
    for i in range(len(vec1)):
        total += vec1[i] * vec2[i]
    return total
```

##### Cosine Similarity

The dot product is a useful way to measure the similarity between two vectors.
However, it is influenced by the magnitude (length) of those vectors — which represents strength or confidence — rather than just their direction.

In semantic search, we usually care about how similar the meanings (directions) are, not how strong or confident each vector is.

That’s where cosine similarity comes in — it measures the angle between two vectors, effectively ignoring their lengths.

The formula for cosine similarity is:

```
cosine_similarity = dot_product(A, B) / (magnitude(A) × magnitude(B))
```

###### Relationship Between Dot Product and Cosine Similarity

The dot product of two vectors can be expressed in two equivalent ways:

1. Algebraic form (components-wise):

```
dot_product(A,B) = A_x \times B_x + A_y \times B_y
```

2. Geometric form (based on the angle between them):

```
dot_product(A,B) = |A| |B| cos(θ)
```

Where,

- ∣A∣ = magnitude (length) of vector A
- ∣B∣ = magnitude (length) of vector B
- θ = angle between the two vectors

Magnitude is also called **Euclidean norm**. The code for the same is:

```py
def euclidean_norm(vec):
    total = 0.0
    for x in vec:
        total += x**2

    return total**0.5
```

Following is the formula to calculate magnitude. This should be reminiscent of the Pythagorean theorem.

```
|A| = sqrt((A_x)^2 + (A_y)^2)
```

By substituting the two forms of the dot product, we get:

```
cos(θ) = \frac{A_x \times B_x + A_y \times B_y} / {|A| \times |B|}
```

This gives us the formula for cosine similarity:

```
cosine_similarity = dot_product(A, B) / (magnitude(A) × magnitude(B))
```

Now, the value of cosine_similarity ranges from `-1.0` to `1.0`. Where

- `1.0`: Vector point in exactly the same direction (perfectly similar)
- `0.0`: Vectors are perpendicular (no similarity)
- `-1.0`: Vectors point in opposite directions (perfectly dissimilar)

In summary it works in two steps:

1. Calculate similarity: The dot product measures how much vectors align
2. Remove length bias: Dividing by magnitudes removes the effect of vector size

###### Why Cosine Similarity

We want to use the same metric for similarity comparison as used by the model during the training process. The model we using is `all-MiniLM-L6-v2` which used cosine similarity during training.

Another great news is that cosine similarity is used by most embedding models

#### Creating Embeddings in Python

We can create embeddings by using a module and [SentenceTransformer](https://sbert.net/#). We using [sentence-transformers/all-MiniLM-L6-v2](https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2) in the project.

We can load modal like this

```py
from sentence_transformers import SentenceTransformer
model = SentenceTransformer('sentence-transformers/all-MiniLM-L6-v2')

embeddings = model.encode(sentences)
```

#### Locality-Sensitive Hashing

Checking a query against all the avialble data makes the query slow. A technique called Locality Sensitive Hashing (LSH), can speed up the search however the accurance can be reduced.

In Locality Sensitive Hashing (LSH) we pre-group the similar vectors into buckets. When user query, we search only in specific buckets, thus speeding up the search. However, any relevant vector /search item is in different bucket then we will miss out the item in search result.

It's a trade off between speed and accuracy (in ML terms you get lower recall). It should be used only when computation speed is a priority over perfect accuracy.

#### Vector Databases

We use a vector database instead of storing embeddings in arrays or the file system because it provides several important advantages.

First, it enables sub-linear time similarity search through efficient vector indexing, making lookups much faster than brute-force comparisons.

Second, it offers persistent and scalable storage, allowing us to store large volumes of embeddings—far beyond the capacity of a single machine.

Third, it supports concurrent access, which is essential in production environments where multiple users or services may query the system simultaneously.

Vector database is inheritely different from traditional database. Where traditional database uses row/columns to store data, vector databases uses High-dimensional vectors to store data.

##### Indexing in Vector Databases

Vector databases also use specialized indexing techniques to speed up similarity searches, such as:

- HNSW: [Hierarchical navigable small world](https://en.wikipedia.org/wiki/Hierarchical_navigable_small_world)
- IVF: [Inverted File Flat Vector](https://docs.oracle.com/en/database/oracle/oracle-database/26/vecse/understand-inverted-file-flat-vector-indexes.html)
- LSH: [Locality-sensitive hashing](https://en.wikipedia.org/wiki/Locality-sensitive_hashing)

##### Popular Vector Databases

- PGVector: Open-source vector similarity search for PostgreSQL
- sqlite-vec: Open-source vector similarity search for SQLite
- LanceDB: Local-first, simple setup, small–medium scale
- Weaviate: Full-featured, GraphQL API, complex schema

#### Chunking

f we create an embedding for a short paragraph, it works well. However, when we try to embed a long paragraph into a single embedding, several issues arise:

- Semantic dilution: Combining multiple topics into one embedding causes the overall meaning to become less precise.
- Token limit: The model can only process a limited number of tokens.
- Reduced precision: Specific concepts tend to get "averaged out," weakening their distinct representation.
- Irrelevant matches: Certain parts of the long text may lead to poor or unrelated search matches.

Chunking solve this problem by splitting long documents into smaller pieces.

##### Fixed-Size Chunking

The simplest chunking approach is to split the text into fixed-size pieces based on character count or word count.

- Predictable size: All chunks are roughly the same length
- Simple implementation: Easy to understand and implement
- Fast runtime performance: Fast chunking with minimal computation
- Token control: Can ensure chunks fit model limits

However, fixed-size chunking has one major drawback — a paragraph can be split at arbitrary points, leading to a loss of context.

To address this, we use a technique called chunk overlap, where consecutive chunks share a portion of text to maintain continuity.

The optimal overlap size depends on the nature of the data and should be verified experimentally, but a 20% overlap is a good starting point.

##### Semantic Chunking

In semantic chunking we try to break text according to natural language breaks like sentance or paragraph. Each chunk will contain complete throught. We can still use overlap to increase context.

So in a typical flow of large text or document we want to:

1. Chunk document text based on semantic or natural language breaks. This can be done using a regex like `r"(?<=[.!?])\s+"`
2. Each chunk is passed to model to create an embedding for the chunk
3. Create a relationship between embedding and parent document
4. Search user query by first converting it into a embedding from same model and then query each embedding individually

##### ColBERT

[ColBERT](https://arxiv.org/abs/2004.12832) creates one embedding per word, with each word contextualized.

ColBERT is an example of MVR (Multi-Vector Retrieval) where a document or chunk is represented by multiple vectors.

It require more storage and computational power

##### Late Chunking

[Late Chunking](https://jina.ai/news/late-chunking-in-long-context-embedding-models/) creates a single embedding for the enttire document (or as much of it as possible), and then uses that embedding to create context-aware embeddings for each word.

Now, each word contributes more meaningful information to the final embedding because its role in the text is already understood.

Again, this is more complex and require more computational power like [ColBERT](#ColBERT). Most real world scenarios can be handled by regular chunking.

### Hybrid Search

There are two different types of search

1. [Keyword Search](#keyword-search)
2. [Semantic Search](#semantic-search)

Both searches have thier own strength. For example when some one searching for exact query or exact year or exact time then [Keyword Search](#keyword-search) will perform better. However, if someone searching for a theme or concept then [Semantic Search](#semantic-search).

So which one to use. We probably need to use both search and combine their result to give user more accurate results.

Their are two major ways by which we can combine the result of both the searches. They are:

1. Weighted Combination
2. Reciprocal Rank Fusion

#### Weighted Combination

We cannot combine the score of keyword and semantic search because both the scores are on different scales alltogether.

- Keyword search is on `0-100+`
- Scmantic search is on cosine score `0-1`

There are two steps to combine the score:

1. Normalize both, Keyword search & Semantic search scores
2. Combine with weighted score

##### Min-Max Normaliztion

The simplest way to normalize the two score is using Min-Max normalization.
Following is the formula:

```
normalized_score = (score - min_score) / (max_score - min_score)
```

for example for the input `[23.2, 8.7, 2.1]`,

- `min_score = 2.1` and
- `max_score = 23.2`

Putting on formula:

```
(23.2 - 2.1) / (23.2 - 2.1) = 1
(8.7 - 2.1) / (23.2 - 2.1) = .313
(2.1 - 2.1) / (23.2 - 2.1) = 0.00
```

Similarly for semantic score `[.623, .453, .231]`

- `min_score = .231` &
- `max_score = .623`

```
(.623 - .231) / (.623 - .231) = 1
(.453 - .231) / (.623 - .231) = .566
(.231 - .231) / (.623 - .231) = 0.00
```

##### Combination

Formula

```
score = (alpha * keyword_score) + ((1-alpha) * semantic_score)
```

Alpha is just a constant that we can use to dynamically control the weighting between the two scores

Alpha range between `0-1`.

- When `alpha=1`, 100% weighting to keyword score
- When `alpha=0`, 100% weighting to semantic score
- When `alpha=0.5`, equal weighting to both

#### Reciprocal Rank Fussion

The issue with [Min-Max Normalization](#min-max-normaliztion) is that it does not work well with major outlier or both search result has different kinds of distributions.

The other way is to use Reciprocal Rank Fussion. In RRF, we completely skip scores of semantic and keyword search and only consider thier ranks.

```
score = 1 / (k + rank)
```

Where `k` is a constant to controls how distribution of score between higher ranked vs lower ranked once get effected. Example:

- Lower `k` values like `20` gives more weight to top ranked results, creating a steep drop-off in scores
- Higher `k` values like `100` creates a more gradual decline.

A good default for k is `60`

### LLM in Search

#### LLM for improving query

User enter query which is not formated or has typos or it is just a framgment of thought. We can use LLM to user enhance query before we perform hybrid search.

Using LLM we can do following with query:

- Fixes typos
- Expands meaning
- Breaks apart complex queries
- Add missing context

#### LLM for Reranking query

To be extremly good at search, we go through 3 stages of searching

1. Stage 1: Fast BM25/Cosine similarity to find top 25 results like using [Hybrid Search](#hybrid-search)
   1.1 [Keyword Search](#keyword-search)
   1.2 [Semantic Search](#semantic-search)
   1.3 Combine result using [Reciprocal Rank Fussion](#reciprocal-rank-fussion) or [Weighted Combination](#weighted-combination)
2. Stage 2: Slow re-ranking to find the top 5 out of 25 results

For reranking we using query and full document to rerank the results. It's more accurate but it's much slower.

There are two main ways we can do Stage two:

1. Use LLM to re-rank. We pass the top results from stage 1 in batch so that LLM can rank them relative to each other
2. [Cross-Encoder Re-Ranking](https://sbert.net/examples/cross_encoder/applications/README.html):

##### Cross-Encoder Re-Ranking

Out semantic search embeddings were created with bi-encoder, which embeds queries and document seprarately.

In cross-encoder, we embed the query and document together in a single input. Output is the score and we don't need to calculate cosine similarity as cross-encoders see the full context of how the query and document interact, so they can catch subtle relationships that bi-encoders miss.

Advantage is that cross-encoder are much faster and cheaper than LLMs. Cross-encoder is a regression model. A regression model is a type of machine learning model that predicts a number / predict a continuous value — not a category or label.

Second advantage of cross-encoder is that they can be fine-tuned on your specific domain relatively easily.

Cohere API uses cross-encoders.

### Evaluation

#### Manual Evaluation

Ask following question

- What's here that shouldn't be?
- What's not here that should be?
- Would I click on these results?
- Did the result give you enough information to know whether the movies are relevant?
- Would it be better to return fewer results because the last few usually aren't relevant?
- Would it be better to return more results because there are more highly relevant queries?

#### Automatic Evaluation

To evaluate result you need a curated dataset, also called "Golden Dataset". Usually created with the help of domain expert. It's like your test data for your test cases.

Ideally it should contain

- `query`: Actual user query
- Expected result

##### Precision Metrics

One of the metric which we check with evaluation is precision metrics. We try to anser the question

**Among the results I recived, How many are revelvant?**
Formula:

```
precision = relevant_retrieved / total_retrieved
```

- Higher precision = less junk in the results.

[Precision@K or P@K](<https://en.wikipedia.org/wiki/Evaluation_measures_(information_retrieval)#Precision_at_k>) is a common metric that measure precision of top k results

##### Recall Metrics

Recall measures the completeness of result. It's like asking, **did we get all the results which we expect to get?**

```
recall = relevant_retrieved / total_relevant
```

Recall is critical, it defines the user experience. User experience decresase if they don't find what they expect to find.

Also, in critical applications like medical or legal discovery & safety systems need high recall as missing information is critical.

Some times [Precision Metrics](#precision-metrics) & recall can not be optimized together. In that case one need to do tradeoff.

##### F1 Score

[F1 Score](https://en.wikipedia.org/wiki/F-score) is a single metric that combine [percision](#precision-metrics) & [recall](#recall-metrics).

Formula:

```
f1 = 2 * (precision * recall) / (precision + recall)x
```

F1 score is the [harmonic mean](https://en.wikipedia.org/wiki/Harmonic_mean) of precision and recall. It gives you one number that represents the overall performance of your search system.

Harmonic mean is better than a regular mean because it punishes extreme disparities more. F1 is better when:

- Precision and recall are equally important
- You want a single metric to optimize
- You're comparing different systems

#### LLM Evaluation

Manual evaluation is slow. Automated metrics miss nuance. So we can add LLM based evalution.

LLM is not a domain export but with the help of domain expert we can define what success look like. We need to

- define clear evaluation criteria
- Specify what makes a result relevant
- Artivulate your quality standards

LLM might not be as good as human domain expert but still better then developer who might not have the domain knowledge. Another advantage is that we can implement LLM for scale and speed.

**Implementation Strategy**
Start with experts – Define clear evaluation criteria
Create detailed prompts – Include domain knowledge
Validate on samples – Check LLM agrees with experts
Use for scale – Let LLM handle bulk evaluation
Spot-check results – Have experts review surprising scores
