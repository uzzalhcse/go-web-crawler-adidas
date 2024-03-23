package main

import (
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
	"IY2911",
	"II5763",
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
		if product.Name == "" {
			log.Printf("Failed to fetch product info for ID %s", productID)
			continue
		}

		if err := saveProductInfoJSON(product, productID); err != nil {
			log.Printf("Failed to save product info for ID %s: %v", productID, err)
			continue
		}
		if err := saveProductInfoSpreadsheet(product); err != nil {
			log.Printf("Failed to save product info for ID %s: %v", productID, err)
			continue
		}
		log.Printf("Product info for ID %s saved successfully", productID)
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
				//"--headless",
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
