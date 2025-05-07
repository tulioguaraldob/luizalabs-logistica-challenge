package order

import "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"

type UserPurchase struct {
	UserID uint
	Name   string
	Orders []*Purchase
}

type Purchase struct {
	UserID   uint
	Name     string
	Order    *entities.Order
	Products []*entities.OrderProduct
	Total    float64
}
