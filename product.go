package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type ProductInfo struct {
	Category string
	Name     string
	Price    string
	Sizes    []string
}

func fetchProductInfo(url string) ProductInfo {
	var productInfo ProductInfo

	// Fetch URL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	// Load HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Error reading HTML document: %v", err)
	}

	// Extract category
	productInfo.Category = doc.Find(".categoryName").Text()

	// Extract product name
	productInfo.Name = doc.Find(".itemTitle").Text()

	// Extract product price
	productInfo.Price = doc.Find(".price-value").Text()

	// Extract available sizes
	doc.Find(".sizeSelectorListItemButton").Each(func(i int, s *goquery.Selection) {
		size := s.Text()
		if size != "disable" {
			productInfo.Sizes = append(productInfo.Sizes, size)
		}
	})

	return productInfo
}
