package api

import (
	"github.com/go-chi/chi"
	"github.com/wlrudi19/elastic-engine/infrastructure/middlewares"
)

func NewProductRouter(productHandler ProductHandler) *chi.Mux {
	router := chi.NewRouter()

	authMiddleware := middlewares.Authenticate

	router.Route("/api/products", func(r chi.Router) {
		r.With(authMiddleware).Get("/create", productHandler.CreateProductHandler)
		r.Get("/findbyId", productHandler.FindProductHandler)
		r.Get("/findall", productHandler.FindProductAllHandler)
		r.With(authMiddleware).Get("/deletebyId", productHandler.DeleteProductHandler)
		r.With(authMiddleware).Get("/update/{product_id}", productHandler.UpdateProductHandler)
	})

	return router
}
