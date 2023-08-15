package api

import (
	"github.com/go-chi/chi"
)

func NewProductRouter(productHandler ProductHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/products/create", productHandler.CreateProductHandler)
	router.Get("/api/products/findbyId", productHandler.FindProductHandler)
	router.Get("/api/products/findall", productHandler.FindProductAllHandler)
	router.Put("/api/products/deletebyId", productHandler.DeleteProductAllHandler)

	return router
}
