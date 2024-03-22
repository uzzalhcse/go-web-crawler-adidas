package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"os"
	"time"
)

func main() {
	const (
		seleniumPath     = "./assets/selenium-server-standalone-3.5.3.jar"
		chromeDriverPath = "./assets/chromedriver.exe"
		port             = 8080
	)

	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath),
		selenium.Output(os.Stderr),
	}
	//selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// Maximize the window to make it full screen.
	if err := wd.MaximizeWindow(""); err != nil {
		log.Fatalf("Failed to maximize window: %v", err)
	}

	baseURL := "https://shop.adidas.jp"
	productID := "IP0418"
	url := fmt.Sprintf("%s/products/%s/", baseURL, productID)

	if err := wd.Get(url); err != nil {
		log.Fatalf("Failed to load page: %v", err)
	}

	// Click on the carousel list item to render coordinated product details
	if err := clickCarouselListItem(wd); err != nil {
		log.Fatalf("Failed to click on carousel list item: %v", err)
	} else {
		// Now you can proceed with coordinated product information extraction logic...
		coordinatedProductInfo := fetchCoordinatedProductInfo(wd, baseURL)
		fmt.Println("Coordinated Product Name:", coordinatedProductInfo.Name)
		fmt.Println("Pricing:", coordinatedProductInfo.Price)
		fmt.Println("Product Number:", coordinatedProductInfo.ProductNumber)
		fmt.Println("Image URL:", coordinatedProductInfo.ImageURL)
		fmt.Println("Product Page URL:", coordinatedProductInfo.ProductPageURL)
	}

}

// Clicks on the carousel list item to render coordinated product details
func clickCarouselListItem(wd selenium.WebDriver) error {

	elem, err := wd.FindElement(selenium.ByCSSSelector, ".coordinate_image")
	if err != nil {
		return err
	}

	// Execute JavaScript to scroll the element into view
	_, err = wd.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{elem})
	if err != nil {
		return err
	}

	// Click on the carousel list item
	if err := elem.Click(); err != nil {
		//return err
	}

	// Wait for the page to load after clicking
	time.Sleep(2 * time.Second) // Adjust the duration based on your page load time

	return nil
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
	productNumberElement, err := wd.FindElement(selenium.ByCSSSelector, ".coordinate_item_tile")
	if err != nil {
		log.Fatalf("Failed to find product number element: %v", err)
	}
	productNumber, err := productNumberElement.GetAttribute("data-articleid")
	if err != nil {
		log.Fatalf("Failed to get product number attribute: %v", err)
	}
	coordinatedProductInfo.ProductNumber = productNumber

	// Extract image URL
	imageURLElement, err := wd.FindElement(selenium.ByCSSSelector, ".coordinate_image_body.test-img")
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
