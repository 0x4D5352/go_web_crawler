package main

import (
	"fmt"
	"net/url"
	"strings"
)

// Strips out all unncessary information from a given URL and returns the host and full path in lowercase.
// e.g. https://blog.example.com/foo/bar?query=true will become "blog.example.com/foo/bar"
func normalizeURL(rawURL string) (string, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL: %w", err)
	}
	fullPath := url.Host + strings.TrimSuffix(url.Path, "/")
	fullPath = strings.ToLower(fullPath)
	return fullPath, nil
}
