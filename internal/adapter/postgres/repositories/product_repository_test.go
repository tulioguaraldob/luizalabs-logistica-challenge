package repositories_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/adapter/postgres/repositories"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

func Test_Add_ProductRepository(t *testing.T) {
	mockProduct := &entities.Product{
		ID: 150,
	}

	tests := []struct {
		description   string
		expectedQuery string
		isErrExpected bool
	}{
		{
			description:   "should return no error when add product",
			expectedQuery: `INSERT INTO products (id) VALUES ($1)`,
			isErrExpected: false,
		},
		{
			description:   "should return error",
			expectedQuery: `INSERT INTO products (id) VALUES ($1)`,
			isErrExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				panic(err)
			}
			defer db.Close()

			query := regexp.QuoteMeta(tt.expectedQuery)

			if tt.isErrExpected {
				mock.ExpectExec(query).WithArgs(mockProduct.ID).WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectExec(query).WithArgs(mockProduct.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			productRepository := repositories.NewProductRepository(db)
			if err := productRepository.Add(mockProduct); err != nil && tt.isErrExpected {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
