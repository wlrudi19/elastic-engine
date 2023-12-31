package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/wlrudi19/elastic-engine/app/product/model"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, tx *sql.Tx, product model.Product) error
	FindProduct(ctx context.Context, id int) (model.FindProductResponse, error)
	FindProductAll(ctx context.Context) ([]model.Product, error)
	DeleteProduct(ctx context.Context, tx *sql.Tx, id int) error
	UpdateProduct(ctx context.Context, tx *sql.Tx, id int, fields model.UpdateProductRequest) error
	WithTransaction() (*sql.Tx, error)
}

type productrepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productrepository{
		db: db,
	}
}

func (pr *productrepository) WithTransaction() (*sql.Tx, error) {
	tx, err := pr.db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (pr *productrepository) CreateProduct(ctx context.Context, tx *sql.Tx, product model.Product) error {
	log.Printf("[QUERY] creating product: %s", product.Name)

	var id int
	sql := "insert into products (name,description,amount,stok) values ($1, $2, $3, $4) RETURNING id"
	err := tx.QueryRowContext(ctx, sql, product.Name, product.Description, product.Amount, product.Stok).Scan(&id)

	if err != nil {
		log.Printf("[QUERY] failed insert into database :%v", err)
		return err
	}

	product.Id = int(id)
	return nil
}

func (pr *productrepository) FindProduct(ctx context.Context, id int) (model.FindProductResponse, error) {
	log.Printf("[QUERY] finding product with id: %d", id)

	var product model.FindProductResponse

	sql := "select name, description, amount, stok from products p where deleted_on isnull and id = $1"
	err := pr.db.QueryRowContext(ctx, sql, id).Scan(
		&product.Name,
		&product.Description,
		&product.Amount,
		&product.Stok,
	)

	if err != nil {
		log.Printf("[QUERY] failed to finding product, %v", err)
		return product, err
	}

	return product, nil
}

func (pr *productrepository) FindProductAll(ctx context.Context) ([]model.Product, error) {
	log.Printf("[QUERY] find all products")

	sql := "select id, name, description, amount, stok from products where deleted_on isnull"
	rows, err := pr.db.QueryContext(ctx, sql)

	if err != nil {
		log.Printf("[QUERY]] failed to finding products, %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Amount,
			&product.Stok,
		)

		if err != nil {
			log.Fatalf("[QUERY] failed to finding product row: %v", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (pr *productrepository) DeleteProduct(ctx context.Context, tx *sql.Tx, id int) error {
	log.Printf("[QUERY] deleting product with id: %d", id)

	var deletedOn sql.NullTime

	checkSQL := "SELECT deleted_on FROM products WHERE id = $1"
	err := tx.QueryRowContext(ctx, checkSQL, id).Scan(&deletedOn)

	if err != nil {
		log.Printf("[QUERY] failed to deleting product: %v", err)
		return err
	}

	if deletedOn.Valid {
		err = errors.New("product has been deleted before")
		log.Printf("[QUERY] failed to deleting product: %v", err)
		return err
	}

	sql := "update products set deleted_on = now() where id = $1"
	_, err = tx.ExecContext(ctx, sql, id)

	if err != nil {
		log.Printf("[QUERY] failed deleting product, %v", err)
		return err
	}

	return nil
}

func (pr *productrepository) UpdateProduct(ctx context.Context, tx *sql.Tx, id int, fields model.UpdateProductRequest) error {
	log.Printf("[QUERY] updating product with id: %d", id)

	updateBuilder := squirrel.Update("products").
		Where(squirrel.Eq{"id": id})

	if fields.Name != nil {
		updateBuilder = updateBuilder.Set("name", *fields.Name)
	}
	if fields.Description != nil {
		updateBuilder = updateBuilder.Set("description", *fields.Description)
	}
	if fields.Amount != nil {
		updateBuilder = updateBuilder.Set("amount", *fields.Amount)
	}
	if fields.Stok != nil {
		updateBuilder = updateBuilder.Set("stok", *fields.Stok)
	}

	query, args, err := updateBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, query, args...)

	if err != nil {
		log.Printf("[QUERY] failed to update product, %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[QUERY] failed to update product, %v", err)
		return err
	}

	if rowsAffected == 0 {
		err := errors.New("sql: no rows in result set")
		log.Printf("[QUERY] product not found, %v", err)
		return err
	}

	return nil
}
