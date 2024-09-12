package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	webpage, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("error - failed to fetch webpage for %s: %w", rawURL, err)
	}
	defer webpage.Body.Close()
	if webpage.StatusCode > 399 {
		return "", fmt.Errorf("error - website returned status code: %w", webpage.Status)
	}
	contentType := webpage.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("got non-HTML response: %s", contentType)
	}
	body, err := io.ReadAll(webpage.Body)
	if err != nil {
		return "", fmt.Errorf("error - failed to read webpage body: %w", err)
	}
	if len(body) == 0 {
		return "", errors.New("error - webpage body was empty")
	}
	return string(body), nil
}
