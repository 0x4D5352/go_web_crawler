package main

import (
	"errors"
	"golang.org/x/net/html"
	"log"
	"net/url"
	"strings"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	if len(htmlBody) == 0 {
		return []string{}, errors.New("Error: HTML missing body")
	}
	htmlReader := strings.NewReader(htmlBody)
	htmlNode, err := html.Parse(htmlReader)
	if err != nil {
		log.Fatal(err)
	}
	var urls []string
	var traverseTree func(*html.Node)
	traverseTree = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attribute := range node.Attr {
				if attribute.Key == "href" {
					newURL, err := url.Parse(attribute.Val)
					if err != nil {
						log.Fatal(err)
					}
					var fullURL string
					if newURL.Host != "" {
						fullURL = newURL.Scheme + "://" + newURL.Host + newURL.Path
					} else {
						fullURL = rawBaseURL
						if strings.Split(newURL.Path, "")[0] != "/" {
							fullURL += "/"
						}
						fullURL += newURL.Path
					}
					urls = append(urls, fullURL)
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
