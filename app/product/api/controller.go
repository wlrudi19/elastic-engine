package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/wlrudi19/elastic-engine/app/product/model"
	"github.com/wlrudi19/elastic-engine/app/product/service"
	httputils "github.com/wlrudi19/elastic-engine/helper/http"
)

type ProductHandler interface {
	CreateProductHandler(writer http.ResponseWriter, req *http.Request)
	FindProductHandler(writer http.ResponseWriter, req *http.Request)
	FindProductAllHandler(writer http.ResponseWriter, req *http.Request)
	DeleteProductAllHandler(writer http.ResponseWriter, req *http.Request)
	UpdateProductHandler(writer http.ResponseWriter, req *http.Request)
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
			{
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
			{
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
			{
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
	var jsonReq model.ProductRequest

	err := json.NewDecoder(req.Body).Decode(&jsonReq)

	if err != nil {
		respon := []httputils.StandardError{
			{
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
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			respon := []httputils.StandardError{
				{
					Code:   "404",
					Title:  "Not found",
					Detail: "Product not found",
					Object: httputils.ErrorObject{},
				},
			}
			httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
			return
		}

		respon := []httputils.StandardError{
			{
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
			{
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
	httpStatus := http.StatusOK

	httputils.WriteResponse(writer, responFix, httpStatus, contentType)
}

func (h *producthandler) FindProductAllHandler(writer http.ResponseWriter, req *http.Request) {
	if req.ContentLength != 0 {
		respon := []httputils.StandardError{
			{
				Code:   "400",
				Title:  "Bad Request",
				Detail: "Tidak boleh ada inputan di body JSON",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
		return
	}

	var products []model.Product
	products, err := h.ProductLogic.FindProductAllLogic(context.TODO())

	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			respon := []httputils.StandardError{
				{
					Code:   "404",
					Title:  "Not found",
					Detail: "Product not found",
					Object: httputils.ErrorObject{},
				},
			}
			httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
			return
		}

		respon := []httputils.StandardError{
			{
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
		Message:   "Products finding successfully",
	}

	envelope := httputils.StandardEnvelope{
		Status: &status,
		Data:   &products,
	}

	responFix, err := json.Marshal(envelope)
	if err != nil {
		respon := []httputils.StandardError{
			{
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
	httpStatus := http.StatusOK

	httputils.WriteResponse(writer, responFix, httpStatus, contentType)
}

func (h *producthandler) DeleteProductAllHandler(writer http.ResponseWriter, req *http.Request) {
	var jsonReq model.ProductRequest

	err := json.NewDecoder(req.Body).Decode(&jsonReq)

	if err != nil {
		respon := []httputils.StandardError{
			{
				Code:   "400",
				Title:  "Bad Request",
				Detail: "Permintaan tidak valid. Format JSON tidak sesuai",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
		return
	}

	err = h.ProductLogic.DeleteProductLogic(context.TODO(), jsonReq.Id)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			respon := []httputils.StandardError{
				{
					Code:   "404",
					Title:  "Not found",
					Detail: "Product not found",
					Object: httputils.ErrorObject{},
				},
			}
			httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
			return
		} else if strings.Contains(err.Error(), "product has been deleted before") {
			respon := []httputils.StandardError{
				{
					Code:   "400",
					Title:  "Bad Request",
					Detail: "Product telah dihapus sebelumnya",
					Object: httputils.ErrorObject{},
				},
			}
			httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
			return
		}

		respon := []httputils.StandardError{
			{
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
		Message:   "Product deleted successfully",
	}

	envelope := httputils.StandardEnvelope{
		Status: &status,
		Errors: nil,
	}

	responFix, err := json.Marshal(envelope)
	if err != nil {
		respon := []httputils.StandardError{
			{
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
	httpStatus := http.StatusOK

	httputils.WriteResponse(writer, responFix, httpStatus, contentType)
}

func (h *producthandler) UpdateProductHandler(writer http.ResponseWriter, req *http.Request) {
	var jsonReq model.UpdateProductRequest

	err := json.NewDecoder(req.Body).Decode(&jsonReq)

	if err != nil {
		respon := []httputils.StandardError{
			{
				Code:   "400",
				Title:  "Bad Request",
				Detail: "Permintaan tidak valid. Format JSON tidak sesuai",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
		return
	}

	productId, err := strconv.Atoi(chi.URLParam(req, "product_id"))

	if err != nil {
		respon := []httputils.StandardError{
			{
				Code:   "400",
				Title:  "Bad Request",
				Detail: "Permintaan tidak valid. Format JSON tidak sesuai",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
		return
	}

	err = h.ProductLogic.UpdateProductLogic(context.TODO(), productId, jsonReq)

	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			respon := []httputils.StandardError{
				{
					Code:   "404",
					Title:  "Not found",
					Detail: "Product not found",
					Object: httputils.ErrorObject{},
				},
			}
			httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
			return
		}

		respon := []httputils.StandardError{
			{
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
		Message:   "Product updated successfully",
	}

	envelope := httputils.StandardEnvelope{
		Status: &status,
		Errors: nil,
	}

	responFix, err := json.Marshal(envelope)
	if err != nil {
		respon := []httputils.StandardError{
			{
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
	httpStatus := http.StatusOK

	httputils.WriteResponse(writer, responFix, httpStatus, contentType)
}
