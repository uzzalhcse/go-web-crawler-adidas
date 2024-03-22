package main

type Product struct {
	ID          string
	Info        ProductInfo
	Coordinates []CoordinatedProductInfo
}

type ProductInfo struct {
	Category    string
	Name        string
	Price       string
	Sizes       []string
	Breadcrumbs []string
}

type CoordinatedProductInfo struct {
	Name           string
	Price          string
	ProductNumber  string
	ImageURL       string
	ProductPageURL string
}
