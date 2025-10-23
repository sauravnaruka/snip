package search

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Client struct {
	movies      []Movie
	stopWordMap map[string]struct{}
}

type Movie struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type MovieData struct {
	Movies []Movie `json:"movies"`
}

func NewClient(dataFilePath string, stopWordsFilePath string) (*Client, error) {
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

	return &Client{
		movies:      movieData.Movies,
		stopWordMap: stopWordMap,
	}, nil
}
