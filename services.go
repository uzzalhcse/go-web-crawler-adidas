package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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

	// Scroll to the bottom of the page
	if err := scrollPageToBottom(wd); err != nil {
		log.Printf("Failed to scroll to bottom of page: %v", err)
	}

	product.Info = getProductInfo(wd)
	product.Coordinates = getCoordinatedProductInfo(wd)
	product.Description = getProductDescription(wd)
	product.SizeChart = parseSizeChartHTML(wd)
	product.ProductMeta = parseProductMeta(wd)

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
func getProductDescription(wd selenium.WebDriver) ProductDescription {
	var description ProductDescription

	// Fetching subheading
	subheadingElem, err := wd.FindElement(selenium.ByCSSSelector, ".itemFeature")
	if err != nil {
		log.Printf("Failed to find subheading element: %v", err)
	} else {
		description.Title, _ = subheadingElem.Text()
	}

	// Fetching main text
	mainTextElem, err := wd.FindElement(selenium.ByCSSSelector, ".commentItem-mainText")
	if err != nil {
		log.Printf("Failed to find main text element: %v", err)
	} else {
		description.MainText, _ = mainTextElem.Text()
	}

	// Fetching article features
	articleFeatures, _ := wd.FindElements(selenium.ByCSSSelector, ".articleFeaturesItem.test-feature")
	for _, featureElem := range articleFeatures {
		feature, _ := featureElem.Text()
		description.ArticleFeatures = append(description.ArticleFeatures, feature)
	}

	return description
}

func parseSizeChartHTML(wd selenium.WebDriver) SizeChart {
	var sizeChart SizeChart

	// Find all table rows in the size chart
	rows, err := wd.FindElements(selenium.ByCSSSelector, ".sizeChartTRow")
	if err != nil {
		log.Printf("Failed to find size chart rows: %v", err)
		return sizeChart
	}

	// Extract column headers
	columnHeaders, err := wd.FindElements(selenium.ByCSSSelector, ".sizeChartTHeaderCell")
	if err != nil {
		log.Printf("Failed to find size chart column headers: %v", err)
		return sizeChart
	}

	var headerRow []string
	for _, header := range columnHeaders {
		text, err := header.Text()
		if err != nil {
			log.Printf("Failed to get column header text: %v", err)
			continue
		}
		headerRow = append(headerRow, text)
	}
	sizeChart.Measurements = append(sizeChart.Measurements, headerRow)

	// Iterate over each row
	for _, row := range rows {
		// Find all table cells in the row
		cells, err := row.FindElements(selenium.ByCSSSelector, ".sizeChartTCell")
		if err != nil {
			log.Printf("Failed to find size chart cells: %v", err)
			continue
		}

		var measurements []string

		// Iterate over each cell
		for _, cell := range cells {
			// Get the text content of the cell
			text, err := cell.Text()
			if err != nil {
				log.Printf("Failed to get size chart cell text: %v", err)
				continue
			}

			// Append the text content to the measurements slice
			measurements = append(measurements, text)
		}

		// Check if the row contains measurements
		if len(measurements) > 0 {
			sizeChart.Measurements = append(sizeChart.Measurements, measurements)
		}
	}

	return sizeChart
}

func parseProductMeta(wd selenium.WebDriver) ProductMeta {
	var productMeta ProductMeta

	// Extract overall rating
	overallRatingElem, _ := wd.FindElement(selenium.ByCSSSelector, ".BVRRRatingNormalOutOf .BVRRNumber.BVRRRatingNumber")
	overallRatingText, _ := overallRatingElem.Text()
	productMeta.OverallRating = overallRatingText

	// Extract number of reviews
	numReviewsElem, _ := wd.FindElement(selenium.ByCSSSelector, ".BVRRNumber.BVRRBuyAgainTotal")
	numReviewsText, _ := numReviewsElem.Text()
	productMeta.NumberOfReviews = numReviewsText

	// Extract recommended rate
	recommendedRateElem, _ := wd.FindElement(selenium.ByCSSSelector, ".BVRRNumber.BVRRBuyAgainRecommend")
	recommendedRateText, _ := recommendedRateElem.Text()
	productMeta.RecommendedRate = recommendedRateText

	// Extract item ratings
	var itemRatings []ItemRating
	itemRatingElems, _ := wd.FindElements(selenium.ByCSSSelector, ".BVRRSecondaryRatingsContainer .BVRRRatingEntry")
	for _, itemRatingElem := range itemRatingElems {

		labelElem, _ := itemRatingElem.FindElement(selenium.ByCSSSelector, ".BVRRLabel")
		label, _ := labelElem.Text()
		ratingElem, _ := itemRatingElem.FindElement(selenium.ByCSSSelector, ".BVRRRatingRadioImage img")
		rating, _ := ratingElem.GetAttribute("title")

		itemRatings = append(itemRatings, ItemRating{Label: label, Rating: rating})
	}
	productMeta.ItemRatings = itemRatings

	// Extract user reviews
	var userReviews []Review
	reviewElems, _ := wd.FindElements(selenium.ByCSSSelector, ".BVRRContentReview")
	for _, reviewElem := range reviewElems {
		review := Review{}

		// Extract review date
		dateElem, _ := reviewElem.FindElement(selenium.ByCSSSelector, ".BVRRReviewDate")
		date, _ := dateElem.Text()
		review.Date = date

		// Extract review title
		titleElem, _ := reviewElem.FindElement(selenium.ByCSSSelector, ".BVRRValue.BVRRReviewTitle")
		title, _ := titleElem.Text()
		review.Title = title

		// Extract review description
		descElem, _ := reviewElem.FindElement(selenium.ByCSSSelector, ".BVRRReviewTextContainer")
		desc, _ := descElem.Text()
		review.Description = desc

		// Extract review rating
		ratingElem, _ := reviewElem.FindElement(selenium.ByCSSSelector, ".BVRRNumber.BVRRRatingNumber")
		ratingText, _ := ratingElem.Text()
		review.Rating = ratingText

		// Extract reviewer ID from the name attribute
		anchorElem, _ := reviewElem.FindElement(selenium.ByCSSSelector, "a.BVRRUserProfileImageLink")
		hrefAttr, _ := anchorElem.GetAttribute("href")

		// Extract the reviewer ID from the href attribute
		reviewerID := parseReviewerIDFromHrefAttr(hrefAttr)

		// Set the reviewer ID
		review.ReviewerID = reviewerID

		// Append review to the slice
		userReviews = append(userReviews, review)
	}
	productMeta.UserReviews = userReviews

	return productMeta
}
func parseReviewerIDFromHrefAttr(hrefAttr string) string {
	// Split the href attribute value by "/"
	parts := strings.Split(hrefAttr, "/")
	// The reviewer ID should be the second-to-last part
	reviewerID := parts[len(parts)-2]
	return reviewerID
}
func scrollPageToBottom(wd selenium.WebDriver) error {
	// Execute JavaScript to scroll to the bottom of the page
	_, err := wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil)
	if err != nil {
		return fmt.Errorf("failed to scroll to bottom of page: %w", err)
	}
	// Add a delay to allow content to load after scrolling
	time.Sleep(2 * time.Second)
	return nil
}

func saveProductInfoJSON(product Product, productID string) error {
	// Marshal product to JSON
	productJSON, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal product info for ID %s: %w", productID, err)
	}

	// Create dist folder if it doesn't exist
	if err := os.MkdirAll("dist", 0755); err != nil {
		return fmt.Errorf("failed to create dist folder: %w", err)
	}

	// Write JSON to file in dist folder
	filename := fmt.Sprintf("dist/product_%s.json", productID)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	if _, err := file.Write(productJSON); err != nil {
		return fmt.Errorf("failed to write JSON to file %s: %w", filename, err)
	}

	return nil
}
