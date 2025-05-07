package order

import (
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

type Service interface {
	GetOrderById(id uint) (*entities.Order, error)
	GetOrdersInInterval(startDate, endDate time.Time) ([]*entities.Order, error)
	GetAllOrders() ([]*entities.Order, error)
	GetAllOrdersWithProducts() ([]*Purchase, error)
}
