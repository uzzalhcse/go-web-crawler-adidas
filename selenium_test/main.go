package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
)

func main() {
	// Start a Selenium WebDriver server instance (if one is not already running).
	const (
		seleniumPath    = "assets/selenium-server-standalone-3.4.jar"
		geckoDriverPath = "assets/geckodriver-v0.18.0-linux64"
		port            = 4444
	)
	//opts := []selenium.ServiceOption{
	//	selenium.StartFrameBuffer(),
	//	selenium.GeckoDriver(geckoDriverPath),
	//	selenium.Output(os.Stderr),
	//}
	//selenium.SetDebug(true)
	//service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	//if err != nil {
	//	panic(err)
	//}
	//defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	baseURL := "https://shop.adidas.jp"
	productID := "II5763"
	url := fmt.Sprintf("%s/products/%s/", baseURL, productID)
	// Navigate to the Adidas website where coordinated product information is available.
	if err := wd.Get(url); err != nil {
		panic(err)
	}

	// Wait for the carousel items to appear
	elem, err := wd.FindElement(selenium.ByCSSSelector, ".carouselListitem")
	if err != nil {
		log.Fatalf("Failed to find carousel items: %v", err)
	}

	fmt.Println("Found carousel items:", elem)

	// Now you can proceed with coordinated product information extraction logic...
	coordinatedProductInfo := fetchCoordinatedProductInfo(wd, baseURL)
	fmt.Println("Coordinated Product Name:", coordinatedProductInfo.Name)
	fmt.Println("Pricing:", coordinatedProductInfo.Price)
	fmt.Println("Product Number:", coordinatedProductInfo.ProductNumber)
	fmt.Println("Image URL:", coordinatedProductInfo.ImageURL)
	fmt.Println("Product Page URL:", coordinatedProductInfo.ProductPageURL)
}

type CoordinatedProductInfo struct {
	Name           string
	Price          string
	ProductNumber  string
	ImageURL       string
	ProductPageURL string
}

func fetchCoordinatedProductInfo(wd selenium.WebDriver, baseURL string) CoordinatedProductInfo {
	var coordinatedProductInfo CoordinatedProductInfo

	// Extract coordinated product name
	nameElement, err := wd.FindElement(selenium.ByCSSSelector, ".coordinate_item_container .title")
	if err != nil {
		log.Fatalf("Failed to find coordinated product name: %v", err)
	}
	name, err := nameElement.Text()
	if err != nil {
		log.Fatalf("Failed to get coordinated product name text: %v", err)
	}
	coordinatedProductInfo.Name = name

	// Extract pricing
	priceElement, err := wd.FindElement(selenium.ByCSSSelector, ".coordinate_item_container .price-value")
	if err != nil {
		log.Fatalf("Failed to find pricing: %v", err)
	}
	price, err := priceElement.Text()
	if err != nil {
		log.Fatalf("Failed to get pricing text: %v", err)
	}
	coordinatedProductInfo.Price = price

	// Extract product number
	productNumberElement, err := wd.FindElement(selenium.ByCSSSelector, ".coordinate_item_container .coordinate_item_tile")
	if err != nil {
		log.Fatalf("Failed to find product number element: %v", err)
	}
	productNumber, err := productNumberElement.GetAttribute("data-articleid")
	if err != nil {
		log.Fatalf("Failed to get product number attribute: %v", err)
	}
	coordinatedProductInfo.ProductNumber = productNumber

	// Extract image URL
	imageURLElement, err := wd.FindElement(selenium.ByCSSSelector, ".coordinate_item_container .coordinate_item_image")
	if err != nil {
		log.Fatalf("Failed to find image URL element: %v", err)
	}
	imageURL, err := imageURLElement.GetAttribute("src")
	if err != nil {
		log.Fatalf("Failed to get image URL attribute: %v", err)
	}
	coordinatedProductInfo.ImageURL = baseURL + imageURL

	// Extract product page URL
	productPageURLElement, err := wd.FindElement(selenium.ByCSSSelector, ".coordinate_item_container .test-link_a")
	if err != nil {
		log.Fatalf("Failed to find product page URL element: %v", err)
	}
	productPageURL, err := productPageURLElement.GetAttribute("href")
	if err != nil {
		log.Fatalf("Failed to get product page URL attribute: %v", err)
	}
	coordinatedProductInfo.ProductPageURL = baseURL + productPageURL

	return coordinatedProductInfo
}
