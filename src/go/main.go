package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const mangahato_url = "https://mangahato.com"

func main() {
	manga := "manga-sakura-sakura-morishige-raw.html"
	url := mapURL(manga)
	saveDir := filepath.Join(".", baseFilename(manga))
	fmt.Println(saveDir)
	downloadManga(url, saveDir)
}

func baseFilename(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

func downloadManga(mangaURL string, saveDir string) {
	res, err := http.Get(mangaURL)
	check(err)
	defer res.Body.Close()

	// create goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(res.Body)
	check(err)

	mangaName := document.Find(".manga-info h1").First().Text()
	fmt.Println("Manga name:", mangaName)

	// find all chapter URLs
	s := document.Find("#tab-chapper tbody tr")
	fmt.Println("Number of chapters:", len(s.Nodes))

	chapterRelativeURLs := make([]string, len(s.Nodes))

	s.Each(func(index int, element *goquery.Selection) {
		href, exists := element.Find("td a").First().Attr("href")
		if exists {
			// fmt.Println(index, href)
			chapterRelativeURLs[index] = trimURL(href)
		}
	})

	// reverse array
	// https://stackoverflow.com/a/19239850/8364403
	for i, j := 0, len(chapterRelativeURLs)-1; i < j; i, j = i+1, j-1 {
		chapterRelativeURLs[i], chapterRelativeURLs[j] = chapterRelativeURLs[j], chapterRelativeURLs[i]
	}

	for i := 0; i < len(chapterRelativeURLs); i++ {
		fmt.Println(i, chapterRelativeURLs[i])
		chapterDir := filepath.Join(saveDir, padLeft(strconv.Itoa(i), 3, "0"))
		downloadChapter(mapURL(chapterRelativeURLs[i]), chapterDir)
		// break
	}
}

func trimURL(url string) string {
	return strings.Trim(url, "\n\r")
}

func downloadChapter(chapterURL string, saveDir string) {
	res, err := http.Get(chapterURL)
	check(err)
	defer res.Body.Close()

	// create goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(res.Body)
	check(err)

	// htmlString, err := document.Html()
	// check(err)
	// // fmt.Println(htmlString)

	// // write document to file
	// // https://gobyexample.com/writing-files
	// f, err := os.Create("./index.html")
	// check(err)

	// n, err := f.WriteString(htmlString)
	// check(err)

	// fmt.Println("n =", n)

	// defer f.Close()

	imgElements := document.Find("center .chapter-content img")
	imgURLs := make([]string, len(imgElements.Nodes))

	imgElements.Each(func(index int, imgElement *goquery.Selection) {
		imgURL, exists := imgElement.Attr("data-src")
		if exists {
			// trim leading/trailing \n\r
			// https://stackoverflow.com/a/55945601/8364403
			imgURLs[index] = trimURL((imgURL))
		} else {
			fmt.Println("data-src does not exit at index:", index)
		}
	})

	for i := 0; i < len(imgURLs); i++ {
		fmt.Println(i, imgURLs[i])

		if len(imgURLs[i]) > 0 {
			fileExt := filepath.Ext(imgURLs[i])

			err := os.MkdirAll(saveDir, os.ModePerm)
			check(err)

			imgPath := filepath.Join(saveDir, padLeft(strconv.Itoa(i), 3, "0")+fileExt)
			fmt.Println(imgPath)

			if fileExists(imgPath) {
				continue
			}

			imageBytes := requestImage(imgURLs[i])

			if imageBytes != nil {
				err := ioutil.WriteFile(imgPath, imageBytes, 0644)
				check(err)
			}
		}
		// break
	}

}

// https://golangcode.com/check-if-a-file-exists/
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func padLeft(s string, n int, c string) string {
	charsLeft := n - len(s)
	if charsLeft <= 0 {
		return s
	}

	return strings.Repeat(c, charsLeft) + s
}

func mapURL(documentFile string) string {
	return mangahato_url + "/" + documentFile
}

func check(e error) {
	if e != nil {
		// log.Fatal(e)
		fmt.Println(e)
	}
}

func requestImage(url string) []byte {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	check(err)

	// here is the important step
	req.Header.Set("referer", mangahato_url)
	res, err := client.Do(req)
	check(err)

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println(res.StatusCode, url)
		return nil
	}

	// fmt.Println("StatusCode", res.StatusCode)
	// fmt.Println("ResponseLength", res.ContentLength)

	imageContent, err := ioutil.ReadAll(res.Body)
	check(err)

	return imageContent
}
