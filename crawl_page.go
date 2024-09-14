package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	parsedCurrent, err := url.Parse(rawCurrentURL)
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
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
	fmt.Printf("Grabbing content from %s...\n", normalizedURL)
	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		// TODO: find a better way to handle this edge case. better to not error out??
		if strings.Contains(err.Error(), "application/xml") {
			fmt.Println("found xml page, skipping...")
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("Grabbing links from %s...\n", normalizedURL)
	pageURLs, err := cfg.getURLsFromHTML(pageHTML)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Crawling links...")
	for _, page := range pageURLs {
		cfg.wg.Add(1)
		go cfg.crawlPage(page)

	}
	fmt.Printf("Finished with %s!", normalizedURL)
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if cfg.pages[normalizedURL] > 0 {
		isFirst = false
		cfg.pages[normalizedURL] += 1
	} else {
		isFirst = true
		cfg.pages[normalizedURL] = 1
	}
	return isFirst
}
