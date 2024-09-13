package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
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

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}
