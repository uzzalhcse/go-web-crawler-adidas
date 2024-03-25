package main

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/tebeka/selenium"
	"log"
	"os"
	"strings"
	"time"
)

func wait(second time.Duration) {
	time.Sleep(second * time.Second)
}

func autoScroll(wd selenium.WebDriver, selector string) error {
	el, err := wd.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		return fmt.Errorf("dom scrolling error: %v", err)
	}
	//fmt.Println("autoScroll...")
	if _, err := wd.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{el}); err != nil {
		log.Printf("Failed to scroll element into view: %v", err)
	}
	wait(1)
	return nil
}
func getAttribute(wd selenium.WebDriver, selector, attribute string) string {
	elem, err := wd.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		log.Printf("Failed to find element for attribute: %v", err)
		return ""
	}
	attr, err := elem.GetAttribute(attribute)
	if err != nil {
		log.Printf("Failed to get element attribute: %v", err)
		return ""
	}
	return attr
}
func getText(wd selenium.WebDriver, selector string) string {
	elem, err := wd.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		log.Printf("Failed to find element for text: %v", err)
		return ""
	}
	text, err := elem.Text()
	if err != nil {
		log.Printf("Failed to get element text: %v", err)
	}
	return text
}

func parseReviewerIDFromId(id string) string {
	parts := strings.Split(id, "_")
	// The reviewer ID should be the last part of the ID string
	reviewerID := parts[len(parts)-1]
	return reviewerID
}

func saveProductInfoJSON(products []Product) error {
	// Create dist folder if it doesn't exist
	if err := os.MkdirAll("dist/json", 0755); err != nil {
		return fmt.Errorf("failed to create dist folder: %w", err)
	}
	// Open or create the JSON file
	filename := fmt.Sprintf("dist/json/products.json")
	var file *os.File
	var err error
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filename, err)
		}
		defer file.Close()

		// Write the opening bracket of the JSON array
		if _, err := file.WriteString("[\n"); err != nil {
			return fmt.Errorf("failed to write to file %s: %w", filename, err)
		}
	} else {
		file, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", filename, err)
		}
		defer file.Close()

		// Check if file is empty
		stat, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file stat %s: %w", filename, err)
		}
		if stat.Size() > 2 { // Check if file has data other than "[" and "]"
			// Move the file pointer back to remove the last closing bracket and newline
			_, err = file.Seek(-2, os.SEEK_END)
			if err != nil {
				return fmt.Errorf("failed to seek file %s: %w", filename, err)
			}
		} else {
			// Write a newline before adding new data
			if _, err := file.WriteString("\n"); err != nil {
				return fmt.Errorf("failed to write to file %s: %w", filename, err)
			}
		}
	}

	for i, product := range products {
		// Marshal product to JSON
		productJSON, err := json.Marshal(product)
		if err != nil {
			return fmt.Errorf("failed to marshal product info for index %d: %w", i, err)
		}

		// Write the JSON object to the file
		if _, err := file.Write(productJSON); err != nil {
			return fmt.Errorf("failed to write to file %s: %w", filename, err)
		}
		if i < len(products)-1 {
			// Add comma and newline after each JSON object except the last one
			if _, err := file.WriteString(",\n"); err != nil {
				return fmt.Errorf("failed to write to file %s: %w", filename, err)
			}
		}
	}

	// Write the closing bracket of the JSON array
	if _, err := file.WriteString("\n]"); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filename, err)
	}

	return nil
}

func saveProductInfoSpreadsheet(products []Product) error {
	// Create dist folder if it doesn't exist
	if err := os.MkdirAll("dist/sheets", 0755); err != nil {
		return fmt.Errorf("failed to create dist folder: %w", err)
	}
	// Open the existing Excel file or create a new one if it doesn't exist
	filename := fmt.Sprintf("dist/sheets/product.xlsx")
	fileExcel, err := xlsx.OpenFile(filename)
	if err != nil {
		fileExcel = xlsx.NewFile()
	}

	// Get the sheet "Products" or create a new one if it doesn't exist
	var sheet *xlsx.Sheet
	if len(fileExcel.Sheets) == 0 {
		sheet, err = fileExcel.AddSheet("Products")
		if err != nil {
			fmt.Println("Error adding sheet:", err)
			return err
		}
	} else {
		sheet = fileExcel.Sheets[0] // Get the first sheet
	}

	// If the sheet is empty, add headers
	if sheet.MaxRow == 0 {
		row := sheet.AddRow()
		headers := []string{"ID", "Category", "Name", "Price", "Sizes", "Breadcrumbs", "Coordinates", "Description Title", "Description MainText", "Description Itemization", "SizeChart", "OverallRating", "NumberOfReviews", "RecommendedRate", "KWs", "ItemRatings", "UserReviews"}
		for _, header := range headers {
			cell := row.AddCell()
			cell.SetString(header)
		}
	}

	for _, product := range products {

		row := sheet.AddRow()
		// Convert Coordinates to JSON
		coordinatesJSON, err := json.Marshal(product.Coordinates)
		if err != nil {
			fmt.Println("Error marshalling Coordinates to JSON:", err)
		}

		// Convert ItemRatings to JSON
		itemRatingsJSON, err := json.Marshal(product.ProductMeta.ItemRatings)
		if err != nil {
			fmt.Println("Error marshalling ItemRatings to JSON:", err)
		}

		// Convert UserReviews to JSON
		userReviewsJSON, err := json.Marshal(product.ProductMeta.UserReviews)
		if err != nil {
			fmt.Println("Error marshalling UserReviews to JSON:", err)
		}

		row.AddCell().SetString(product.ID)
		row.AddCell().SetString(product.Category)
		row.AddCell().SetString(product.Name)
		row.AddCell().SetString(product.Price)
		row.AddCell().SetString(strings.Join(product.Sizes, ","))
		row.AddCell().SetString(strings.Join(product.Breadcrumbs, ","))
		row.AddCell().SetString(string(coordinatesJSON)) // Add Coordinates JSON
		row.AddCell().SetString(product.DescriptionTitle)
		row.AddCell().SetString(product.DescriptionMainText)
		row.AddCell().SetString(strings.Join(product.DescriptionItemization, ","))
		row.AddCell().SetString(fmt.Sprintf("%v", product.SizeChart))
		row.AddCell().SetString(product.ProductMeta.OverallRating)
		row.AddCell().SetString(product.ProductMeta.NumberOfReviews)
		row.AddCell().SetString(product.ProductMeta.RecommendedRate)
		row.AddCell().SetString(strings.Join(product.Tags, ","))
		row.AddCell().SetString(string(itemRatingsJSON)) // Add ItemRatings JSON
		row.AddCell().SetString(string(userReviewsJSON)) // Add UserReviews JSON

	}
	// Save the Excel file
	err = fileExcel.Save(filename)
	if err != nil {
		fmt.Println("Error saving Excel file:", err)
		return err
	}
	fmt.Println("Product saved into excel file successfully")

	return nil
}
