package user

import "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"

type Repository interface {
	Get(id uint) (*entities.User, error)
	Add(user *entities.User) error
	AddUserOrder(userId uint, orders ...*entities.Order) error
}
