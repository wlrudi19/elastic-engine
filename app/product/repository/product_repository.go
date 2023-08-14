package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/wlrudi19/elastic-engine/app/product/model"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, tx *sql.Tx, product model.Product) error
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

	// if err != nil {
	// 	log.Fatalf("failed insert into database :%v", err)
	// 	//return err
	// }

	if err != nil {
		pqErr, isPqError := err.(*pq.Error)
		if isPqError && pqErr.Code == "23505" {
			// Handle duplicate key error
			log.Printf("[%s] Product with name '%s' already exists", ctx.Value("productName"), product.Name)
			return fmt.Errorf("product with name '%s' already exists", product.Name)
		}

		log.Fatalf("failed insert into database: %v", err)
		return err
	}

	product.Id = int(id)
	log.Printf("[%s] created product success with id: %d", ctx.Value("productId"), product.Id)
	return nil
}
