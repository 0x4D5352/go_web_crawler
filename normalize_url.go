package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return url.Host + strings.TrimSuffix(url.Path, "/"), nil
}
