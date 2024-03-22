package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath     = "./assets/selenium-server-standalone-3.5.3.jar"
	chromeDriverPath = "./assets/chromedriver.exe"
	port             = 8080
	baseURL          = "https://shop.adidas.jp"
)

var productIDs = []string{"IP0418", "IP0418", "IP0418", "IP0418"}

func main() {
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath),
		selenium.Output(os.Stderr),
	}

	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		log.Fatalf("Failed to start Selenium service: %v", err)
	}
	defer service.Stop()

	wd, err := selenium.NewRemote(
		selenium.Capabilities{"browserName": "chrome"},
		fmt.Sprintf("http://localhost:%d/wd/hub", port),
	)
	if err != nil {
		log.Fatalf("Failed to create WebDriver: %v", err)
	}
	defer wd.Quit()

	if err := wd.MaximizeWindow(""); err != nil {
		log.Fatalf("Failed to maximize window: %v", err)
	}

	for _, productID := range productIDs {
		url := fmt.Sprintf("%s/products/%s/", baseURL, productID)
		if err := wd.Get(url); err != nil {
			log.Printf("Failed to load page for product ID %s: %v", productID, err)
			continue
		}

		carouselListItems, err := wd.FindElements(selenium.ByCSSSelector, ".carouselListitem")
		if err != nil {
			log.Printf("Failed to find carousel list items for product ID %s: %v", productID, err)
			continue
		}

		for _, item := range carouselListItems {
			if err := clickCarouselListItem(wd, item); err != nil {
				log.Printf("Failed to click on carousel list item for product ID %s: %v", productID, err)
				continue
			}
			coordinatedProductInfo := fetchCoordinatedProductInfo(wd)
			fmt.Println("Product ID:", productID)
			fmt.Println("Coordinated Product Name:", coordinatedProductInfo.Name)
			fmt.Println("Pricing:", coordinatedProductInfo.Price)
			fmt.Println("Product Number:", coordinatedProductInfo.ProductNumber)
			fmt.Println("Image URL:", coordinatedProductInfo.ImageURL)
			fmt.Println("Product Page URL:", coordinatedProductInfo.ProductPageURL)

			time.Sleep(3 * time.Second)
		}
	}
}

func clickCarouselListItem(wd selenium.WebDriver, elem selenium.WebElement) error {
	if _, err := wd.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{elem}); err != nil {
		return fmt.Errorf("failed to scroll element into view: %v", err)
	}
	if err := elem.Click(); err != nil {
		return fmt.Errorf("failed to click element: %v", err)
	}
	return nil
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

	if err := waitForLoad(wd, time.Second*3, ".coordinate_item_container .title"); err != nil {
		log.Fatalf("Failed to find coordinated product name: %v", err)
	}

	coordinatedProductInfo.Name = getText(wd, ".coordinate_item_container .title")
	coordinatedProductInfo.Price = getText(wd, ".coordinate_item_container .price-value")
	coordinatedProductInfo.ProductNumber = getAttribute(wd, ".coordinate_item_tile", "data-articleid")
	coordinatedProductInfo.ImageURL = baseURL + getAttribute(wd, ".coordinate_image_body.test-img", "src")
	coordinatedProductInfo.ProductPageURL = baseURL + getAttribute(wd, ".coordinate_item_container .test-link_a", "href")

	return coordinatedProductInfo
}

func waitForLoad(wd selenium.WebDriver, timeout time.Duration, selector string) error {
	return wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		elem, err := wd.FindElement(selenium.ByCSSSelector, selector)
		if err != nil || elem == nil {
			return false, nil
		}
		return true, nil
	}, timeout)
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
