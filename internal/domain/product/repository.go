package product

import "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"

type Repository interface {
	Get(id uint) (*entities.Product, error)
	Add(product *entities.Product) error
}
