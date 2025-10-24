package search

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type InvertedIndex struct {
	// Index maps tokens to sets of document IDs
	Index map[string]map[int]struct{}

	// DocMap maps document IDs to their full Movie objects
	DocMap map[int]Movie
}

func newInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Index:  make(map[string]map[int]struct{}),
		DocMap: make(map[int]Movie),
	}
}

func (c *Client) BuildInvertedIndex() error {
	for _, movie := range c.movies {
		text := movie.Title + " " + movie.Description
		tokens := preProcessQuery(text, c.stopWordMap)
		c.invertedIndex.addDocument(movie.Id, tokens, movie)
	}

	return c.invertedIndex.save()
}

func (idx *InvertedIndex) addDocument(docID int, tokens []string, movie Movie) {
	uniqueTokens := make(map[string]struct{})
	for _, token := range tokens {
		uniqueTokens[token] = struct{}{}
	}

	for token := range uniqueTokens {
		if idx.Index[token] == nil {
			idx.Index[token] = make(map[int]struct{})
		}

		idx.Index[token][docID] = struct{}{}
	}

	idx.DocMap[docID] = movie
}

func (idx *InvertedIndex) getDocumentIDs(token string) []int {
	docIDSet, ok := idx.Index[token]
	if !ok {
		return nil
	}

	docIDs := make([]int, 0, len(docIDSet))
	for id := range docIDSet {
		docIDs = append(docIDs, id)
	}

	sort.Ints(docIDs)
	return docIDs
}

func (idx *InvertedIndex) getMovieByID(id int) (Movie, bool) {
	movie, ok := idx.DocMap[id]
	return movie, ok
}

func (idx *InvertedIndex) save() error {
	return idx.saveToPath("cache/inverted_index.gob")
}

func (idx *InvertedIndex) saveToPath(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(idx.Index); err != nil {
		return fmt.Errorf("failed to encode index: %w", err)
	}

	if err := encoder.Encode(idx.DocMap); err != nil {
		return fmt.Errorf("failed to encode docmap: %w", err)
	}

	return nil
}

func (idx *InvertedIndex) load() error {
	return idx.loadFromPath("cache/inverted_index.gob")
}

func (idx *InvertedIndex) loadFromPath(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	if err := decoder.Decode(&idx.Index); err != nil {
		return fmt.Errorf("failed to decode index: %w", err)
	}

	if err := decoder.Decode(&idx.DocMap); err != nil {
		return fmt.Errorf("failed to decode DocMap: %w", err)
	}

	return nil
}
