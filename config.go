package main

import (
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
	robots             RobotsTxt
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
