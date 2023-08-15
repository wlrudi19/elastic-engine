package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/wlrudi19/elastic-engine/app/product/model"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, tx *sql.Tx, product model.Product) error
	FindProduct(ctx context.Context, tx *sql.Tx, id int) (model.FindProductResponse, error)
}

type productrepository struct {
}

func NewProductRepository() ProductRepository {
	return &productrepository{}
}

func (pr *productrepository) CreateProduct(ctx context.Context, tx *sql.Tx, product model.Product) error {
	log.Printf("[%s] creating product: %s", ctx.Value("productName"), product.Name)

	var id int
	sql := "insert into products (name,description,amount,stok) values ($1, $2, $3, $4) RETURNING id"
	err := tx.QueryRowContext(ctx, sql, product.Name, product.Description, product.Amount, product.Stok).Scan(&id)

	if err != nil {
		log.Fatalf("failed insert into database :%v", err)
		return err
	}

	product.Id = int(id)
	log.Printf("[%s] created product success with id: %d", ctx.Value("productId"), product.Id)
	return nil
}

func (pr *productrepository) FindProduct(ctx context.Context, tx *sql.Tx, id int) (model.FindProductResponse, error) {
	log.Printf("[%s] find product by id: %d", ctx.Value("productId"), id)

	var product model.FindProductResponse
	sql := "select name, description, amount, stok from products p where deleted_on isnull and id = $1"
	err := tx.QueryRowContext(ctx, sql, id).Scan(
		&product.Name,
		&product.Description,
		&product.Amount,
		&product.Stok,
	)

	if err != nil {
		log.Fatalf("failed to find product :%v", err)
		return product, err
	}

	log.Printf("[%s] success find product by id: %d", ctx.Value("productId"), id)
	return product, nil
}
