package search

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Client struct {
	dataFilePath *string
}

type Movie struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"genre"`
}

type MovieData struct {
	Movies []Movie `json:"movies"`
}

func NewClient(pathToDataFile string) *Client {
	c := Client{
		dataFilePath: &pathToDataFile,
	}

	return &c
}

func (c *Client) SearchMovie(query string) ([]Movie, error) {
	data, err := os.ReadFile(*c.dataFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var movieData MovieData
	if err := json.Unmarshal(data, &movieData); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	var results []Movie
	queryLower := strings.ToLower(query)

	for _, movie := range movieData.Movies {
		if strings.Contains(strings.ToLower(movie.Title), queryLower) {
			results = append(results, movie)
		}
	}

	if len(results) == 0 {
		fmt.Println("No movies found.")
	}

	for i, movie := range results {
		fmt.Printf("%d. %s\n", i+1, movie.Title)
	}

	return results, nil
}
