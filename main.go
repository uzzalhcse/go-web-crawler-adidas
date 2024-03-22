package main

import (
	"fmt"
	"strings"
)

func main() {
	url := "https://shop.adidas.jp/products/IE3439/"
	breadcrumbString := fetchBreadcrumbString(url)
	fmt.Println("breadcrumbString:", breadcrumbString)

	productInfo := fetchProductInfo(url)
	fmt.Println("Category:", productInfo.Category)
	fmt.Println("Product Name:", productInfo.Name)
	fmt.Println("Pricing:", productInfo.Price)
	fmt.Println("Available Sizes:", strings.Join(productInfo.Sizes, ", "))
}
