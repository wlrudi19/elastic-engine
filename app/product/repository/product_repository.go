package repository

import (
	"context"
	"database/sql"

	"github.com/wlrudi19/elastic-engine/app/product/model"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, tx *sql.Tx, product model.Product) (model.Product, error)
}

type product struct {
}

func NewProductRepository() ProductRepository {
	return &product{}
}

func (p *product) CreateProduct(ctx context.Context, tx *sql.Tx, product model.Product) (model.Product, error) {
	SQL := "insert into products (name,description,amount,stok) values ($1, $2, $3, $4) RETURNING id"
	row := tx.QueryRowContext(ctx, SQL, product.Name, product.Description, product.Amount, product.Stok)

	var id int
	if err := row.Scan(&id); err != nil {
		return model.Product{}, err
	}

	product.Id = int(id)
	return product, nil
}
