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

type PurchaseResponse struct {
	UserID uint             `json:"user_id"`
	Name   string           `json:"name"`
	Orders []*OrderResponse `json:"orders"`
}

type OrderResponse struct {
	OrderID  uint               `json:"order_id"`
	Total    float64            `json:"total"`
	Date     time.Time          `json:"date"`
	Products []*ProductResponse `json:"products"`
}

type ProductResponse struct {
	ProductID uint    `json:"product_id"`
	Value     float64 `json:"value"`
}

func FromOrderToResponse(order *entities.Order) *Response {
	return &Response{
		ID:     order.ID,
		UserID: order.UserID,
		Date:   order.Date,
	}
}

func FromPurchaseToResponse(purchase *Purchase) *PurchaseResponse {
	productsRes := make([]*ProductResponse, 0)
	for _, product := range purchase.Products {
		productRes := &ProductResponse{
			ProductID: product.ID,
			Value:     product.Value,
		}

		productsRes = append(productsRes, productRes)
	}

	// ordersRes := make([]*OrderResponse, 0)
	// orderRes := &OrderResponse{
	// 	OrderID:  purchase.Order.ID,
	// 	Total:    purchase.Total,
	// 	Date:     purchase.Order.Date,
	// 	Products: productsRes,
	// }

	res := &PurchaseResponse{
		UserID: purchase.UserID,
		Name:   purchase.Name,
		Orders: []*OrderResponse{
			{
				OrderID:  purchase.Order.ID,
				Total:    purchase.Total,
				Date:     purchase.Order.Date,
				Products: productsRes,
			},
		},
	}

	return res
}
