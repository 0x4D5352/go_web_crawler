package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("no website provided")
	}
	if len(args) > 3 {
		log.Fatal("too many arguments provided")
	}

	rawURL := args[0]
	maxConcurrency := 5
	if len(args) > 1 {
		userConcurrency, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("error - unable to convert %s to integer: %w", args[1], err)
		}
		maxConcurrency = userConcurrency
	}
	maxPages := 50
	if len(args) > 2 {
		userPages, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("error - unable to convert %s to integer: %w", args[2], err)
		}
		maxPages = userPages
	}
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("invalid url")
	}
	fmt.Printf("starting crawl of: %s\n", rawURL)
	cfg := config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}
	fmt.Println("------------------")
	cfg.checkRobotsTxt(rawURL)
	cfg.wg.Add(1)
	cfg.crawlPage(rawURL)
	cfg.wg.Wait()
	fmt.Println("------------------")
	fmt.Println("crawl finished!")
	fmt.Println("------------------")
	printReport(cfg.pages, baseURL.String())
	fmt.Println("------------------")
}
