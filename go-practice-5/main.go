package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

var products = []Product{
	{1, "iPhone 17", "phones", 1000000},
	{2, "SAMSUNG Galaxy S25", "phones", 999999},
	{3, "MacBook Pro", "laptops", 799999},
	{4, "Acer Nitro 5", "laptops", 399999},
	{5, "LG TV", "tvs", 99999},
	{6, "SAMSUNG TV", "tvs", 199999},
}

func main() {
	http.HandleFunc("/products", getProductsHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	category := r.URL.Query().Get("category")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	sortParam := r.URL.Query().Get("sort")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	filtered := make([]Product, 0, len(products))
	for _, p := range products {
		if category != "" && !strings.EqualFold(p.Category, category) {
			continue
		}
		if minPriceStr != "" {
			if min, err := strconv.Atoi(minPriceStr); err == nil && p.Price < min {
				continue
			}
		}
		if maxPriceStr != "" {
			if max, err := strconv.Atoi(maxPriceStr); err == nil && p.Price > max {
				continue
			}
		}
		filtered = append(filtered, p)
	}

	if sortParam == "price_asc" {
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].Price < filtered[j].Price })
	} else if sortParam == "price_desc" {
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].Price > filtered[j].Price })
	}

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)
	if offset > len(filtered) {
		filtered = []Product{}
	} else {
		filtered = filtered[offset:]
	}
	if limit > 0 && limit < len(filtered) {
		filtered = filtered[:limit]
	}

	duration := time.Since(start)
	w.Header().Set("X-Query-Time", duration.String())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}
