package order

import (
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

type Repository interface {
	Get(id uint) (*entities.Order, error)
	GetByInterval(startDate, endDate time.Time) ([]*entities.Order, error)
	Add(order *entities.Order) error
}
