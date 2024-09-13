package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("no website provided")
	}
	if len(args) > 1 {
		log.Fatal("too many arguments provided")
	}

	rawURL := args[0]
	fmt.Printf("starting crawl of: %s\n", rawURL)
	pages := make(map[string]int)
	crawlPage(rawURL, rawURL, pages)
	for key, value := range pages {
		fmt.Printf("%s: seen %d times.\n", key, value)
	}
}
