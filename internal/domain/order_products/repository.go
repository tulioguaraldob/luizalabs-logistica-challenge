package orderproducts

import (
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

type Repository interface {
	Add(orderProduct *entities.OrderProduct) error
	GetByOrderID(orderId uint) ([]*entities.OrderProduct, error)
}
