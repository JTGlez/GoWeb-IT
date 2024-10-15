package main

import (
	"log"
	"net/http"

	storage "github.com/JTGlez/GoWeb-IT/internal/database"
	"github.com/JTGlez/GoWeb-IT/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("Loading products from JSON file...")
	err := storage.LoadProducts("../internal/database/static/products.json")
	if err != nil {
		log.Fatalf("Failed to load products: %v", err)
	}
	log.Println("Products loaded successfully.")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", handlers.PingHandler)

	r.Route("/products", func(r chi.Router) {
		r.Get("/", handlers.ProductsHandler)
		r.Get("/{id}", handlers.ProductsByIDHandler)
		r.Get("/search", handlers.SearchProductsHandler)
		r.Post("/", handlers.AddProductHandler)
	})

	port := ":8080"
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
