package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	baseURL := "https://shop.adidas.jp"
	productID := "II5763"
	url := fmt.Sprintf("%s/products/%s/", baseURL, productID)
	productInfo := fetchProductInfo(url, baseURL)
	fmt.Println("Category:", productInfo.Category)
	fmt.Println("Product Name:", productInfo.Name)
	fmt.Println("Pricing:", productInfo.Price)
	fmt.Println("Available Sizes:", strings.Join(productInfo.Sizes, ", "))
	fmt.Println("Image URLs:")
	for _, imageURL := range productInfo.ImageURLs {
		fmt.Println(imageURL)
	}
	fmt.Println("Sense of Size:", productInfo.SenseOfSize)
}

type ProductInfo struct {
	Category    string
	Name        string
	Price       string
	Sizes       []string
	ImageURLs   []string
	SenseOfSize string
}

type CoordinateProductInfo struct {
	Items []CoordinateItem `json:"items"`
}

type CoordinateItem struct {
	ArticleID string   `json:"article_id"`
	Price     string   `json:"price"`
	ImageURL  string   `json:"image_url"`
	Sizes     []string `json:"sizes"`
}

func fetchProductInfo(url string, baseURL string) ProductInfo {
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

	// Extract image URLs
	doc.Find(".selectableImageListItem").Each(func(i int, s *goquery.Selection) {
		imageURL, _ := s.Find("img").Attr("src")
		if imageURL != "" {
			imageURL = baseURL + imageURL
			productInfo.ImageURLs = append(productInfo.ImageURLs, imageURL)
		}
	})

	// Extract sense of size
	productInfo.SenseOfSize = doc.Find(".sizeFitBar .label").Text()

	return productInfo
}

func fetchCoordinateProductInfo1(url string, baseURL string) CoordinateProductInfo {
	var coordinateInfo CoordinateProductInfo

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

	// Extract coordinate items
	doc.Find(".coordinate_item_tile").Each(func(i int, s *goquery.Selection) {
		articleID, _ := s.Attr("data-articleid")
		price, _ := s.Attr("data-price")
		imageURL, _ := s.Find("img").Attr("src")
		sizes := make([]string, 0)
		s.Find(".textItemButton-text").Each(func(i int, s *goquery.Selection) {
			sizes = append(sizes, s.Text())
		})

		coordinateInfo.Items = append(coordinateInfo.Items, CoordinateItem{
			ArticleID: articleID,
			Price:     price,
			ImageURL:  baseURL + imageURL,
			Sizes:     sizes,
		})
	})

	return coordinateInfo
}
func fetchCoordinatedProductInfo(url string, baseURL string) ProductInfo {
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

	// Find the carousel list items and simulate clicking on the first item
	firstCarouselItem := doc.Find(".carouselListitem").First()
	if firstCarouselItem.Length() > 0 {
		// Get the data-articleid attribute of the first item
		articleID, _ := firstCarouselItem.Find(".coordinate_item_tile").Attr("data-articleid")

		// Construct the URL for the coordinated product item
		coordinatedProductURL := baseURL + "/products/" + articleID

		log.Printf("Coordinated product URL: %s", coordinatedProductURL)

		// Fetch the coordinated product item page
		coordinatedResp, err := http.Get(coordinatedProductURL)
		if err != nil {
			log.Fatalf("Error fetching coordinated product URL: %v", err)
		}
		defer coordinatedResp.Body.Close()

		// Load the HTML document of the coordinated product item page
		coordinatedDoc, err := goquery.NewDocumentFromReader(coordinatedResp.Body)
		if err != nil {
			log.Fatalf("Error reading coordinated product HTML document: %v", err)
		}

		// Extract coordinated product information
		// Modify this part according to the structure of the coordinated product page
		log.Printf("Coordinated product page title: %s", coordinatedDoc.Find("title").Text())

		// Add extraction logic here...

		// For debugging, log the extracted information
		log.Printf("Extracted coordinated product info: %+v", productInfo)
	}

	return productInfo
}
