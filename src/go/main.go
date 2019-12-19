package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const BASE_URL = "https://mangahato.com"

func main() {
	url := BASE_URL + "/manga-sakura-sakura-morishige-raw.html"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// create goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	mangaName := document.Find(".manga-info h1").First().Text()
	fmt.Println("Manga name:", mangaName)

	// find all chapter URLs
	s := document.Find("#tab-chapper tbody tr")
	fmt.Println("Number of chapters:", len(s.Nodes))
	s.Each(func(index int, element *goquery.Selection) {
		href, exists := element.Find("td a").First().Attr("href")
		if exists {
			fmt.Println(index, href)
		}
	})
}
