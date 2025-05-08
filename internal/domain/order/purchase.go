package order

import "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"

// Purchase is an aggregation that combines the user
// with respective order and its products related with the order
type Purchase struct {
	UserID   uint
	Name     string
	Order    *entities.Order
	Products []*entities.OrderProduct
	Total    float64
}
