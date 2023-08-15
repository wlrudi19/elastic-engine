package service

import (
	"context"
	"database/sql"
	"log"

	"github.com/wlrudi19/elastic-engine/app/product/model"
	"github.com/wlrudi19/elastic-engine/app/product/repository"
)

type ProductLogic interface {
	CreateProductLogic(ctx context.Context, req model.CreateProductRequest) error
	FindProductLogic(ctx context.Context, id int) (model.FindProductResponse, error)
}

type productlogic struct {
	ProductRepository repository.ProductRepository
	db                *sql.DB
}

func NewProductLogic(productRepository repository.ProductRepository, db *sql.DB) ProductLogic {
	return &productlogic{
		ProductRepository: productRepository,
		db:                db,
	}
}

func (l *productlogic) CreateProductLogic(ctx context.Context, req model.CreateProductRequest) error {
	log.Printf("[%s][LOGIC] create new product: %s", ctx.Value("productName"), req.Name)

	tx, err := l.db.Begin()

	if err != nil {
		log.Printf("[LOGIC] failed to create product :%v", err)
		return err
	}

	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Amount:      req.Amount,
		Stok:        req.Stok,
	}

	err = l.ProductRepository.CreateProduct(ctx, tx, product)

	if err != nil {
		log.Printf("[LOGIC] failed to create product :%v", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	log.Printf("[%s][LOGIC] created product success with id: %d", ctx.Value("productId"), product.Id)
	return nil
}

func (l *productlogic) FindProductLogic(ctx context.Context, id int) (model.FindProductResponse, error) {
	log.Printf("[%s][LOGIC] find product with id: %d", ctx.Value("productId"), id)
	var product model.FindProductResponse

	tx, err := l.db.Begin()

	if err != nil {
		log.Printf("[LOGIC] failed to find product :%v", err)
		return product, err
	}

	product, err = l.ProductRepository.FindProduct(ctx, tx, id)

	if err != nil {
		log.Printf("[LOGIC] failed to find product :%v", err)
		tx.Rollback()
		return product, err
	}

	tx.Commit()
	log.Printf("[%s][LOGIC] product find successfulyy, id: %d", ctx.Value("productId"), id)
	return product, nil
}
