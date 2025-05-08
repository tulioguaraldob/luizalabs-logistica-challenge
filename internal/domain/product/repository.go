package product

import "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"

type Repository interface {
	Add(product *entities.Product) error
}
