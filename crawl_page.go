package main

import (
	"fmt"
	"log"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
	cfg.mu.Lock()
	maxPageLimitReached := len(cfg.pages) >= cfg.maxPages
	cfg.mu.Unlock()
	if maxPageLimitReached {
		fmt.Println("max page count reached! stopping crawl.")
		return
	}
	parsedCurrent, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.baseURL.Host != parsedCurrent.Host {
		fmt.Printf("%s is not part of %s, skipping!\n", parsedCurrent.Host, cfg.baseURL.Host)
		return
	}
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}
	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst {
		fmt.Printf("%s already visited, skipping!\n", normalizedURL)
		return
	}
	fmt.Printf("grabbing content from %s...\n", normalizedURL)
	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting page contents: %w", err)
		return
	}
	fmt.Printf("grabbing links from %s...\n", normalizedURL)
	pageURLs, err := cfg.getURLsFromHTML(pageHTML)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("crawling links...")
	for _, page := range pageURLs {
		cfg.wg.Add(1)
		go cfg.crawlPage(page)

	}
	fmt.Printf("finished with %s!\n", normalizedURL)
}
