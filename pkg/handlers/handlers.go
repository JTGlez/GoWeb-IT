package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JTGlez/GoWeb-IT/pkg/models"
	"github.com/JTGlez/GoWeb-IT/pkg/storage"
	"github.com/go-chi/chi/v5"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {

	products := []models.Product{}
	for _, product := range storage.Store {
		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}

func ProductsByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
	}

	product, exists := storage.Store[id]

	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(product)
}

func SearchProductsHandler(w http.ResponseWriter, r *http.Request) {
	priceGtStr := r.URL.Query().Get("priceGt")

	priceGt, err := strconv.ParseFloat(priceGtStr, 64)
	if err != nil {
		http.Error(w, "Invalid priceGt value", http.StatusBadRequest)
		return
	}

	products := []models.Product{}
	for _, product := range storage.Store {
		if product.Price > priceGt {
			products = append(products, product)
		}
	}

	json.NewEncoder(w).Encode(products)
}
