package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

const (
	getOrderProductsQuery           string = `SELECT * FROM order_products op WHERE op.id = $1 ORDER BY op.id DESC`
	getOrderProductsByIntervalQuery string = `SELECT * FROM order_products op WHERE op.date >= $1 AND op.date <= $2 ORDER BY op.date ASC`
	createOrderProductsQuery        string = `INSERT INTO order_products (order_id, product_id, value) VALUES ($1, $2, $3)`
	getOrderProductsByOrderIdQuery  string = `SELECT * FROM order_products op WHERE op.order_id = $1`
)

type orderProductRepository struct {
	db *sql.DB
}

func NewOrderProductRepository(db *sql.DB) *orderProductRepository {
	return &orderProductRepository{
		db: db,
	}
}

func (r *orderProductRepository) Get(id uint) (*entities.OrderProduct, error) {
	orderProduct := new(entities.OrderProduct)
	row := r.db.QueryRowContext(context.Background(), getOrderQuery, id)
	if err := row.Scan(
		&orderProduct.ID,
		&orderProduct.OrderID,
		&orderProduct.ProductID,
		&orderProduct.Value,
	); err != nil {
		return nil, err
	}

	return orderProduct, nil
}

func (r *orderProductRepository) GetByInterval(startDate, endDate time.Time) ([]*entities.OrderProduct, error) {
	rows, err := r.db.QueryContext(context.Background(), getOrdersByIntervalQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entities.OrderProduct, 0)
	for rows.Next() {
		orderProduct := new(entities.OrderProduct)
		if err := rows.Scan(
			&orderProduct.ID,
			&orderProduct.OrderID,
			&orderProduct.ProductID,
			&orderProduct.Value,
		); err != nil {
			return nil, err
		}
		orders = append(orders, orderProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderProductRepository) Add(orderProduct *entities.OrderProduct) error {
	if _, err := r.db.ExecContext(
		context.Background(),
		createOrderProductsQuery,
		orderProduct.OrderID,
		orderProduct.ProductID,
		orderProduct.Value,
	); err != nil {
		return err
	}

	return nil
}

func (r *orderProductRepository) GetByOrderID(orderId uint) ([]*entities.OrderProduct, error) {
	rows, err := r.db.QueryContext(context.Background(), getOrderProductsByOrderIdQuery, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entities.OrderProduct, 0)
	for rows.Next() {
		orderProduct := new(entities.OrderProduct)
		if err := rows.Scan(
			&orderProduct.ID,
			&orderProduct.OrderID,
			&orderProduct.ProductID,
			&orderProduct.Value,
		); err != nil {
			return nil, err
		}
		orders = append(orders, orderProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
