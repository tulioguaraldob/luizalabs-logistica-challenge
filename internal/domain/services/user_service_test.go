package services_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/errors"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services/mocks/order"
	orderproducts "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services/mocks/order_products"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services/mocks/product"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services/mocks/user"
	"go.uber.org/mock/gomock"
)

func Test_GetUserById_UserService(t *testing.T) {
	mockUserId := 150

	mockUser := &entities.User{
		ID:   uint(mockUserId),
		Name: "Tulio Guaraldo",
	}

	tests := []struct {
		description string
		setMocks    func(
			mur *user.MockRepository,
			mor *order.MockRepository,
			mpr *product.MockRepository,
			mopr *orderproducts.MockRepository,
		)
		expectedErr error
	}{
		{
			description: "should return no error when user exists",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
				mur.
					EXPECT().
					Get(uint(mockUserId)).
					Return(mockUser, nil)
			},
			expectedErr: nil,
		},
		{
			description: "should return error",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
				mur.
					EXPECT().
					Get(uint(mockUserId)).
					Return(nil, errors.ErrUserNotFound)
			},
			expectedErr: errors.ErrUserNotFound,
		},
		{
			description: "should return error when user infos are empty",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
				mur.
					EXPECT().
					Get(uint(mockUserId)).
					Return(&entities.User{
						ID:   0,
						Name: "",
					}, nil)
			},
			expectedErr: errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mur := user.NewMockRepository(ctrl)
			mor := order.NewMockRepository(ctrl)
			mpr := product.NewMockRepository(ctrl)
			mopr := orderproducts.NewMockRepository(ctrl)

			tt.setMocks(
				mur,
				mor,
				mpr,
				mopr,
			)

			userService := services.NewUserService(
				mur,
				mor,
				mpr,
				mopr,
			)

			u, err := userService.GetUserByID(uint(mockUserId))
			if err != nil {
				assert.NotNil(t, err)
				assert.EqualError(t, tt.expectedErr, err.Error())
				return
			}

			assert.NotNil(t, u)
			assert.NoError(t, err)
			assert.Equal(t, mockUser.ID, u.ID)
			assert.Equal(t, mockUser.Name, u.Name)
		})
	}
}

func Test_LoadUsersDataFile_UserService(t *testing.T) {
	mockUser := &entities.User{
		ID:   70,
		Name: "Palmer Prosacco",
	}

	mockProduct := &entities.Product{
		ID: 3,
	}

	mockOrder := &entities.Order{
		ID:     753,
		UserID: 70,
		Date:   time.Date(2021, 3, 8, 0, 0, 0, 0, time.UTC),
	}

	mockOrderProduct := &entities.OrderProduct{
		OrderID:   753,
		ProductID: 3,
		Value:     1836.74,
	}

	mockUsers := []*entities.User{
		mockUser,
		{
			ID:   75,
			Name: "Bobbie Batz",
		},
	}

	mockProducts := []*entities.Product{
		mockProduct,
		{
			ID: 2,
		},
	}

	mockOrders := []*entities.Order{
		mockOrder,
		{
			ID:     798,
			UserID: 75,
			Date:   time.Date(2021, 11, 16, 0, 0, 0, 0, time.UTC),
		},
	}

	mockOrderProducts := []*entities.OrderProduct{
		mockOrderProduct,
		{
			OrderID:   798,
			ProductID: 2,
			Value:     1578.57,
		},
	}

	tests := []struct {
		description string
		mockedFile  string
		setMocks    func(
			mur *user.MockRepository,
			mor *order.MockRepository,
			mpr *product.MockRepository,
			mopr *orderproducts.MockRepository,
		)
		expectedProcessedLines int
	}{
		{
			description: "should process a single line and save parsed data from file",
			mockedFile:  "./mocks/user/mock_single_data_file.txt",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
				mur.EXPECT().Add(mockUser).Return(nil)
				mpr.EXPECT().Add(mockProduct).Return(nil)
				mor.EXPECT().Add(mockOrder).Return(nil)
				mopr.EXPECT().Add(mockOrderProduct).Return(nil)
			},
			expectedProcessedLines: 1,
		},
		{
			description: "should process multiple valid lines and save parsed data from file",
			mockedFile:  "./mocks/user/mock_mult_data_file.txt",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
				for _, mockUser := range mockUsers {
					mur.EXPECT().Add(mockUser).Return(nil)
				}
				for _, mockProduct := range mockProducts {
					mpr.EXPECT().Add(mockProduct).Return(nil)
				}
				for _, mockOrder := range mockOrders {
					mor.EXPECT().Add(mockOrder).Return(nil)
				}
				for _, mockOrderProduct := range mockOrderProducts {
					mopr.EXPECT().Add(mockOrderProduct).Return(nil)
				}
			},
			expectedProcessedLines: 2,
		},
		{
			description: "should process a valid line and handle user error",
			mockedFile:  "./mocks/user/mock_failed_data_file.txt",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
				mur.EXPECT().Add(mockUser).Return(assert.AnError)
				mpr.EXPECT().Add(mockProduct).Return(nil)
				mor.EXPECT().Add(mockOrder).Return(nil)
				mopr.EXPECT().Add(mockOrderProduct).Return(nil)
			},
			expectedProcessedLines: 1,
		},
		{
			description: "should process a valid line and handle product error",
			mockedFile:  "./mocks/user/mock_failed_data_file.txt",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
				mur.EXPECT().Add(mockUser).Return(nil)
				mpr.EXPECT().Add(mockProduct).Return(assert.AnError)
				mor.EXPECT().Add(mockOrder).Return(nil)
				mopr.EXPECT().Add(mockOrderProduct).Return(nil)
			},
			expectedProcessedLines: 1,
		},
		{
			description: "should process a valid line and handle order error",
			mockedFile:  "./mocks/user/mock_failed_data_file.txt",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
				mur.EXPECT().Add(mockUser).Return(nil)
				mpr.EXPECT().Add(mockProduct).Return(nil)
				mor.EXPECT().Add(mockOrder).Return(assert.AnError)
				mopr.EXPECT().Add(mockOrderProduct).Return(nil)
			},
			expectedProcessedLines: 1,
		},
		{
			description: "should process a valid line and handle order products error",
			mockedFile:  "./mocks/user/mock_failed_data_file.txt",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
				mur.EXPECT().Add(mockUser).Return(nil)
				mpr.EXPECT().Add(mockProduct).Return(nil)
				mor.EXPECT().Add(mockOrder).Return(nil)
				mopr.EXPECT().Add(mockOrderProduct).Return(assert.AnError)
			},
			expectedProcessedLines: 1,
		},
		{
			description: "should process an empty file and return zero lines processed",
			mockedFile:  "./mocks/user/mock_empty_data_file.txt",
			setMocks: func(
				mur *user.MockRepository,
				mor *order.MockRepository,
				mpr *product.MockRepository,
				mopr *orderproducts.MockRepository,
			) {
			},
			expectedProcessedLines: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mur := user.NewMockRepository(ctrl)
			mor := order.NewMockRepository(ctrl)
			mpr := product.NewMockRepository(ctrl)
			mopr := orderproducts.NewMockRepository(ctrl)

			tt.setMocks(
				mur,
				mor,
				mpr,
				mopr,
			)

			userService := services.NewUserService(
				mur,
				mor,
				mpr,
				mopr,
			)

			file, err := os.Open(tt.mockedFile)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			processedLines := userService.LoadUsersDataFile(file)
			assert.Equal(t, tt.expectedProcessedLines, processedLines)
		})
	}
}
