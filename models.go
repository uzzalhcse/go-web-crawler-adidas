package main

type Product struct {
	ID          string
	Info        ProductInfo
	Coordinates []CoordinatedProductInfo
	Description ProductDescription
	SizeChart   SizeChart
	ProductMeta ProductMeta
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
type ProductDescription struct {
	Title           string
	MainText        string
	ArticleFeatures []string
}
type SizeChart struct {
	CategoryNames []string
	Measurements  [][]string
}

// Review struct represents a single review
type Review struct {
	Date        string
	Rating      string
	Title       string
	Description string
	ReviewerID  string
}

// ItemRating struct represents the rating for a specific item
type ItemRating struct {
	Label  string
	Rating string
}

// ProductInfo struct represents information about the product
type ProductMeta struct {
	OverallRating   string
	NumberOfReviews string
	RecommendedRate string // Percentage of reviewers who recommend the product
	ItemRatings     []ItemRating
	UserReviews     []Review
}
