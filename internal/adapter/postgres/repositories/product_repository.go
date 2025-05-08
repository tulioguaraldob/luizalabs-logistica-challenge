package repositories

import (
	"context"
	"database/sql"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

const (
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
