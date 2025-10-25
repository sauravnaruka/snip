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
