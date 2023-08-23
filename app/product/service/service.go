package service

import (
	"context"
	"log"

	"github.com/wlrudi19/elastic-engine/app/product/model"
	"github.com/wlrudi19/elastic-engine/app/product/repository"
)

type ProductLogic interface {
	CreateProductLogic(ctx context.Context, req model.CreateProductRequest) error
	FindProductLogic(ctx context.Context, id int) (model.FindProductResponse, error)
	FindProductAllLogic(ctx context.Context) ([]model.Product, error)
	DeleteProductLogic(ctx context.Context, id int) error
	UpdateProductLogic(ctx context.Context, id int, fields model.UpdateProductRequest) error
}

type productlogic struct {
	ProductRepository repository.ProductRepository
}

func NewProductLogic(productRepository repository.ProductRepository) ProductLogic {
	return &productlogic{
		ProductRepository: productRepository,
	}
}

func (l *productlogic) CreateProductLogic(ctx context.Context, req model.CreateProductRequest) error {
	log.Printf("[LOGIC] create new product: %s", req.Name)

	tx, err := l.ProductRepository.WithTransaction()

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
	log.Printf("[LOGIC] created product success with id: %d", product.Id)
	return nil
}

func (l *productlogic) FindProductLogic(ctx context.Context, id int) (model.FindProductResponse, error) {
	log.Printf("[LOGIC] find product with id: %d", id)

	var product model.FindProductResponse

	product, err := l.ProductRepository.FindProduct(ctx, id)

	if err != nil {
		log.Printf("[LOGIC] failed to find product :%v", err)
		return product, err
	}

	log.Printf("[LOGIC] product find successfulyy, id: %d", id)
	return product, nil
}

func (l *productlogic) FindProductAllLogic(ctx context.Context) ([]model.Product, error) {
	log.Printf("[LOGIC] find all products")

	var products []model.Product

	products, err := l.ProductRepository.FindProductAll(ctx)

	if err != nil {
		log.Printf("[LOGIC] failed to find product s:%v", err)
		return products, err
	}

	log.Printf("[LOGIC] products find successfulyy")
	return products, nil
}

func (l *productlogic) DeleteProductLogic(ctx context.Context, id int) error {
	log.Printf("[LOGIC] delete product with id: %d", id)

	tx, err := l.ProductRepository.WithTransaction()

	if err != nil {
		log.Printf("[LOGIC] failed to delete product :%v", err)
		return err
	}

	err = l.ProductRepository.DeleteProduct(ctx, tx, id)

	if err != nil {
		log.Printf("[LOGIC] failed to delete product :%v", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	log.Printf("[LOGIC] deleted product success with id: %d", id)
	return nil
}

func (l *productlogic) UpdateProductLogic(ctx context.Context, id int, fields model.UpdateProductRequest) error {
	log.Printf("[LOGIC] update product with id: %d", id)

	tx, err := l.ProductRepository.WithTransaction()

	if err != nil {
		log.Printf("[LOGIC] failed to update product :%v", err)
		return err
	}

	err = l.ProductRepository.UpdateProduct(ctx, tx, id, fields)

	if err != nil {
		log.Printf("[LOGIC] failed to update product :%v", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	log.Printf("[LOGIC] update product success with id: %d", id)
	return nil
}
