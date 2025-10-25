/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	BM25_K1 := 1.5
	BM25_K1_STR := os.Getenv("BM25_K1")
	if BM25_K1_STR == "" {
		fmt.Printf("BM25_K1 is empty in env file. Taking the default\n")
	} else {
		envBM25K1, err := strconv.ParseFloat(BM25_K1_STR, 64)
		if err != nil {
			fmt.Printf("Error converting env BM25_K1 value '%s' with error %v\n", BM25_K1_STR, err)
		} else {
			BM25_K1 = envBM25K1
		}
	}

	BM25_B := 0.75
	BM25_B_STR := os.Getenv("BM25_B")
	if BM25_B_STR == "" {
		fmt.Printf("BM25_B is empty in env file. Taking the default\n")
	} else {
		envBM25B, err := strconv.ParseFloat(BM25_B_STR, 64)
		if err != nil {
			fmt.Printf("Error converting env BM25_B value '%s' with error %v\n", BM25_B_STR, err)
		} else {
			BM25_B = envBM25B
		}
	}

	searchClient, err := search.NewClient(dataFilePath, stopWordsFilePath, BM25_K1, BM25_B)
	if err != nil {
		log.Fatal("Error creating search client %w", err)
	}

	rootCmd := cmd.GetRootCommand(searchClient)

	if len(os.Args) > 1 {
		// Run Cobra command directly
		rootCmd.SetArgs(os.Args[1:])
		if err := rootCmd.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	rootCmd.Help()

	// startRepl(rootCmd)
}
