package main

import (
	"fmt"
	"reflect"
)

func printReport(pages map[string]int, baseURL string) {
	fmt.Printf(`
=============================
  REPORT for %s
=============================
`, baseURL)
	sortedPages := sortPages(pages)
	for _, page := range sortedPages {
		fmt.Printf("Found %d internal links to %s\n", page.count, page.url)
	}
}

func sortPages(pages map[string]int) []page {
	var sortedPages []page
	for hostPath, count := range pages {
		newPage := page{
			count: count,
			url:   "https://" + hostPath,
		}
		sortedPages = append(sortedPages, newPage)
	}
	swapF := reflect.Swapper(sortedPages)
	for i, _ := range sortedPages {
		swapped := false
		for j, p := range sortedPages {
			if j >= len(sortedPages)-i-1 {
				break
			}
			if p.count < sortedPages[j+1].count {
				swapF(j, j+1)
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
	return sortedPages
}

type page struct {
	count int
	url   string
}

/*
	if len(sortedPages) == 0 || sortedPages[0].count >= count || count < sortedPages[len(sortedPages)-1].count {
		sortedPages = append(sortedPages, newPage)
	} else {
		sortedPages = append([]page{newPage}, sortedPages...)
	}
*/
