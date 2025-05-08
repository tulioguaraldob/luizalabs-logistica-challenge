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

func Test_Add_OrderProductRepository(t *testing.T) {
	mockOrderProduct := &entities.OrderProduct{
		ID:        10,
		OrderID:   137,
		ProductID: 120,
		Value:     99.99,
	}

	tests := []struct {
		description   string
		expectedQuery string
		isErrExpected bool
	}{
		{
			description:   "should return no error and add order product",
			expectedQuery: `INSERT INTO order_products (order_id, product_id, value) VALUES ($1, $2, $3)`,
			isErrExpected: false,
		},
		{
			description:   "should return error",
			expectedQuery: `INSERT INTO order_products (order_id, product_id, value) VALUES ($1, $2, $3)`,
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
				mock.ExpectExec(query).WithArgs(
					mockOrderProduct.OrderID,
					mockOrderProduct.ProductID,
					mockOrderProduct.Value,
				).WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectExec(query).WithArgs(
					mockOrderProduct.OrderID,
					mockOrderProduct.ProductID,
					mockOrderProduct.Value,
				).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			orderProductRepository := repositories.NewOrderProductRepository(db)
			if err := orderProductRepository.Add(mockOrderProduct); err != nil && tt.isErrExpected {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func Test_GetByOrderID_OrderProductRepository(t *testing.T) {
	mockOrderID := 10
	mockOrderProducts := []*entities.OrderProduct{
		{ID: 1, OrderID: uint(mockOrderID), ProductID: 100, Value: 99.90},
		{ID: 2, OrderID: uint(mockOrderID), ProductID: 200, Value: 19.90},
	}

	tests := []struct {
		description   string
		expectedQuery string
		expectedRows  *sqlmock.Rows
		isErrExpected bool
	}{
		{
			description:   "should return no error and return order products",
			expectedQuery: `SELECT * FROM order_products op WHERE op.order_id = $1`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"order_id",
				"product_id",
				"value",
			}).AddRow(
				mockOrderProducts[0].ID,
				mockOrderProducts[0].OrderID,
				mockOrderProducts[0].ProductID,
				mockOrderProducts[0].Value,
			).AddRow(
				mockOrderProducts[1].ID,
				mockOrderProducts[1].OrderID,
				mockOrderProducts[1].ProductID,
				mockOrderProducts[1].Value,
			),
			isErrExpected: false,
		},
		{
			description:   "should return no error and return empty order products",
			expectedQuery: `SELECT * FROM order_products op WHERE op.order_id = $1`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"order_id",
				"product_id",
				"value",
			}),
			isErrExpected: false,
		},
		{
			description:   "should return error on query",
			expectedQuery: `SELECT * FROM order_products op WHEREoporder_id = $1`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"order_id",
				"product_id",
				"value",
			}),
			isErrExpected: true,
		},
		{
			description:   "should return error on scan rows",
			expectedQuery: `SELECT * FROM order_products op WHERE op.order_id = $1`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"order_id",
				"product_id",
				"value",
				"mocked",
			}).AddRow(
				mockOrderProducts[0].ID,
				mockOrderProducts[0].OrderID,
				mockOrderProducts[0].ProductID,
				mockOrderProducts[0].Value,
				[]byte{},
			),
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
			mock.ExpectQuery(query).WithArgs(uint(mockOrderID)).WillReturnRows(tt.expectedRows)

			orderProductRepository := repositories.NewOrderProductRepository(db)
			orderProducts, err := orderProductRepository.GetByOrderID(uint(mockOrderID))

			if tt.isErrExpected {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, orderProducts)
			for i := 0; i < len(orderProducts); i++ {
				assert.Equal(t, mockOrderProducts[i].ID, orderProducts[i].ID)
				assert.Equal(t, mockOrderProducts[i].OrderID, orderProducts[i].OrderID)
				assert.Equal(t, mockOrderProducts[i].ProductID, orderProducts[i].ProductID)
				assert.Equal(t, mockOrderProducts[i].Value, orderProducts[i].Value)
			}
		})
	}
}
