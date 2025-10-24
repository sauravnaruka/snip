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

## Search

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

### TF-IDF

#### Inverted Index

An inverted index make search fast, it's like a SQL database index but for text. Instead of searching for all the documents each time a user searches, we build and index for fast lookup.

- Forward index maps location to value. eg: `doc1: [matrix, reload]`
- Inverted index maps value (token) to location. eg: `matrix: [doc1, doc5]`

We build inverted index in `internal/search` module.

1. We use struct `InvertedIndex` to store

- `Index` which maps tokens to set of document IDs
- `DocMap` which maps document IDs to the movie object

2. When `BuildInvertedIndex` called from cmd module it iterate through all the movies, create token from title and description, then call `addDocument`

3. `addDcoument` create a mapping of token and corresponding document id
