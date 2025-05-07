package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

const (
	getProductQuery    string = `SELECT id, value FROM products p WHERE p.id = $1 ORDER BY p.id DESC`
	createProductQuery string = `INSERT INTO products (id) VALUES ($1)`
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Get(id uint) (*entities.Product, error) {
	product := new(entities.Product)
	row := r.db.QueryRowContext(context.Background(), getProductQuery, id)
	if err := row.Scan(
		&product.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("error scanning product: %w", err)
	}
	return product, nil
}

func (r *productRepository) Add(product *entities.Product) error {
	if _, err := r.db.ExecContext(
		context.Background(),
		createProductQuery,
		product.ID,
	); err != nil {
		return err
	}
	return nil
}
