package main

import (
	"fmt"
	"github.com/tebeka/selenium/chrome"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/tebeka/selenium"
)

const (
	baseURL          = "https://shop.adidas.jp"
	baseURLDummy     = "http://localhost:8080"
	chromeDriverPath = "./assets/chromedriver.exe"
	maxWorkers       = 20
	port             = 4444
)

var ids = []string{
	"IF9280",
	"ID8708",
	"ID5480",
	"B75807", "BD7633", "ID2350", "IE5836", "ID5103", "GW3774", "IF3219",
	"IE5485", "EG4959", "ID1994", "HQ6900", "IE4230", "IF3233", "IE3437", "HP8739", "ID0985", "HQ6787",
	"DB3021", "ID1600", "BD7632", "IF3235", "H06260", "HP2201", "IE0480", "IE4931", "II5763", "IE4783",
	"IH7502", "FX5500", "FX5499", "IK9149", "FX9028", "IG6421", "IG8296", "IG8482", "IE3710", "HQ6893",
	"IE4195", "IY8077", "ID5961", "IF1953", "IG1504", "IG6047", "IG6049", "IF0202", "IG6201", "IG8301",
	"IS0541", "IE0422", "IF6491", "CQ2809", "IF3813", "IG6440", "GY5695", "FX5509", "HQ6785", "GX3937",
	"ID1995", "EG4958", "IH7496", "ID6990", "IZ4926", "FX5508", "FX5501", "ID0986", "ID4637", "IG6191",
	"IF3914", "IU0106", "IZ4923", "HZ5730", "HQ8718", "ID0983", "IG8295", "IG5929", "IJ7055", "FV0321",
	"IF8047", "IF9773", "GY0042", "IG4036", "ID4950", "IT2491", "IF6162", "IT2492", "IG8661", "IY8075",
	"IE3709", "EG4957", "IF6514", "IY8076", "IF8770", "IG6192", "IG3136", "IG1561", "ID2564", "GZ1537",
}

var products []Product

func main() {
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		log.Fatalf("Error starting the ChromeDriver server: %v", err)
	}
	defer service.Stop()

	wd := createWebDriver(port)
	defer wd.Quit()

	//ids := fetchProductIds(wd)
	startTime := time.Now()
	//fetchProductAsync(ids)
	fetchProduct(wd, ids)

	if err := saveProductInfoSpreadsheet(products); err != nil {
		log.Printf("Failed to save Xlxs %v", err)
	}
	if err := saveProductInfoJSON(products); err != nil {
		log.Printf("Failed to save Json %v", err)
	}
	totalTime := time.Since(startTime).Seconds()
	log.Printf("Total time elapsed: %.2f Seconds", totalTime)

}
func fetchProductAsync(ids []string) {
	var wg sync.WaitGroup
	jobs := make(chan string)

	fmt.Println("Len", len(ids))

	// Start worker pool
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for productID := range jobs {
				wd := createWebDriver(port) // Create a new WebDriver instance
				// Fetch product info in the new tab
				product := fetchProductInfo(wd, productID)
				products = append(products, product)

				//fetchDummy(wd, productID)

				log.Printf("Product info for ID %s Fetched successfully", productID)
				err := wd.Quit()
				if err != nil {
					fmt.Println("Error on Quiting Browser")
				} // Quit the WebDriver instance after use
			}
		}()
	}

	// Add jobs to the channel
	for i, productID := range ids {
		fmt.Println("Sending Job:", i)
		jobs <- productID
		time.Sleep(500 * time.Millisecond)
	}
	close(jobs)

	wg.Wait()
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
				"--headless",
				"--no-sandbox",
				"--log-level=3",
				"--disable-gpu",
				"--disable-dev-shm-usage",
				"--disable-web-security",
			},
		},
	)
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(port)+"/wd/hub")
	if err != nil {
		log.Fatalf("Failed to create WebDriver: %v", err)
	}
	return wd
}

func fetchDummy(wd selenium.WebDriver, productID string) {
	startTime := time.Now()
	url := baseURLDummy + "?product=" + productID

	log.Printf("Fetching Product info for Url %s", url)
	if err := wd.Get(url); err != nil {
		log.Printf("Failed to load page for product ID %s: %v", productID, err)
	}
	title := getText(wd, ".title")
	fmt.Println(title)
	dynamicTitle := getText(wd, ".dynamic-title")
	fmt.Println(dynamicTitle)
	totalTime := time.Since(startTime).Seconds()
	log.Printf("Time elapsed: %.2f Secounds", totalTime)
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
	err := autoScroll(wd, ".js-articlePromotion")
	if err != nil {
		log.Printf("Failed to scroll the page: %v", err)
	}

	product = getProductInfo(wd)
	product.ID = productID
	//product.Coordinates = getCoordinatedProductInfo(wd)
	product.SizeChart = parseSizeChartHTML(wd)
	product.ProductMeta = parseProductMeta(wd)

	totalTime := time.Since(startTime).Seconds()
	log.Printf("Time elapsed: %.2f Secounds", totalTime)
	return product
}
