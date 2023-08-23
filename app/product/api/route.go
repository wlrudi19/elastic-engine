package api

import (
	"github.com/go-chi/chi"
	"github.com/wlrudi19/elastic-engine/infrastructure/middlewares"
)

func NewProductRouter(productHandler ProductHandler) *chi.Mux {
	router := chi.NewRouter()

	authMiddleware := middlewares.NewAuth()

	router.Route("/api/products", func(r chi.Router) {
		r.With(authMiddleware.Authenticate).Get("/create", productHandler.CreateProductHandler)
		r.With(authMiddleware.Authenticate).Get("/findbyId", productHandler.FindProductHandler)
		r.With(authMiddleware.Authenticate).Get("/findall", productHandler.FindProductAllHandler)
		r.With(authMiddleware.Authenticate).Get("/deletebyId", productHandler.DeleteProductHandler)
		r.With(authMiddleware.Authenticate).Get("/update/{product_id}", productHandler.UpdateProductHandler)
	})

	return router
}
