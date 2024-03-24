package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		product := r.URL.Query().Get("product")
		if product == "" {
			product = "Dynamic Title"
		}

		time.Sleep(5 * time.Second)

		html := fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>%s</title>
			</head>
			<body>
				<h1 class="title">Nice Websites</h1>
				<h1 class="dynamic-title">The Product ID is:  %s</h1>
			</body>
			</html>
		`, product, product)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, html)
		// Log the request to the terminal
		log.Printf("Request received: %s %s", r.Method, r.URL)
	})

	fmt.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
