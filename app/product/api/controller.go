package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wlrudi19/elastic-engine/app/product/model"
	"github.com/wlrudi19/elastic-engine/app/product/service"
	httputils "github.com/wlrudi19/elastic-engine/helper/http"
)

type ProductHandler interface {
	CreateProductHandler(writer http.ResponseWriter, req *http.Request)
	FindProductHandler(writer http.ResponseWriter, req *http.Request)
}

type producthandler struct {
	ProductLogic service.ProductLogic
}

func NewProductHandler(productLogic service.ProductLogic) ProductHandler {
	return &producthandler{
		ProductLogic: productLogic,
	}
}

func (h *producthandler) CreateProductHandler(writer http.ResponseWriter, req *http.Request) {
	var jsonReq model.CreateProductRequest

	err := json.NewDecoder(req.Body).Decode(&jsonReq)

	if err != nil {
		respon := []httputils.StandardError{
			httputils.StandardError{
				Code:   "400",
				Title:  "Bad Request",
				Detail: "Permintaan tidak valid. Format JSON tidak sesuai",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
		return
	}

	err = h.ProductLogic.CreateProductLogic(context.TODO(), jsonReq)
	if err != nil {
		respon := []httputils.StandardError{
			httputils.StandardError{
				Code:   "500",
				Title:  "Internal server error",
				Detail: "Terjadi kesalahan internal pada server",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
		return
	}

	status := httputils.StandardStatus{
		ErrorCode: 201,
		Message:   "Product created successfully",
	}

	envelope := httputils.StandardEnvelope{
		Status: &status,
		Errors: nil,
	}

	responFix, err := json.Marshal(envelope)
	if err != nil {
		respon := []httputils.StandardError{
			httputils.StandardError{
				Code:   "500",
				Title:  "Internal server error",
				Detail: "Terjadi kesalahan internal pada server",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
		return
	}

	contentType := httputils.NewContentTypeDecorator("application/json")
	httpStatus := http.StatusCreated

	httputils.WriteResponse(writer, responFix, httpStatus, contentType)
}

func (h *producthandler) FindProductHandler(writer http.ResponseWriter, req *http.Request) {
	var jsonReq model.FindProductRequest

	err := json.NewDecoder(req.Body).Decode(&jsonReq)

	if err != nil {
		respon := []httputils.StandardError{
			httputils.StandardError{
				Code:   "400",
				Title:  "Bad Request",
				Detail: "Permintaan tidak valid. Format JSON tidak sesuai",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
		return
	}

	var product = model.FindProductResponse{}
	product, err = h.ProductLogic.FindProductLogic(context.TODO(), jsonReq.Id)
	if err != nil {
		respon := []httputils.StandardError{
			httputils.StandardError{
				Code:   "500",
				Title:  "Internal server error",
				Detail: "Terjadi kesalahan internal pada server",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
		return
	}

	status := httputils.StandardStatus{
		ErrorCode: 200,
		Message:   "Product finding successfully",
	}

	envelope := httputils.StandardEnvelope{
		Status: &status,
		Data:   &product,
	}

	responFix, err := json.Marshal(envelope)
	if err != nil {
		respon := []httputils.StandardError{
			httputils.StandardError{
				Code:   "500",
				Title:  "Internal server error",
				Detail: "Terjadi kesalahan internal pada server",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
		return
	}

	contentType := httputils.NewContentTypeDecorator("application/json")
	httpStatus := http.StatusCreated

	httputils.WriteResponse(writer, responFix, httpStatus, contentType)
}
