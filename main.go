/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"os"

	"github.com/SauravNaruka/snip/cmd"
	"github.com/SauravNaruka/snip/internal/search"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	dataFilePath := os.Getenv("DB_FILE_PATH")
	if dataFilePath == "" {
		log.Fatal("DB_FILE_PATH must be set in env")
	}

	stopWordsFilePath := os.Getenv("STOP_WORD_FILE_PATH")
	if stopWordsFilePath == "" {
		log.Fatal("STOP_WORD_FILE_PATH must be set in env")
	}

	searchClient, err := search.NewClient(dataFilePath, stopWordsFilePath)
	if err != nil {
		log.Fatal("Error creating search client %w", err)
	}

	rootCmd := cmd.GetRootCommand(searchClient)
	startRepl(rootCmd)
}
