package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath      = "./assets/selenium-server-standalone-3.5.3.jar"
	chromeDriverPath  = "./assets/chromedriver.exe"
	firefoxDriverPath = "./assets/geckodriver.exe"
	baseURL           = "https://shop.adidas.jp"
)

var productIDs = []string{"IP0418"}

func main() {
	port := 4646
	// Start the Selenium service
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath), // Set the path to ChromeDriver
		selenium.Output(os.Stderr),              // Output logs to stderr
	}

	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		log.Fatalf("Failed to start Selenium service: %v", err)
	}
	defer service.Stop()

	// Create capabilities for Chrome browser
	caps := selenium.Capabilities{"browserName": "chrome"}
	//caps.AddChrome(chromeCaps)

	// Connect to the Selenium WebDriver
	wd, err := selenium.NewRemote(
		caps,
		fmt.Sprintf("http://localhost:%d/wd/hub", port),
	)

	if err != nil {
		log.Fatalf("Failed to create WebDriver: %v", err)
	}
	defer wd.Quit()

	for _, productID := range productIDs {
		url := fmt.Sprintf("%s/products/%s/", baseURL, productID)
		if err := wd.Get(url); err != nil {
			log.Printf("Failed to load page for product ID %s: %v", productID, err)
			continue
		}

		// Click on each carousel list item
		for i := 0; i < 3; i++ {
			carouselListItem := fmt.Sprintf(".carouselListitem:nth-child(%d)", i+1)
			elem, err := wd.FindElement(selenium.ByCSSSelector, carouselListItem)
			if err != nil {
				log.Printf("Failed to find carousel list item: %v", err)
				continue
			}

			if err := elem.Click(); err != nil {
				log.Printf("Failed to click on carousel list item: %v", err)
				continue
			}

			time.Sleep(2 * time.Second) // Wait for the content to load (adjust as needed)

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
