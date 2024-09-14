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
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("invalid url")
	}
	fmt.Printf("starting crawl of: %s\n", rawURL)
	cfg := config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 1),
		wg:                 &sync.WaitGroup{},
	}
	fmt.Println("------------------")
	cfg.crawlPage(rawURL)
	cfg.wg.Wait()
	fmt.Println("------------------")
	fmt.Println("crawl finished!")
	fmt.Println("------------------")
	for key, value := range cfg.pages {
		fmt.Printf("%s: seen %d times.\n", key, value)
	}
	fmt.Println("------------------")
}

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}
