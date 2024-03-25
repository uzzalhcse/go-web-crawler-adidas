package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
)

func fetchProductIds(wd selenium.WebDriver) []string {
	totalPages := []int{1}
	var ids []string

	for _, page := range totalPages {
		url := fmt.Sprintf("%v/item/?order=11&gender=mens&limit=100&category=wear&page=%v", baseURL, page)
		fmt.Printf("Fetching Product Page: %v \n", url)
		if err := wd.Get(url); err != nil {
			log.Printf("Failed to load page %s: %v", page, err)
		}

		listItems, err := wd.FindElements(selenium.ByCSSSelector, ".itemCardArea-cards")
		if err != nil {
			log.Printf("Failed to find list items: %v", err)
		}

		for index, item := range listItems {
			scrollIndex := index
			if index > 0 && index < 10 {
				scrollIndex = index * 10
				err := autoScroll(wd, fmt.Sprintf(".itemCardArea-cards:nth-child(%v)", scrollIndex))
				if err != nil {
					log.Printf("Target Dom not availabe yet: %v", fmt.Sprintf(".itemCardArea-cards:nth-child(%v)", scrollIndex))
				}
			}
			elem, err := item.FindElement(selenium.ByCSSSelector, ".image_link")
			if err != nil {
				log.Printf("[%v] Failed to find element for attribute: %v", index, err)
				continue
			}
			id, err := elem.GetAttribute("data-ga-eec-product-id")
			if err != nil {
				log.Printf("[%v] Failed to get element attribute: %v", index, err)
				continue
			}
			ids = append(ids, id)
		}
	}

	fmt.Println(fmt.Sprintf("Total %v Product Found", len(ids)))
	return ids

}

func getProductInfo(wd selenium.WebDriver) Product {
	fmt.Println("Getting Product Detail")
	var productInfo Product

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
			productInfo.Sizes = append(productInfo.Sizes, size)
		}
	}

	// Fetching sub-heading
	subheadingElem, err := wd.FindElement(selenium.ByCSSSelector, ".itemFeature")
	if err != nil {
		log.Printf("DescriptionTitle not available")
	} else {
		productInfo.DescriptionTitle, _ = subheadingElem.Text()
	}

	// Fetching main text
	mainTextElem, err := wd.FindElement(selenium.ByCSSSelector, ".commentItem-mainText")
	if err != nil {
		log.Printf("DescriptionMainText not available")
	} else {
		productInfo.DescriptionMainText, _ = mainTextElem.Text()
	}

	// Fetching article features
	articleFeatures, _ := wd.FindElements(selenium.ByCSSSelector, ".articleFeaturesItem.test-feature")
	for _, featureElem := range articleFeatures {
		feature, _ := featureElem.Text()
		productInfo.DescriptionItemization = append(productInfo.DescriptionItemization, feature)
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
		//Click on the carousel list item
		if err := item.Click(); err != nil {
			continue
		}

		wait(1)
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
	overallRatingElem, err := wd.FindElement(selenium.ByCSSSelector, ".BVRRRatingNormalOutOf .BVRRNumber.BVRRRatingNumber")
	if err != nil {
		fmt.Printf("ProductMeta not available")
		return productMeta
	}
	overallRatingText, _ := overallRatingElem.Text()
	productMeta.OverallRating = overallRatingText

	// Extract number of reviews
	numReviewsElem, _ := wd.FindElement(selenium.ByCSSSelector, ".BVRRNumber.BVRRBuyAgainTotal")
	numReviewsText, _ := numReviewsElem.Text()
	productMeta.NumberOfReviews = numReviewsText

	// Extract recommended rate
	recommendedRateElem, _ := wd.FindElement(selenium.ByCSSSelector, ".BVRRRatingPercentage .BVRRBuyAgainPercentage")
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
		idAttr, _ := reviewElem.GetAttribute("id")

		// Extract the reviewer ID
		reviewerID := parseReviewerIDFromId(idAttr)

		// Set the reviewer ID
		review.ReviewerID = reviewerID

		// Append review to the slice
		userReviews = append(userReviews, review)
	}
	productMeta.UserReviews = userReviews

	return productMeta
}
