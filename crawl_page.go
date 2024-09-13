package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	parsedCurrent, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}
	parsedBase, err := url.Parse(rawBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if parsedBase.Host != parsedCurrent.Host {
		fmt.Printf("%s is not part of %s, skipping!\n", parsedCurrent.Host, parsedBase.Host)
		return
	}
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
	}
	if pages[normalizedURL] > 0 {
		fmt.Printf("%s already visited, skipping!\n", normalizedURL)
		pages[normalizedURL] += 1
		return
	} else {
		pages[normalizedURL] = 1
	}
	fmt.Printf("Grabbing content from %s...\n", normalizedURL)
	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		if strings.Contains(err.Error(), "application/xml") {
			fmt.Println("found xml page, skipping...")
			return
		}
		log.Fatal(err)
	}
	pageURLs, err := getURLsFromHTML(pageHTML, rawBaseURL)
	fmt.Printf("Grabbing links from %s...\n", normalizedURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Crawling links...")
	for _, page := range pageURLs {
		crawlPage(rawBaseURL, page, pages)
	}
}
