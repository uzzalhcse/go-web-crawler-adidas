
# Adidas Product Scraper

## Overview

This project is a web scraper tool designed to extract product information from the Adidas Japan website. It utilizes Selenium WebDriver to automate web browsing and gather data such as product names, categories, prices, sizes, and more. The extracted data can be saved in Excel formats for further analysis and processing.

## Features

-   Scrapes product IDs from the men's category on the Adidas Japan website . For example `https://shop.adidas.jp/item/?order=11&gender=mens&limit=100&page=1
    `

-   Supports fetching product details and meta-information
-   Saves scraped data to Excel files

## Installation

1.  Clone the repository:

    `git clone https://github.com/uzzalhcse/go-web-crawler-adidas.git`

2.  Install dependencies:

    `go mod tidy`

3.  Download ChromeDriver and place it in the `assets` directory.

    -   Ensure that ChromeDriver version matches your installed Chrome browser version. You can download the ChromeDriver for your specific version from [here](https://storage.googleapis.com/chrome-for-testing-public/123.0.6312.58/win64/chromedriver-win64.zip)

## Usage

1.  Run the main.go file:

    `go run main.go`

2.  Wait for the scraping process to complete.

3.  Find the extracted data in the `dist/sheets` directories.


## Acknowledgments

Special thanks to the authors of the following libraries used in this project:

-   [Selenium WebDriver](https://github.com/tebeka/selenium)
-   [xlsx](https://github.com/tealeg/xlsx)