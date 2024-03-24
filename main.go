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

//	var productIDs = []string{
//		"IP0418",
//		"IY2911",
//		"II5763",
//		//"IT2491",
//		//"IZ4922",
//	}
//var productIds = []string{"IF9280", "ID8708", "ID5480", "IE5836", "B75807", "BD7633", "ID2350", "ID5103", "GW3774", "IE5485", "EG4959", "IF3219", "ID1994", "IF3233", "IE4230", "HQ6900", "HP8739", "IE3437", "HQ6787", "ID0985", "DB3021", "IE4931", "ID1600", "IF3235", "IE4783", "II5763", "H06260", "BD7632", "HP2201", "FX5499", "FX5500", "IE0480", "IK9149", "FX9028", "IH7502", "IG8296", "IG6421", "IE3710", "IG8482", "HQ6893"}

func main() {
	port := 8088
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		log.Fatalf("Error starting the ChromeDriver server: %v", err)
	}
	defer service.Stop()

	wd := createWebDriver(port)
	defer wd.Quit()
	ids := fetchProductIds(wd)
	for _, productID := range ids {
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
				"--headless",
				"--no-sandbox",
				"--log-level=3",
			},
		},
	)

	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(port)+"/wd/hub")
	if err != nil {
		log.Fatalf("Failed to create WebDriver: %v", err)
	}
	return wd
}
