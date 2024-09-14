package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.wg.Add(1)
	defer cfg.wg.Done()
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
	// TODO: add addpagevisit as a bool conditional here
	fmt.Printf("Grabbing content from %s...\n", normalizedURL)
	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		if strings.Contains(err.Error(), "application/xml") {
			fmt.Println("found xml page, skipping...")
			return
		}
		log.Fatal(err)
	}
	pageURLs, err := cfg.getURLsFromHTML(pageHTML)
	fmt.Printf("Grabbing links from %s...\n", normalizedURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Crawling links...")
	for _, page := range pageURLs {
		go cfg.crawlPage(page)
		cfg.concurrencyControl <- struct{}{}
	}
	<-cfg.concurrencyControl
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if cfg.pages[normalizedURL] > 0 {
		fmt.Printf("%s already visited, skipping!\n", normalizedURL)
		isFirst = false
		cfg.pages[normalizedURL] += 1
	} else {
		isFirst = true
		cfg.pages[normalizedURL] = 1
	}
	return isFirst
}
