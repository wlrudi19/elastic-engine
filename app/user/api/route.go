package api

import (
	"github.com/go-chi/chi"
	"github.com/wlrudi19/elastic-engine/infrastructure/middlewares"
)

func NewUserRouter(userHandler UserHandler) *chi.Mux {
	router := chi.NewRouter()

	authMiddleware := middlewares.NewAuth()

	router.Route("/api/users", func(r chi.Router) {
		r.With(authMiddleware.Authenticate).Get("/findbyEmail", userHandler.FindUserHandler)
		r.Post("/login", userHandler.LoginUserHandler)
	})

	return router
}
