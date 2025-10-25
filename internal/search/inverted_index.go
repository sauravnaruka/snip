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

	// Track document length
	DocLengthMap map[int]int

	// map[docId]map[token]frequency
	TermFrequencies map[int]map[string]int
}

func newInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Index:           make(map[string]map[int]struct{}),
		DocMap:          make(map[int]Movie),
		DocLengthMap:    make(map[int]int),
		TermFrequencies: make(map[int]map[string]int),
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

	termCounts := make(map[string]int)
	for _, token := range tokens {
		termCounts[token]++
	}
	idx.TermFrequencies[docID] = termCounts

	// termCounts contain unique tokens as keys
	for token := range termCounts {
		if idx.Index[token] == nil {
			idx.Index[token] = make(map[int]struct{})
		}
		idx.Index[token][docID] = struct{}{}
	}

	idx.DocMap[docID] = movie
	idx.DocLengthMap[docID] = len(tokens)
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

	if err := encoder.Encode(idx.TermFrequencies); err != nil {
		return fmt.Errorf("failed to encode TermFrequencies: %w", err)
	}

	if err := encoder.Encode(idx.DocLengthMap); err != nil {
		return fmt.Errorf("failed to encode DocLengthMap: %w", err)
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

	if err := decoder.Decode(&idx.TermFrequencies); err != nil {
		return fmt.Errorf("failed to decode TermFrequencies: %w", err)
	}

	if err := decoder.Decode(&idx.DocLengthMap); err != nil {
		return fmt.Errorf("failed to decode DocLengthMap: %w", err)
	}

	return nil
}

func (idx *InvertedIndex) getTF(docID int, token string) int {
	if freqCounter, ok := idx.TermFrequencies[docID]; ok {
		// Returns 0 if token not found
		return freqCounter[token]
	}
	return 0
}

func (idx *InvertedIndex) getAvgDocLength() float64 {
	totalDocuments := len(idx.DocLengthMap)
	if totalDocuments == 0 {
		return 0.0
	}

	totalLength := 0
	for _, docLength := range idx.DocLengthMap {
		totalLength += docLength
	}

	return float64(totalLength) / float64(totalDocuments)
}
