package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

const (
	getOrderQuery            string = `SELECT * FROM orders o WHERE o.id = $1 ORDER BY o.id DESC`
	getOrdersByIntervalQuery string = `SELECT * FROM orders o WHERE o.date >= $1 AND o.date <= $2 ORDER BY o.date ASC`
	createOrderQuery         string = `INSERT INTO orders (id, user_id, date) VALUES ($1, $2, $3)`
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Get(id uint) (*entities.Order, error) {
	order := new(entities.Order)
	row := r.db.QueryRowContext(context.Background(), getOrderQuery, id)
	if err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Date,
	); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *orderRepository) GetByInterval(startDate, endDate time.Time) ([]*entities.Order, error) {
	rows, err := r.db.QueryContext(context.Background(), getOrdersByIntervalQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entities.Order, 0)
	for rows.Next() {
		order := new(entities.Order)
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Date,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) Add(order *entities.Order) error {
	if _, err := r.db.ExecContext(
		context.Background(),
		createOrderQuery,
		order.ID,
		order.UserID,
		order.Date,
	); err != nil {
		return err
	}

	return nil
}
