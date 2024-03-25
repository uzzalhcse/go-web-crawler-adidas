package main

import (
	"github.com/tebeka/selenium/chrome"
	"log"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
)

const (
	baseURL          = "https://shop.adidas.jp"
	chromeDriverPath = "./driver/chromedriver.exe"
	port             = 4444
)

var products []Product

func main() {
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		log.Fatalf("Error starting the ChromeDriver server: %v", err)
	}
	defer service.Stop()

	wd := createWebDriver(port)
	defer wd.Quit()

	ids := fetchProductIds(wd)
	startTime := time.Now()
	fetchProduct(wd, ids)

	if err := saveProductInfoJSON(products); err != nil {
		log.Printf("Failed to save Json %v", err)
	}
	if err := saveProductInfoSpreadsheet(products); err != nil {
		log.Printf("Failed to save Xlxs %v", err)
	}
	totalTime := time.Since(startTime).Minutes()
	log.Printf("Total time elapsed: %.2f Minutes", totalTime)

}
func fetchProduct(wd selenium.WebDriver, ids []string) {
	for _, productID := range ids {
		product := fetchProductInfo(wd, productID)
		products = append(products, product)

		//fetchDummy(wd, productID)
		log.Printf("Product info for ID %s Fetched successfully", productID)
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
				"--log-level=3",
				//"--disable-gpu",
				//"--disable-dev-shm-usage",
				//"--disable-web-security",
			},
		},
	)
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(port)+"/wd/hub")
	if err != nil {
		log.Fatalf("Failed to create WebDriver: %v", err)
	}
	return wd
}

func fetchProductInfo(wd selenium.WebDriver, productID string) Product {
	startTime := time.Now()
	var product Product
	url := baseURL + "/products/" + productID + "/"

	log.Printf("Fetching Product info for Url %s", url)
	if err := wd.Get(url); err != nil {
		log.Printf("Failed to load page for product ID %s: %v", productID, err)
	}
	//time.Sleep(2 * time.Second)
	err := autoScroll(wd, ".coordinateItems .carouselListitem")
	if err != nil {
		_ = autoScroll(wd, ".js-articlePromotion")
	}

	product = getProductInfo(wd)
	product.ID = productID
	product.Coordinates = getCoordinatedProductInfo(wd)
	product.SizeChart = parseSizeChartHTML(wd)
	product.ProductMeta = parseProductMeta(wd)

	totalTime := time.Since(startTime).Seconds()
	log.Printf("Time elapsed: %.2f Secounds", totalTime)
	return product
}
