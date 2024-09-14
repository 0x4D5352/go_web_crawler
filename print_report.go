package main

import "fmt"

func printReport(pages map[string]int, baseURL string) {
	sortedPages := sortPages(pages)
	for _, page := range sortedPages {
		fmt.Printf("Found %d internal links to %s\n", page.count, page.url)
	}
}

func sortPages(pages map[string]int) []page {
	sortedPages := make([]page, len(pages))
	for url, count := range pages {
		if len(sortedPages) == 0 {
			sortedPages := append(sortedPages, newPage)
	}
	return sortedPages
}

type page struct {
	count int
	url   string
}
			newPage := page{
				count: value,
				url:   "https://" + key,
			}
