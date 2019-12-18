// https://www.devdungeon.com/content/web-scraping-go#parse_urls
package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	// Parse a complex URL
	complexURL := "https://www.example.com/path/to/?query=123&this=that#fragment"
	parsedURL, err := url.Parse(complexURL)
	if err != nil {
		log.Fatal(err)
	}

	// Print out URL pieces
	fmt.Println("Scheme: " + parsedURL.Scheme)
	fmt.Println("Host: " + parsedURL.Host)
	fmt.Println("Path: " + parsedURL.Path)
	fmt.Println("Query string: " + parsedURL.RawQuery)
	fmt.Println("Fragment: " + parsedURL.Fragment)

	// Get the query key/values as a map
	fmt.Println("\nQuery values:")
	queryMap := parsedURL.Query()
	fmt.Println(queryMap)

	// Craft a new URL from scratch
	var customURL url.URL
	customURL.Scheme = "https"
	customURL.Host = "google.com"
	newQueryValues := customURL.Query()
	newQueryValues.Set("key1", "value1")
	newQueryValues.Set("key2", "value2")
	customURL.Fragment = "bookmarkLink"
	customURL.RawQuery = newQueryValues.Encode()

	fmt.Println("\nCustom URL:")
	fmt.Println(customURL.String())
}
