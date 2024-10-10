package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type RobotsTxt struct {
	allowed    []string
	disallowed []string
	sitemap    []string
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
	allowed, disallowed, err := cfg.parseRobotsTxt(contents)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: fill out the allowed and disallowed lists kind of like the getlinksfromURL function
	cfg.robots = RobotsTxt{
		allowed:    allowed,
		disallowed: disallowed,
	}
	// fmt.Printf("Robots.txt:\n%+v\n", cfg.robots)
	fmt.Println("robots.txt parsed!")
}

func (cfg *config) parseRobotsTxt(contents []byte) ([]string, []string, error) {
	body := string(contents)
	lines := strings.Split(body, "\n")
	var allowed []string
	var disallowed []string
	var isOurUserAgent bool
	if strings.Contains(lines[len(lines)-2], "Sitemap:") {
		fmt.Println("sitemap found! just use that.")
		_, sitemap, _ := strings.Cut(lines[len(lines)-2], ": ")
		cfg.robots.sitemap = parseSitemap(sitemap)
	}
	for _, line := range lines {
		// fmt.Println(line)
		if len(line) == 0 {
			continue
		}
		if line == "User-agent: *" {
			isOurUserAgent = true
			continue
		} else if strings.Contains(line, "User-agent:") {
			isOurUserAgent = false
			continue
		}
		// else we're in an allow or disallow line
		if !isOurUserAgent {
			continue
		}
		group, path, exists := strings.Cut(line, ": ")
		if !exists {
			continue
		}
		href, err := url.Parse(path)
		if err != nil {
			return nil, nil, fmt.Errorf("couldn't parse href '%v': %v\n", path, err)
		}
		resolvedURL := cfg.baseURL.ResolveReference(href)
		switch group {
		case "Allow":
			allowed = append(allowed, resolvedURL.String())
		case "Disallow":
			disallowed = append(disallowed, resolvedURL.String())
		case "Sitemap":
			continue
		default:
			return nil, nil, fmt.Errorf("error: unexpected key/value pair %s", line)
		}
	}
	if allowed == nil && disallowed == nil {
		return nil, nil, fmt.Errorf("error: no valid groups in %s", lines)
	}
	return allowed, disallowed, nil
}

func parseSitemap(sitemap string) []string {
	return nil
}
