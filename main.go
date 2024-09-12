package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("no website provided")
	}
	if len(args) > 1 {
		log.Fatal("too many arguments provided")
	}

	rawURL := args[0]
	fmt.Printf("starting crawl of: %s", rawURL)
	body, err := getHTML(rawURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(body)
}

func getHTML(rawURL string) (string, error) {
	webpage, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("error - failed to fetch webpage for %s: %w", rawURL, err)
	}
	defer webpage.Body.Close()
	if webpage.StatusCode > 399 {
		return "", fmt.Errorf("error - website returned status code: %w", webpage.Status)
	}
	if len(webpage.Header["Content-Type"]) == 0 || len(webpage.Header["Content-Type"]) > 1 {
		return "", fmt.Errorf("error - incorrect number of content types: %d", len(webpage.Header["Content-Type"]))
	}
	if webpage.Header["Content-Type"][0] != "text/html" {
		return "", fmt.Errorf("error - incorrect content-type: %w", webpage.Header["Content-Type"])
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
