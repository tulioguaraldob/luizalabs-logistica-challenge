package repositories

import (
	"context"
	"database/sql"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/errors"
)

const (
	getUserQuery          string = `SELECT * FROM users u WHERE u.id = $1 ORDER BY u.id DESC`
	createUserQuery       string = `INSERT INTO users (id, name) VALUES ($1, $2)`
	createUserOrdersQuery string = `INSERT INTO users_orders (user_id, order_id) VALUES ($1, $2)`
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Get(id uint) (*entities.User, error) {
	user := new(entities.User)
	row := r.db.QueryRowContext(context.Background(), getUserQuery, id)
	if err := row.Scan(
		&user.ID,
		&user.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}

		return nil, err
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Add(user *entities.User) error {
	if _, err := r.db.ExecContext(
		context.Background(),
		createUserQuery,
		user.ID,
		user.Name,
	); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) AddUserOrder(userId uint, orders ...*entities.Order) error {
	// Verify if there is any products to be related with the order
	if len(orders) == 0 {
		return errors.ErrEmptyOrders
	}

	for _, order := range orders {
		if _, err := r.db.ExecContext(
			context.Background(),
			createUserOrdersQuery,
			userId,
			order.ID,
		); err != nil {
			return err
		}
	}

	return nil
}
