package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

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

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	if newProduct.Name == "" || newProduct.Quantity == 0 || newProduct.CodeValue == "" || newProduct.Expiration == "" || newProduct.Price == 0 {
		http.Error(w, "All fields except 'is_published' are required", http.StatusBadRequest)
		return
	}

	datePattern := `^\d{2}/\d{2}/\d{4}$`
	match, _ := regexp.MatchString(datePattern, newProduct.Expiration)
	if !match {
		http.Error(w, "Invalid expiration date format. Use XX/XX/XXXX", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("02/01/2006", newProduct.Expiration)
	if err != nil {
		http.Error(w, "Invalid expiration date", http.StatusBadRequest)
		return
	}

	if _, exists := storage.CodeIndex[newProduct.CodeValue]; exists {
		http.Error(w, "Product with that code already exists", http.StatusBadRequest)
		return
	}

	storage.LastID++
	newProduct.ID = storage.LastID

	storage.Store[newProduct.ID] = newProduct
	storage.CodeIndex[newProduct.CodeValue] = newProduct.ID

	log.Printf("New product created in memory: %+v\n", newProduct)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}
