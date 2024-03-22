package main

import (
	"log"
	"time"

	"github.com/tebeka/selenium"
)

func fetchProductInfo(wd selenium.WebDriver, productID string) Product {
	var product Product
	product.ID = productID

	url := baseURL + "/products/" + productID + "/"
	if err := wd.Get(url); err != nil {
		log.Printf("Failed to load page for product ID %s: %v", productID, err)
		return product
	}

	product.Info = getProductInfo(wd)
	product.Coordinates = getCoordinatedProductInfo(wd)

	return product
}

func getProductInfo(wd selenium.WebDriver) ProductInfo {
	var productInfo ProductInfo

	productInfo.Breadcrumbs = fetchBreadcrumbs(wd)
	productInfo.Category = getText(wd, ".categoryName")
	productInfo.Name = getText(wd, ".itemTitle")
	productInfo.Price = getText(wd, ".price-value")

	sizesElements, err := wd.FindElements(selenium.ByCSSSelector, ".sizeSelectorListItemButton")
	if err != nil {
		log.Printf("Failed to find size elements: %v", err)
	} else {
		for _, sizeElement := range sizesElements {
			size, err := sizeElement.Text()
			if err != nil {
				log.Printf("Failed to get size text: %v", err)
				continue
			}
			if size != "disable" {
				productInfo.Sizes = append(productInfo.Sizes, size)
			}
		}
	}

	return productInfo
}

func getCoordinatedProductInfo(wd selenium.WebDriver) []CoordinatedProductInfo {
	var coordinatedProducts []CoordinatedProductInfo

	// Find all carousel list items
	carouselListItems, err := wd.FindElements(selenium.ByCSSSelector, ".coordinateItems .carouselListitem")
	if err != nil {
		log.Printf("Failed to find carousel list items: %v", err)
		return coordinatedProducts
	}

	for _, item := range carouselListItems {
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

		coordinatedProduct := CoordinatedProductInfo{
			Name:           getText(wd, ".coordinate_item_container .title"),
			Price:          getText(wd, ".coordinate_item_container .price-value"),
			ProductNumber:  getAttribute(wd, ".coordinate_item_tile", "data-articleid"),
			ImageURL:       baseURL + getAttribute(wd, ".coordinate_image_body.test-img", "src"),
			ProductPageURL: baseURL + getAttribute(wd, ".coordinate_item_container .test-link_a", "href"),
		}
		coordinatedProducts = append(coordinatedProducts, coordinatedProduct)
	}

	return coordinatedProducts
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

func fetchBreadcrumbs(wd selenium.WebDriver) []string {
	var breadcrumbs []string

	// Find breadcrumb items
	breadcrumbItems, err := wd.FindElements(selenium.ByCSSSelector, ".breadcrumbListItemLink")
	if err != nil {
		log.Printf("Failed to find breadcrumb items: %v", err)
		return breadcrumbs
	}

	// Extract text from breadcrumb items
	for _, item := range breadcrumbItems {
		text, err := item.Text()
		if err != nil {
			log.Printf("Failed to get breadcrumb text: %v", err)
			continue
		}
		breadcrumbs = append(breadcrumbs, text)
	}

	return breadcrumbs
}
