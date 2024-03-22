package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchBreadcrumbString(url string) string {
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

	// Remove all SVG elements from the document
	doc.Find("svg").Remove()

	// Find all breadcrumb list items and concatenate their text content with "/"
	var breadcrumbs []string
	doc.Find("li.breadcrumbListItem").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if text != "" {
			breadcrumbs = append(breadcrumbs, text)
		}
	})

	// Concatenate breadcrumb items with "/"
	breadcrumbString := strings.Join(breadcrumbs, "/")

	return breadcrumbString
}
