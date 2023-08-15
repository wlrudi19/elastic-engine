package api

import (
	"github.com/go-chi/chi"
)

func NewProductRouter(productHandler ProductHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/products/create", productHandler.CreateProductHandler)

	return router
}
