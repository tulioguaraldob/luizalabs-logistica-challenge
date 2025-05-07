package user

import (
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

type UserFileData struct {
	UserID       uint      `json:"user_id"`
	UserName     string    `json:"user_name"`
	OrderID      uint      `json:"order_id"`
	ProductID    uint      `json:"product_id"`
	ProductValue float64   `json:"product_value"`
	OrderDate    time.Time `json:"order_date"`
}

type UserFileResponse struct {
	Message        string `json:"message"`
	ProcessedLines int    `json:"processed_lines"`
}

type Response struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func FromUserToResponse(user *entities.User) *Response {
	return &Response{
		ID:   user.ID,
		Name: user.Name,
	}
}
