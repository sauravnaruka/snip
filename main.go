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

	pathToDataFile := os.Getenv("DB_FILE_PATH")
	if pathToDataFile == "" {
		log.Fatal("DB_FILE_PATH must be set in env")
	}

	searchClient := search.NewClient(pathToDataFile)

	rootCmd := cmd.GetRootCommand(searchClient)
	startRepl(rootCmd)
}
