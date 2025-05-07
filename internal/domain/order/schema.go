package order

import (
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

type Response struct {
	ID     uint      `json:"id"`
	UserID uint      `json:"user_id"`
	Date   time.Time `json:"date"`
}

func FromOrderToResponse(order *entities.Order) *Response {
	return &Response{
		ID:     order.ID,
		UserID: order.UserID,
		Date:   order.Date,
	}
}
