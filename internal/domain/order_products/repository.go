package orderproducts

import (
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

type Repository interface {
	Get(id uint) (*entities.OrderProduct, error)
	GetByInterval(startDate, endDate time.Time) ([]*entities.OrderProduct, error)
	Add(orderProduct *entities.OrderProduct) error
}
