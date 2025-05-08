package services

import (
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/errors"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/order"
	orderproducts "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/order_products"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/user"
)

type orderService struct {
	repository              order.Repository
	orderProductsRepository orderproducts.Repository
	userRepository          user.Repository
}

func NewOrderService(
	repository order.Repository,
	orderProductsRepository orderproducts.Repository,
	userRepository user.Repository,
) *orderService {
	return &orderService{
		repository:              repository,
		orderProductsRepository: orderProductsRepository,
		userRepository:          userRepository,
	}
}

func (s *orderService) GetOrderById(orderId uint) (*entities.Order, error) {
	order, err := s.repository.Get(orderId)
	if err != nil {
		return nil, err
	}

	if order == nil || order.ID == 0 {
		return nil, errors.ErrOrderNotFound
	}

	return order, nil
}

func (s *orderService) GetOrdersInInterval(startDate, endDate time.Time) ([]*entities.Order, error) {
	if endDate.Before(startDate) {
		return nil, errors.ErrInvalidDateInterval
	}

	orders, err := s.repository.GetByInterval(startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errors.ErrNoOrders
	}

	return orders, nil
}

func (s *orderService) GetAllOrders() ([]*entities.Order, error) {
	orders, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errors.ErrNoOrders
	}

	return orders, nil
}

func (s *orderService) GetAllOrdersProducts() ([]*order.Purchase, error) {
	orders, err := s.GetAllOrders()
	if err != nil {
		return nil, err
	}

	purchases := make([]*order.Purchase, 0)
	for _, o := range orders {
		u, err := s.userRepository.Get(o.UserID)
		if err != nil {
			return nil, err
		}

		ops, err := s.orderProductsRepository.GetByOrderID(o.ID)
		if err != nil {
			return nil, err
		}

		total := 0.0
		for _, p := range ops {
			total += float64(p.Value)
		}

		purchase := &order.Purchase{
			UserID:   u.ID,
			Name:     u.Name,
			Order:    o,
			Products: ops,
			Total:    total,
		}

		purchases = append(purchases, purchase)
	}

	return purchases, nil
}

func (s *orderService) GetOrdersProductsByOrderId(orderId uint) ([]*order.Purchase, error) {
	o, err := s.GetOrderById(orderId)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepository.Get(o.UserID)
	if err != nil {
		return nil, err
	}

	ops, err := s.orderProductsRepository.GetByOrderID(o.ID)
	if err != nil {
		return nil, err
	}

	total := 0.0
	for _, p := range ops {
		total += float64(p.Value)
	}

	purchase := &order.Purchase{
		UserID:   u.ID,
		Name:     u.Name,
		Order:    o,
		Products: ops,
		Total:    total,
	}

	return []*order.Purchase{
		purchase,
	}, nil
}

func (s *orderService) GetOrdersProductsByInterval(startDate, endDate time.Time) ([]*order.Purchase, error) {
	orders, err := s.GetOrdersInInterval(startDate, endDate)
	if err != nil {
		return nil, err
	}

	purchases := make([]*order.Purchase, 0)
	for _, o := range orders {
		u, err := s.userRepository.Get(o.UserID)
		if err != nil {
			return nil, err
		}

		ops, err := s.orderProductsRepository.GetByOrderID(o.ID)
		if err != nil {
			return nil, err
		}

		total := 0.0
		for _, p := range ops {
			total += float64(p.Value)
		}

		purchase := &order.Purchase{
			UserID:   u.ID,
			Name:     u.Name,
			Order:    o,
			Products: ops,
			Total:    total,
		}

		purchases = append(purchases, purchase)
	}

	return purchases, nil
}
