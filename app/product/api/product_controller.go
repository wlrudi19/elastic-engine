package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wlrudi19/elastic-engine/app/product/model"
	"github.com/wlrudi19/elastic-engine/app/product/service"
	"github.com/wlrudi19/elastic-engine/helper"
)

type ProductController interface {
	CreateProductController(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type productcontroller struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return &productcontroller{
		ProductService: productService,
	}
}

func (pc *productcontroller) CreateProductController(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productCreateRequest := model.CreateProductRequest{}
	helper.ReadFromRequestBody(request, &productCreateRequest)

	categoryResponse := pc.ProductService.CreateProductService
	webResponse := helper.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
