package main

import (
	"fmt"
	"github.com/tebeka/selenium/chrome"
	"log"
	"strconv"

	"github.com/tebeka/selenium"
)

const (
	baseURL          = "https://shop.adidas.jp"
	chromeDriverPath = "./assets/chromedriver.exe"
)

var productIDs = []string{
	"IP0418",
	//"IY2911",
	//"II5763",
	//"IT2491",
	//"IZ4922",
}

func main() {
	port := 8088
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		log.Fatalf("Error starting the ChromeDriver server: %v", err)
	}
	defer service.Stop()

	wd := createWebDriver(port)
	defer wd.Quit()

	for _, productID := range productIDs {
		product := fetchProductInfo(wd, productID)
		if product.Info.Name == "" {
			log.Printf("Failed to fetch product info for ID %s", productID)
			continue
		}

		fmt.Printf("Product Info for ID: %s\n", productID)
		fmt.Println("Breadcrumbs:", product.Info.Breadcrumbs)
		fmt.Println("Category:", product.Info.Category)
		fmt.Println("Name:", product.Info.Name)
		fmt.Println("Price:", product.Info.Price)
		fmt.Println("Sizes:", product.Info.Sizes)
		fmt.Println()

		for _, coordinatedProduct := range product.Coordinates {
			fmt.Println("Product ID:", productID)
			fmt.Println("Coordinated Product Name:", coordinatedProduct.Name)
			fmt.Println("Pricing:", coordinatedProduct.Price)
			fmt.Println("Product Number:", coordinatedProduct.ProductNumber)
			fmt.Println("Image URL:", coordinatedProduct.ImageURL)
			fmt.Println("Product Page URL:", coordinatedProduct.ProductPageURL)
			fmt.Println()
		}
	}
}

func createWebDriver(port int) selenium.WebDriver {
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	caps.AddChrome(
		chrome.Capabilities{
			Path: "",
			Args: []string{
				"--window-size=1920,1080",
				"--headless",
				"--no-sandbox",
				"--blink-settings=imagesEnabled=false", // Disable images
				"--blink-settings=cssEnabled=false",    // Disable CSS
			},
		},
	)

	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(port)+"/wd/hub")
	if err != nil {
		log.Fatalf("Failed to create WebDriver: %v", err)
	}
	return wd
}
