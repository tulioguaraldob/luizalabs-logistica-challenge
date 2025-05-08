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

func Test_Get_UserRepository(t *testing.T) {
	mockUser := &entities.User{
		ID:   10,
		Name: "Tulio Guaraldo",
	}

	tests := []struct {
		description   string
		expectedQuery string
		expectedRows  *sqlmock.Rows
		isErrExpected bool
	}{
		{
			description:   "should return no error",
			expectedQuery: `SELECT * FROM users u WHERE u.id = $1 ORDER BY u.id DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"ID",
				"Name",
			}).AddRow(
				mockUser.ID,
				mockUser.Name,
			),
			isErrExpected: false,
		},
		{
			description:   "should return error",
			expectedQuery: `SELECT * FROM users u WHEREuid IS NULL`,
			expectedRows: sqlmock.NewRows([]string{
				"ID",
				"Name",
			}).AddRow(
				mockUser.ID,
				mockUser.Name,
			),
			isErrExpected: true,
		},
		{
			description:   "should return error on scan rows",
			expectedQuery: `SELECT * FROM users u WHERE u.id = $1 ORDER BY u.id DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"ID",
				"Name",
				"Mocked",
			}).AddRow(
				mockUser.ID,
				mockUser.Name,
				[]byte{},
			),
			isErrExpected: true,
		},
		{
			description:   "should return error on scan rows",
			expectedQuery: `SELECT * FROM users u WHERE u.id = $1 ORDER BY u.id DESC`,
			expectedRows: sqlmock.NewRows([]string{
				"ID",
				"Name",
			}).AddRow(
				mockUser.ID,
				mockUser.Name,
			).RowError(0, sql.ErrNoRows),
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

			userRepository := repositories.NewUserRepository(db)
			u, err := userRepository.Get(mockUser.ID)

			if tt.isErrExpected {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, u)
		})
	}
}

func Test_Add_UserRepository(t *testing.T) {
	mockUser := &entities.User{
		ID:   10,
		Name: "Tulio Guaraldo",
	}

	tests := []struct {
		description   string
		expectedQuery string
		isErrExpected bool
	}{
		{
			description:   "should return no error when add user",
			expectedQuery: `INSERT INTO users (id, name) VALUES ($1, $2)`,
			isErrExpected: false,
		},
		{
			description:   "should return error",
			expectedQuery: `INSERT INTO users (id, name) VALUES ($1, $2)`,
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
				mock.ExpectExec(query).WithArgs(mockUser.ID, mockUser.Name).WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectExec(query).WithArgs(mockUser.ID, mockUser.Name).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			userRepository := repositories.NewUserRepository(db)
			err = userRepository.Add(mockUser)

			if tt.isErrExpected {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
