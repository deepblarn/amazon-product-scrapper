package main

import (
	"fmt"
	"strconv"
	"github.com/gocolly/colly"
	"encoding/json"
	"io/ioutil"
	"os"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	crawlUrl := " "

	file,_ := os.Open("categories.json")
	byteValue, _ := ioutil.ReadAll(file)
	var categoriesUrls []string
	json.Unmarshal(byteValue, &categoriesUrls)

	fmt.Println()

	c := colly.NewCollector()
	activePage := 0

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML(".a-selected a", func(e *colly.HTMLElement) {
		fmt.Println("PAGE NUMBER-------------- " + e.Text)
		currentPage, _ := strconv.Atoi(e.Text)
		activePage = currentPage
	})

	c.OnHTML("div.s-result-list.s-search-results.sg-row", func(e *colly.HTMLElement) {
		e.ForEach("div.a-section.a-spacing-medium", func(_ int, e *colly.HTMLElement) {
			var productName string

			productName = e.ChildText("span.a-size-medium.a-color-base.a-text-normal")

			if productName == "" {
				// If we can't get any name, we return and go directly to the next element
				return
			}

			fmt.Printf("Product Name: %s \n", productName)
		})

		pageNumber := strconv.Itoa((activePage+1))
		c.Visit(crawlUrl + "&page=" + pageNumber)

	})

	for _, element := range categoriesUrls {
		crawlUrl = element+"&__mk_es_ES=ÅMÅŽÕÑ"
		c.Visit(element+"&__mk_es_ES=ÅMÅŽÕÑ")
	}

}