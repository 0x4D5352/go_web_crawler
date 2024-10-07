package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type RobotsTxt struct {
	source     *url.URL
	allowed    []string
	disallowed []string
}

func (cfg *config) checkRobotsTxt(rawURL string) {
	cleanURL, err := normalizeURL(rawURL)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := http.Get("https://" + cleanURL + "/robots.txt")
	if err != nil {
		// TODO: handle non-404s
		return
	}
	defer robots.Body.Close()
	contents, err := io.ReadAll(robots.Body)
	if err != nil || len(contents) == 0 {
		// TODO: make this smarter
		return
	}
	if string(contents) == "User-agent: *" {
		fmt.Println("just there for show, go hog wild")
		return
	}
	fmt.Println("robots.txt exists, checking their requests...")
	parseRobotsTxt(&contents)
	source, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: fill out the allowed and disallowed lists kind of like the getlinksfromURL function
	// resolvedURL := cfg.baseURL.ResolveReference(href)
	cfg.robots = RobotsTxt{
		source: source,
	}
}

func parseRobotsTxt(contents *[]byte) {
	panic("unimplemented")
}
