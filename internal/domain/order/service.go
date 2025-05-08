package order

import (
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

type Service interface {
	GetOrderById(id uint) (*entities.Order, error)
	GetOrdersInInterval(startDate, endDate time.Time) ([]*entities.Order, error)
	GetAllOrders() ([]*entities.Order, error)
	GetAllOrdersProducts() ([]*Purchase, error)
	GetOrdersProductsByOrderId(orderId uint) ([]*Purchase, error)
	GetOrdersProductsByInterval(startDate, endDate time.Time) ([]*Purchase, error)
}
