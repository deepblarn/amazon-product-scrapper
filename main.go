package main

import (
	"fmt"
	// "strconv"
	"github.com/gocolly/colly"
	"encoding/json"
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
	c := colly.NewCollector()
	categoriesUrls := []string{}
	crawlUrl := "https://www.amazon.es/s?k=nintendo+switch&ref=nb_sb_noss_1"

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("#departments", func(e *colly.HTMLElement) {
		e.ForEach("a.a-link-normal.s-navigation-item", func(_ int, e *colly.HTMLElement) {
			// fmt.Printf("CATEGORY: %s URL %s \n", e.ChildText("span.a-size-base.a-color-base"), e.Request.AbsoluteURL(e.Attr("href")))
			categoriesUrls = append(categoriesUrls,  e.Request.AbsoluteURL(e.Attr("href")))
			categoriesUrlsJSON, _ := json.Marshal(categoriesUrls)
			jsonFile,_ := os.Create("./categories.json")
			defer jsonFile.Close()
			jsonFile.Write(categoriesUrlsJSON)
			jsonFile.Close()
		})
	})
	
	c.Visit(crawlUrl)

	fmt.Println(categoriesUrls)
}