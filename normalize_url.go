package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL: %w", err)
	}
	fullPath := url.Host + strings.TrimSuffix(url.Path, "/")
	fullPath = strings.ToLower(fullPath)
	return fullPath, nil
}
