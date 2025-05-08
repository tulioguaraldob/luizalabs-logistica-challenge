package repositories_test

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/adapter/postgres/repositories"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

func Test_Get_OrderRepository(t *testing.T) {
	mockOrder := &entities.Order{
		ID:     1,
		UserID: 10,
		Date:   time.Now(),
	}

	tests := []struct {
		description   string
		expectedQuery string
		expectedRows  *sqlmock.Rows
		isErrExpected bool
	}{
		{
			description:   "should return no error",
			expectedQuery: `SELECT * FROM orders o WHERE o.id = $1 ORDER BY o.id DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
			}).AddRow(
				mockOrder.ID,
				mockOrder.UserID,
				mockOrder.Date,
			),
			isErrExpected: false,
		},
		{
			description:   "should return error",
			expectedQuery: `SELECT * FROM orders o WHEREoid IS NULL`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
			}).AddRow(
				mockOrder.ID,
				mockOrder.UserID,
				mockOrder.Date,
			),
			isErrExpected: true,
		},
		{
			description:   "should return error on scan rows",
			expectedQuery: `SELECT * FROM orders o WHERE o.id = $1 ORDER BY o.id DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
				"mocked",
			}).AddRow(
				mockOrder.ID,
				mockOrder.UserID,
				mockOrder.Date,
				[]byte{},
			),
			isErrExpected: true,
		},
		{
			description:   "should return error on no rows",
			expectedQuery: `SELECT * FROM orders o WHERE o.id = $1 ORDER BY o.id DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"ID",
				"user_id",
				"date",
			}).RowError(0, sql.ErrNoRows),
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
			mock.ExpectQuery(query).WithArgs(mockOrder.ID).WillReturnRows(tt.expectedRows)

			orderRepository := repositories.NewOrderRepository(db)
			o, err := orderRepository.Get(mockOrder.ID)

			if tt.isErrExpected {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, o)
			assert.Equal(t, mockOrder.ID, o.ID)
			assert.Equal(t, mockOrder.UserID, o.UserID)
			assert.Equal(t, mockOrder.Date.Truncate(time.Second), o.Date.Truncate(time.Second))
		})
	}
}

func Test_GetByInterval_OrderRepository(t *testing.T) {
	startDate := time.Now().Add(time.Minute * -90)
	endDate := time.Now()

	mockOrders := []*entities.Order{
		{ID: 1, UserID: 10, Date: startDate.Add(time.Minute * 30)},
		{ID: 2, UserID: 20, Date: startDate.Add(time.Minute * 20)},
		{ID: 3, UserID: 10, Date: endDate.Add(-time.Minute * 10)},
	}

	tests := []struct {
		description   string
		expectedQuery string
		expectedRows  *sqlmock.Rows
		isErrExpected bool
	}{
		{
			description:   "should return no error",
			expectedQuery: `SELECT * FROM orders o WHERE o.date >= $1 AND o.date <= $2 ORDER BY o.date DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
			}).AddRow(
				mockOrders[0].ID,
				mockOrders[0].UserID,
				mockOrders[0].Date,
			).AddRow(
				mockOrders[1].ID,
				mockOrders[1].UserID,
				mockOrders[1].Date,
			).AddRow(
				mockOrders[2].ID,
				mockOrders[2].UserID,
				mockOrders[2].Date,
			),
			isErrExpected: false,
		},
		{
			description:   "should return no error and return empty slice if no orders",
			expectedQuery: `SELECT * FROM orders o WHERE o.date >= $1 AND o.date <= $2 ORDER BY o.date DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
			}),
			isErrExpected: false,
		},
		{
			description:   "should return error on query",
			expectedQuery: `SELECT * FROM orders o WHEREo.date > $1 AND o.date <= $2 ORDER BY o.date DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
			}),
			isErrExpected: true,
		},
		{
			description:   "should return error on scan rows",
			expectedQuery: `SELECT * FROM orders o WHERE o.date >= $1 AND o.date <= $2 ORDER BY o.date DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
				"mocked",
			}).AddRow(
				mockOrders[0].ID,
				mockOrders[0].UserID,
				mockOrders[0].Date,
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
			mock.ExpectQuery(query).WithArgs(startDate, endDate).WillReturnRows(tt.expectedRows)

			orderRepository := repositories.NewOrderRepository(db)
			orders, err := orderRepository.GetByInterval(startDate, endDate)

			if tt.isErrExpected {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, orders)
			for i := 0; i < len(orders); i++ {
				assert.Equal(t, mockOrders[i].ID, orders[i].ID)
				assert.Equal(t, mockOrders[i].UserID, orders[i].UserID)
				assert.Equal(t, mockOrders[i].Date, orders[i].Date)
			}
		})
	}
}

func Test_Add_OrderRepository(t *testing.T) {
	mockOrder := &entities.Order{
		ID:     2,
		UserID: 20,
		Date:   time.Now(),
	}

	tests := []struct {
		description   string
		expectedQuery string
		isErrExpected bool
	}{
		{
			description:   "should add order successfully",
			expectedQuery: `INSERT INTO orders (id, user_id, date) VALUES ($1, $2, $3)`,
			isErrExpected: false,
		},
		{
			description:   "should return error on database operation",
			expectedQuery: `INSERT INTO orders (id, user_id, date) VALUES ($1, $2, $3)`,
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
				mock.ExpectExec(query).WithArgs(mockOrder.ID, mockOrder.UserID, mockOrder.Date).WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectExec(query).WithArgs(mockOrder.ID, mockOrder.UserID, mockOrder.Date).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			orderRepository := repositories.NewOrderRepository(db)
			err = orderRepository.Add(mockOrder)

			if tt.isErrExpected {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func Test_GetAll_OrderRepository(t *testing.T) {
	startDate := time.Now().Add(time.Minute * -90)
	endDate := time.Now()

	mockOrders := []*entities.Order{
		{ID: 1, UserID: 10, Date: startDate.Add(time.Minute * 30)},
		{ID: 2, UserID: 20, Date: startDate.Add(time.Minute * 20)},
		{ID: 3, UserID: 10, Date: endDate.Add(-time.Minute * 10)},
	}

	tests := []struct {
		description   string
		expectedQuery string
		expectedRows  *sqlmock.Rows
		isErrExpected bool
	}{
		{
			description:   "should return no error",
			expectedQuery: `SELECT * FROM orders o ORDER BY o.id DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
			}).AddRow(
				mockOrders[0].ID,
				mockOrders[0].UserID,
				mockOrders[0].Date,
			).AddRow(
				mockOrders[1].ID,
				mockOrders[1].UserID,
				mockOrders[1].Date,
			).AddRow(
				mockOrders[2].ID,
				mockOrders[2].UserID,
				mockOrders[2].Date,
			),
			isErrExpected: false,
		},
		{
			description:   "should return no error and return empty slice if no orders",
			expectedQuery: `SELECT * FROM orders o ORDER BY o.id DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
			}),
			isErrExpected: false,
		},
		{
			description:   "should return error on query",
			expectedQuery: `SELECT * FROM orders o ORDER BYoid DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
			}),
			isErrExpected: true,
		},
		{
			description:   "should return error on scan rows",
			expectedQuery: `SELECT * FROM orders o ORDER BY o.id DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"id",
				"user_id",
				"date",
				"mocked",
			}).AddRow(
				mockOrders[0].ID,
				mockOrders[0].UserID,
				mockOrders[0].Date,
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
			mock.ExpectQuery(query).WillReturnRows(tt.expectedRows)

			orderRepository := repositories.NewOrderRepository(db)
			orders, err := orderRepository.GetAll()

			if tt.isErrExpected {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, orders)
			for i := 0; i < len(orders); i++ {
				assert.Equal(t, mockOrders[i].ID, orders[i].ID)
				assert.Equal(t, mockOrders[i].UserID, orders[i].UserID)
				assert.Equal(t, mockOrders[i].Date, orders[i].Date)
			}
		})
	}
}
