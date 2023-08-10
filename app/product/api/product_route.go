package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wlrudi19/elastic-engine/helper"
)

func NewProductRouter(productController ProductController) *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/products/create", productController.CreateProductController)

	router.PanicHandler = helper.ErrorHandler

	return router
}
