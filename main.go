package main

import (
	"log"

	"github.com/hsmtkk/selloptiondelta/cmd"
)

func main() {
	rootCmd := cmd.RootCommand
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
