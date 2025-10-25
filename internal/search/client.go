package search

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Client struct {
	movies        []Movie
	stopWordMap   map[string]struct{}
	invertedIndex *InvertedIndex
	BM25_K1       float64
}

func NewClient(dataFilePath string, stopWordsFilePath string, BM25_K1 float64) (*Client, error) {
	movieDataRaw, err := os.ReadFile(dataFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var movieData MovieData
	if err := json.Unmarshal(movieDataRaw, &movieData); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	stopWordsData, err := os.ReadFile(stopWordsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading stop words file: %w", err)
	}

	stopWords := strings.Fields(string(stopWordsData))
	stopWordMap := make(map[string]struct{}, len(stopWords))
	for _, w := range stopWords {
		stopWordMap[strings.ToLower(w)] = struct{}{}
	}

	client := &Client{
		movies:        movieData.Movies,
		stopWordMap:   stopWordMap,
		invertedIndex: newInvertedIndex(),
		BM25_K1:       BM25_K1,
	}
	client.invertedIndex.load()

	return client, nil
}
