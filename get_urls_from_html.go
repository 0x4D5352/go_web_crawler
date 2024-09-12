package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	if len(htmlBody) == 0 {
		return []string{}, errors.New("Error: HTML missing body")
	}
	htmlReader := strings.NewReader(htmlBody)
	htmlNode, err := html.Parse(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("Error parsing HTML: %w", err)
	}
	var urls []string
	var traverseTree func(*html.Node)
	traverseTree = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attribute := range node.Attr {
				if attribute.Key == "href" {
					href, err := url.Parse(attribute.Val)
					if err != nil {
						fmt.Printf("couldn't parse href '%v': %v\n", attribute.Val, err)
						continue
					}

					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
					break
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverseTree(child)
		}
	}
	traverseTree(htmlNode)

	return urls, nil
}
