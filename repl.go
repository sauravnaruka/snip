package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func startRepl(rootCmd *cobra.Command) {
	reader := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to Snip REPL!")
	fmt.Println("Type 'help' to see available commands or 'exit'/'quit' to leave.")
	fmt.Println()

	for {
		fmt.Print("Snip >")
		if !reader.Scan() {
			break
		}

		args := cleanInput(reader.Text())
		if len(args) == 0 {
			continue
		}

		if args[0] == "exit" || args[0] == "quit" {
			fmt.Println("Exiting Snip!")
			break
		}

		rootCmd.SetArgs(args)
		rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
			f.Changed = false
		})

		if err := rootCmd.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		fmt.Println()
	}

	if err := reader.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return []string{}
	}

	words := strings.Fields(trimmed)

	if len(words) > 0 {
		words[0] = strings.ToLower(words[0])
	}

	return words
}
