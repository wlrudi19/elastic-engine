package api

import "github.com/go-chi/chi"

func NewProductRouter(userHandler UserHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/api/users/findbyEmail", userHandler.FindUserHandler)
	router.Post("/api/users/login", userHandler.LoginUserHandler)

	return router
}
