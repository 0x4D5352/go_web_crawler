package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/url"
	"regexp"
	"strings"
)

// Pulls all URL links, both within the website and linked from other websites.
func (cfg *config) getURLsFromHTML(htmlBody string) ([]string, error) {
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
					resolvedURL := cfg.baseURL.ResolveReference(href)
					if cfg.robots.allowed == nil && cfg.robots.disallowed == nil {
						urls = append(urls, resolvedURL.String())
						break
					}
					newURL := resolvedURL.String()
					fmt.Println(newURL)
					if len(cfg.robots.disallowed) > 0 {
						for _, disallowed := range cfg.robots.disallowed {
							re := regexp.MustCompile(escapeQuestionMark(disallowed))
							found := re.MatchString(newURL)
							if err != nil {
								fmt.Printf("error when comparing %s to %s\n", disallowed, newURL)
								continue
							}
							if found {
								// fmt.Printf("not allowed to crawl %s, as it matches %s\n", newURL, disallowed)
								continue
							}
							// else it's not in the disallowed list
							// fmt.Println("safe to crawl?")
						}
					}
					if len(cfg.robots.allowed) > 0 {
						for _, allowed := range cfg.robots.allowed {
							re := regexp.MustCompile(escapeQuestionMark(allowed))
							found := re.MatchString(newURL)
							if err != nil {
								fmt.Printf("error when comparing %s to %s\n", allowed, newURL)
								continue
							}
							if !found {
								// fmt.Printf("not allowed to crawl %s, as it does not match %s\n", newURL, allowed)
								continue
							}
							// else it's in the allowed list
							// fmt.Println("safe to crawl!")
							urls = append(urls, newURL)
							break
						}
					}
					// fmt.Println("safe to crawl!")
					urls = append(urls, newURL)
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

// I cheated and used chatgpt on this, forgive me
func escapeQuestionMark(input string) string {
	// Replace only the '?' with an escaped version
	return strings.ReplaceAll(input, "?", "\\?")
}
