package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/gocolly/colly"
)

var crawlUrl = "https://www.amazon.es/s?k=nintendo+switch&ref=nb_sb_noss_1"
var categoriesUrls = []string{}
var productLinks = []string{}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	fetchCategories()
	visitAllPagesFromCategories()
	fetchItemInfo()
}

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func fetchCategories() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	c.OnHTML("#departments", func(e *colly.HTMLElement) {
		e.ForEach("a.a-link-normal.s-navigation-item", func(_ int, e *colly.HTMLElement) {
			categoriesUrls = append(categoriesUrls, e.Request.AbsoluteURL(e.Attr("href")))
		})
	})

	c.Visit(crawlUrl)
}

func visitAllPagesFromCategories() {
	c := colly.NewCollector()
	activePage := 0
	c.Limit(&colly.LimitRule{Parallelism: 10})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	c.OnHTML(".a-selected a", func(e *colly.HTMLElement) {
		fmt.Println("PAGE NUMBER-------------- " + e.Text)
		currentPage, _ := strconv.Atoi(e.Text)
		activePage = currentPage
	})

	c.OnHTML("div.s-result-list.s-search-results.sg-row", func(e *colly.HTMLElement) {
		e.ForEach("div.a-section.a-spacing-medium", func(_ int, e *colly.HTMLElement) {
			var productLink string
			e.ForEach("a.a-link-normal.a-text-normal", func(_ int, e *colly.HTMLElement) {
				productLink = e.Request.AbsoluteURL(e.Attr("href"))
				if productLink == "" {
					return
				}
				productLinks = append(productLinks, productLink)
			})
		})

		pageNumber := strconv.Itoa((activePage + 1))
		c.Visit(crawlUrl + "&page=" + pageNumber)
	})

	for _, element := range categoriesUrls {
		crawlUrl = element
		c.Visit(element)
	}
}

func fetchItemInfo() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	c.OnHTML("span.priceBlockDealPriceString", func(e *colly.HTMLElement) {
		fmt.Println("PRICE-------------- " + e.Text)
	})

	c.OnHTML("span#productTitle", func(e *colly.HTMLElement) {
		fmt.Println("PRODUCT-------------- " + e.Text)
	})

	for _, element := range productLinks {
		fmt.Println("--------------------")
		c.Visit(element)
		fmt.Println("--------------------")
	}
}
