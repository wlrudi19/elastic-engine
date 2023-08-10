package service

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/wlrudi19/elastic-engine/app/product/model"
	"github.com/wlrudi19/elastic-engine/app/product/repository"
	"github.com/wlrudi19/elastic-engine/helper"
)

type ProductService interface {
	CreateProductService(ctx context.Context, request model.CreateProductRequest) model.ProductResponse
}

type productservice struct {
	ProductRepository repository.ProductRepository
	DB                *sql.DB
	Validate          *validator.Validate
}

func NewProductService(productRepository repository.ProductRepository, DB *sql.DB, validate *validator.Validate) ProductService {
	return &productservice{
		ProductRepository: productRepository,
		DB:                DB,
		Validate:          validate,
	}
}

func (ps *productservice) CreateProductService(ctx context.Context, request model.CreateProductRequest) model.ProductResponse {
	err := ps.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := ps.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	product := model.Product{
		Name:        request.Name,
		Description: request.Description,
		Amount:      request.Amount,
		Stok:        request.Stok,
	}

	product, err = ps.ProductRepository.CreateProduct(ctx, tx, product)
	helper.PanicIfError(err)

	return model.ToProductResponse(product)
}
