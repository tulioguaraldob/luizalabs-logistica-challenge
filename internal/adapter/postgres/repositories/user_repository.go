package repositories

import (
	"context"
	"database/sql"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/errors"
)

const (
	getUserQuery    string = `SELECT * FROM users u WHERE u.id = $1 ORDER BY u.id DESC`
	createUserQuery string = `INSERT INTO users (id, name) VALUES ($1, $2)`
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
