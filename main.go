package main

import (
	"fmt"
	"log"
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
