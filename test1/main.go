package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"strconv"
	"time"
)

const (
	seleniumPath      = "./assets/selenium-server-standalone-3.5.3.jar"
	chromeDriverPath  = "./assets/chromedriver.exe"
	firefoxDriverPath = "./assets/geckodriver.exe"
	baseURL           = "https://shop.adidas.jp"
)

var productIDs = []string{
	"IP0418",
	"IY2911",
	//"II5763",
	//"IT2491",
	//"IZ4922",
}

func main() {
	//browserPath := GetBrowserPath("chromium")
	port := 8088

	var opts []selenium.ServiceOption
	service, err := selenium.NewChromeDriverService(chromeDriverPath,
		port, opts...)

	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// NB: Path is the important part here
	caps.AddChrome(
		chrome.Capabilities{
			Path: "",
			Args: []string{
				"--window-size=1920,1080",
				//"--headless",
				"--no-sandbox",
				"--blink-settings=imagesEnabled=false", // Disable images
				"--blink-settings=cssEnabled=false",    // Disable CSS
			},
		},
	)

	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(port)+"/wd/hub")
	if err != nil {
		panic(err)
	}

	wd.Refresh()

	defer wd.Quit()
	defer service.Stop()

	for _, productID := range productIDs {
		url := fmt.Sprintf("%s/products/%s/", baseURL, productID)
		if err := wd.Get(url); err != nil {
			log.Printf("Failed to load page for product ID %s: %v", productID, err)
			continue
		}

		// Find all carousel list items
		carouselListItems, err := wd.FindElements(selenium.ByCSSSelector, ".coordinateItems .carouselListitem")
		if err != nil {
			log.Printf("Failed to find carousel list items: %v", err)
		}
		if len(carouselListItems) < 1 {
			log.Printf("No carouselListItems to click: %v", len(carouselListItems))
			continue
		}

		// Loop through each carousel list item
		for _, item := range carouselListItems {
			// Limit the loop to a certain number of items if needed
			// if i >= 2 {
			//     break
			// }

			// Scroll the item into view
			if _, err := wd.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{item}); err != nil {
				log.Printf("Failed to scroll element into view: %v", err)
				continue
			}

			// Click on the carousel list item
			if err := item.Click(); err != nil {
				log.Printf("Failed to click on carousel list item: %v", err)
				continue
			}

			// Wait for the content to load (adjust as needed)
			time.Sleep(2 * time.Second)

			// Fetch coordinated product information
			coordinatedProductInfo := fetchCoordinatedProductInfo(wd)
			fmt.Println("Product ID:", productID)
			fmt.Println("Coordinated Product Name:", coordinatedProductInfo.Name)
			fmt.Println("Pricing:", coordinatedProductInfo.Price)
			fmt.Println("Product Number:", coordinatedProductInfo.ProductNumber)
			fmt.Println("Image URL:", coordinatedProductInfo.ImageURL)
			fmt.Println("Product Page URL:", coordinatedProductInfo.ProductPageURL)
			fmt.Println()
		}

	}
}

type CoordinatedProductInfo struct {
	Name           string
	Price          string
	ProductNumber  string
	ImageURL       string
	ProductPageURL string
}

func fetchCoordinatedProductInfo(wd selenium.WebDriver) CoordinatedProductInfo {
	var coordinatedProductInfo CoordinatedProductInfo

	coordinatedProductInfo.Name = getText(wd, ".coordinate_item_container .title")
	coordinatedProductInfo.Price = getText(wd, ".coordinate_item_container .price-value")
	coordinatedProductInfo.ProductNumber = getAttribute(wd, ".coordinate_item_tile", "data-articleid")
	coordinatedProductInfo.ImageURL = baseURL + getAttribute(wd, ".coordinate_image_body.test-img", "src")
	coordinatedProductInfo.ProductPageURL = baseURL + getAttribute(wd, ".coordinate_item_container .test-link_a", "href")

	return coordinatedProductInfo
}

func getText(wd selenium.WebDriver, selector string) string {
	elem, err := wd.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		log.Fatalf("Failed to find element for text: %v", err)
	}
	text, err := elem.Text()
	if err != nil {
		log.Fatalf("Failed to get element text: %v", err)
	}
	return text
}

func getAttribute(wd selenium.WebDriver, selector, attribute string) string {
	elem, err := wd.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		log.Fatalf("Failed to find element for attribute: %v", err)
	}
	attr, err := elem.GetAttribute(attribute)
	if err != nil {
		log.Fatalf("Failed to get element attribute: %v", err)
	}
	return attr
}
