package main

type Product struct {
	ID                     string
	Name                   string
	Category               string
	Price                  string
	Sizes                  []string
	Breadcrumbs            []string
	DescriptionTitle       string
	DescriptionMainText    string
	DescriptionItemization []string
	Coordinates            []CoordinatedProductInfo
	SizeChart              SizeChart
	ProductMeta            ProductMeta
}

type CoordinatedProductInfo struct {
	Name           string
	Price          string
	ProductNumber  string
	ImageURL       string
	ProductPageURL string
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
