package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/errors"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/order"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services"
	mockorder "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services/mocks/order"
	orderproducts "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services/mocks/order_products"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services/mocks/user"
	"go.uber.org/mock/gomock"
)

func Test_GetOrderById_OrderService(t *testing.T) {
	mockOrderId := 99

	mockOrder := &entities.Order{
		ID:     uint(mockOrderId),
		UserID: 50,
		Date:   time.Now(),
	}

	tests := []struct {
		description string
		setMocks    func(
			mor *mockorder.MockRepository,
			mopr *orderproducts.MockRepository,
			mur *user.MockRepository,
		)
		expectedErr error
	}{
		{
			description: "should return no error when order exists",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					Get(uint(mockOrderId)).
					Return(mockOrder, nil)
			},
			expectedErr: nil,
		},
		{
			description: "should return error when order not found",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					Get(uint(mockOrderId)).
					Return(nil, errors.ErrOrderNotFound)
			},
			expectedErr: errors.ErrOrderNotFound,
		},
		{
			description: "should return error on nil order",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					Get(uint(mockOrderId)).
					Return(nil, nil)
			},
			expectedErr: errors.ErrOrderNotFound,
		},
		{
			description: "should return error on empty order",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					Get(uint(mockOrderId)).
					Return(&entities.Order{}, nil)
			},
			expectedErr: errors.ErrOrderNotFound,
		},
		{
			description: "should return error from repository",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					Get(uint(mockOrderId)).
					Return(nil, assert.AnError)
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mor := mockorder.NewMockRepository(ctrl)
			mopr := orderproducts.NewMockRepository(ctrl)
			mur := user.NewMockRepository(ctrl)

			tt.setMocks(mor, mopr, mur)

			orderService := services.NewOrderService(mor, mopr, mur)

			o, err := orderService.GetOrderById(uint(mockOrderId))
			if err != nil {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedErr, err.Error())
				assert.Nil(t, o)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, o)
			assert.Equal(t, mockOrder.ID, o.ID)
			assert.Equal(t, mockOrder.UserID, o.UserID)
		})
	}
}

func Test_GetOrdersInInterval_OrderService(t *testing.T) {
	startDate := time.Date(2025, 05, 01, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 05, 05, 0, 0, 0, 0, time.UTC)

	mockOrders := []*entities.Order{
		{ID: 1, UserID: 10, Date: time.Date(2025, 05, 02, 10, 0, 0, 0, time.UTC)},
		{ID: 2, UserID: 20, Date: time.Date(2025, 05, 04, 15, 0, 0, 0, time.UTC)},
	}

	tests := []struct {
		description string
		startDate   time.Time
		endDate     time.Time
		setMocks    func(
			mor *mockorder.MockRepository,
			mopr *orderproducts.MockRepository,
			mur *user.MockRepository,
		)
		expectedOrders []*entities.Order
		expectedErr    error
	}{
		{
			description: "should return orders",
			startDate:   startDate,
			endDate:     endDate,
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetByInterval(startDate, endDate).
					Return(mockOrders, nil)
			},
			expectedOrders: mockOrders,
			expectedErr:    nil,
		},
		{
			description: "should return error for invalid date",
			startDate:   endDate,
			endDate:     startDate,
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
			},
			expectedOrders: nil,
			expectedErr:    errors.ErrInvalidDateInterval,
		},
		{
			description: "should return no error and empty orders",
			startDate:   startDate,
			endDate:     endDate,
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetByInterval(startDate, endDate).
					Return(make([]*entities.Order, 0), nil)
			},
			expectedOrders: nil,
			expectedErr:    errors.ErrNoOrders,
		},
		{
			description: "should return error",
			startDate:   startDate,
			endDate:     endDate,
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetByInterval(startDate, endDate).
					Return(nil, assert.AnError)
			},
			expectedOrders: nil,
			expectedErr:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mor := mockorder.NewMockRepository(ctrl)
			mopr := orderproducts.NewMockRepository(ctrl)
			mur := user.NewMockRepository(ctrl)

			tt.setMocks(mor, mopr, mur)

			orderService := services.NewOrderService(mor, mopr, mur)

			orders, err := orderService.GetOrdersInInterval(tt.startDate, tt.endDate)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedErr, err.Error())
				assert.Nil(t, orders)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOrders, orders)
		})
	}
}

func Test_GetAllOrders_OrderService(t *testing.T) {
	mockOrders := []*entities.Order{
		{ID: 1, UserID: 10, Date: time.Now()},
		{ID: 2, UserID: 20, Date: time.Now().Add(time.Hour * 72)},
	}

	tests := []struct {
		description string
		setMocks    func(
			mor *mockorder.MockRepository,
			mopr *orderproducts.MockRepository,
			mur *user.MockRepository,
		)
		expectedOrders []*entities.Order
		expectedErr    error
	}{
		{
			description: "should return all orders",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetAll().
					Return(mockOrders, nil)
			},
			expectedOrders: mockOrders,
			expectedErr:    nil,
		},
		{
			description: "should return no error and no orders",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetAll().
					Return(make([]*entities.Order, 0), nil)
			},
			expectedOrders: nil,
			expectedErr:    errors.ErrNoOrders,
		},
		{
			description: "should return error",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetAll().
					Return(nil, assert.AnError)
			},
			expectedOrders: nil,
			expectedErr:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mor := mockorder.NewMockRepository(ctrl)
			mopr := orderproducts.NewMockRepository(ctrl)
			mur := user.NewMockRepository(ctrl)

			tt.setMocks(mor, mopr, mur)

			orderService := services.NewOrderService(mor, mopr, mur)

			orders, err := orderService.GetAllOrders()

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedErr, err.Error())
				assert.Nil(t, orders)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOrders, orders)
		})
	}
}

func Test_GetAllOrdersProducts_OrderService(t *testing.T) {
	mockOrders := []*entities.Order{
		{ID: 1, UserID: 10, Date: time.Now()},
		{ID: 2, UserID: 20, Date: time.Now().Add(time.Hour * 72)},
	}

	mockUsers := []*entities.User{
		{ID: 10, Name: "Tulio Guaraldo"},
		{ID: 20, Name: "John Rambo"},
	}

	mockOrderProducts := []*entities.OrderProduct{
		{OrderID: 1, ProductID: 100, Value: 9.99},
		{OrderID: 1, ProductID: 101, Value: 5.49},
		{OrderID: 2, ProductID: 200, Value: 20.00},
	}

	expectedPurchases := []*order.Purchase{
		{
			UserID: 10,
			Name:   "Tulio Guaraldo",
			Order:  mockOrders[0],
			Products: []*entities.OrderProduct{
				mockOrderProducts[0],
				mockOrderProducts[1],
			},
			Total: 15.48,
		},
		{
			UserID: 20,
			Name:   "John Rambo",
			Order:  mockOrders[1],
			Products: []*entities.OrderProduct{
				mockOrderProducts[2],
			},
			Total: 20.00,
		},
	}

	tests := []struct {
		description string
		setMocks    func(
			mor *mockorder.MockRepository,
			mopr *orderproducts.MockRepository,
			mur *user.MockRepository,
		)
		expectedPurchases []*order.Purchase
		expectedErr       error
	}{
		{
			description: "should return all order products and user info",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetAll().
					Return(mockOrders, nil)

				mur.
					EXPECT().
					Get(uint(10)).
					Return(mockUsers[0], nil)

				mur.
					EXPECT().
					Get(uint(20)).
					Return(mockUsers[1], nil)

				mopr.
					EXPECT().
					GetByOrderID(uint(1)).
					Return([]*entities.OrderProduct{mockOrderProducts[0], mockOrderProducts[1]}, nil)

				mopr.
					EXPECT().
					GetByOrderID(uint(2)).
					Return([]*entities.OrderProduct{mockOrderProducts[2]}, nil)
			},
			expectedPurchases: expectedPurchases,
			expectedErr:       nil,
		},
		{
			description: "should return error no orders on empty orders",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetAll().
					Return(make([]*entities.Order, 0), nil)
			},
			expectedPurchases: []*order.Purchase{},
			expectedErr:       errors.ErrNoOrders,
		},
		{
			description: "should return error on get all orders",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetAll().
					Return(nil, assert.AnError)
			},
			expectedPurchases: nil,
			expectedErr:       assert.AnError,
		},
		{
			description: "should return error on get user",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetAll().
					Return(mockOrders, nil)

				mur.
					EXPECT().
					Get(uint(10)).
					Return(nil, errors.ErrUserNotFound)
			},
			expectedPurchases: nil,
			expectedErr:       errors.ErrUserNotFound,
		},
		{
			description: "should return error on get gy order id",
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetAll().
					Return(mockOrders, nil)

				mur.
					EXPECT().
					Get(uint(10)).
					Return(mockUsers[0], nil)

				mopr.
					EXPECT().
					GetByOrderID(uint(1)).
					Return(nil, assert.AnError)
			},
			expectedPurchases: nil,
			expectedErr:       assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mor := mockorder.NewMockRepository(ctrl)
			mopr := orderproducts.NewMockRepository(ctrl)
			mur := user.NewMockRepository(ctrl)

			tt.setMocks(mor, mopr, mur)

			orderService := services.NewOrderService(mor, mopr, mur)

			purchases, err := orderService.GetAllOrdersProducts()
			if err != nil {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedErr, err.Error())
				assert.Nil(t, purchases)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expectedPurchases), len(purchases))

			for i, purchase := range purchases {
				assert.Equal(t, tt.expectedPurchases[i].UserID, purchase.UserID)
				assert.Equal(t, tt.expectedPurchases[i].Name, purchase.Name)
				assert.Equal(t, tt.expectedPurchases[i].Order, purchase.Order)
				assert.Equal(t, len(tt.expectedPurchases[i].Products), len(purchase.Products))
			}
		})
	}
}

func Test_GetOrdersProductsByOrderId_OrderService(t *testing.T) {
	mockOrderId := 77

	mockOrder := &entities.Order{
		ID:     uint(mockOrderId),
		UserID: 20,
		Date:   time.Now(),
	}

	mockUser := &entities.User{
		ID:   20,
		Name: "Tulio Guaraldo",
	}

	mockOrderProducts := []*entities.OrderProduct{
		{OrderID: uint(mockOrderId), ProductID: 501, Value: 25.99},
		{OrderID: uint(mockOrderId), ProductID: 502, Value: 12.50},
	}

	expectedPurchase := &order.Purchase{
		UserID: mockUser.ID,
		Name:   mockUser.Name,
		Order:  mockOrder,
		Products: []*entities.OrderProduct{
			mockOrderProducts[0],
			mockOrderProducts[1],
		},
		Total: 38.49,
	}

	tests := []struct {
		description string
		orderId     uint
		setMocks    func(
			mor *mockorder.MockRepository,
			mopr *orderproducts.MockRepository,
			mur *user.MockRepository,
		)
		expectedPurchase *order.Purchase
		expectedErr      error
	}{
		{
			description: "should return order products and user info by order id",
			orderId:     uint(mockOrderId),
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					Get(uint(mockOrderId)).
					Return(mockOrder, nil)

				mur.
					EXPECT().
					Get(mockOrder.UserID).
					Return(mockUser, nil)

				mopr.
					EXPECT().
					GetByOrderID(uint(mockOrderId)).
					Return(mockOrderProducts, nil)
			},
			expectedPurchase: expectedPurchase,
			expectedErr:      nil,
		},
		{
			description: "should return error on order not found",
			orderId:     uint(mockOrderId),
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					Get(uint(mockOrderId)).
					Return(nil, errors.ErrOrderNotFound)
			},
			expectedPurchase: nil,
			expectedErr:      errors.ErrOrderNotFound,
		},
		{
			description: "should return error on user not found",
			orderId:     uint(mockOrderId),
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					Get(uint(mockOrderId)).
					Return(mockOrder, nil)

				mur.
					EXPECT().
					Get(mockOrder.UserID).
					Return(nil, errors.ErrUserNotFound)
			},
			expectedPurchase: nil,
			expectedErr:      errors.ErrUserNotFound,
		},
		{
			description: "should return error on order products not found",
			orderId:     uint(mockOrderId),
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					Get(uint(mockOrderId)).
					Return(mockOrder, nil)

				mur.
					EXPECT().
					Get(mockOrder.UserID).
					Return(mockUser, nil)

				mopr.
					EXPECT().
					GetByOrderID(uint(mockOrderId)).
					Return(nil, assert.AnError)
			},
			expectedPurchase: nil,
			expectedErr:      assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mor := mockorder.NewMockRepository(ctrl)
			mopr := orderproducts.NewMockRepository(ctrl)
			mur := user.NewMockRepository(ctrl)

			tt.setMocks(mor, mopr, mur)

			orderService := services.NewOrderService(mor, mopr, mur)

			purchases, err := orderService.GetOrdersProductsByOrderId(tt.orderId)
			if err != nil {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedErr, err.Error())
				assert.Nil(t, purchases)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, purchases, 1)
			assert.Equal(t, tt.expectedPurchase.UserID, purchases[0].UserID)
			assert.Equal(t, tt.expectedPurchase.Name, purchases[0].Name)
			assert.Equal(t, tt.expectedPurchase.Order, purchases[0].Order)
			assert.Equal(t, len(tt.expectedPurchase.Products), len(purchases[0].Products))
		})
	}
}

func Test_GetOrdersProductsByInterval_OrderService(t *testing.T) {
	startDate := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 5, 5, 23, 59, 59, 0, time.UTC)

	mockOrders := []*entities.Order{
		{ID: 1, UserID: 10, Date: time.Date(2025, time.May, 2, 10, 0, 0, 0, time.UTC)},
		{ID: 2, UserID: 20, Date: time.Date(2025, time.May, 4, 15, 30, 0, 0, time.UTC)},
	}

	mockUsers := []*entities.User{
		{ID: 10, Name: "Tulio Guaraldo"},
		{ID: 20, Name: "John Rambo"},
	}

	mockOrderProducts := []*entities.OrderProduct{
		{OrderID: 1, ProductID: 100, Value: 19.99},
		{OrderID: 1, ProductID: 101, Value: 7.50},
		{OrderID: 2, ProductID: 200, Value: 35.00},
	}

	expectedPurchases := []*order.Purchase{
		{
			UserID: 10,
			Name:   "Tulio Guaraldo",
			Order:  mockOrders[0],
			Products: []*entities.OrderProduct{
				mockOrderProducts[0],
				mockOrderProducts[1],
			},
			Total: 27.49,
		},
		{
			UserID: 20,
			Name:   "John Rambo",
			Order:  mockOrders[1],
			Products: []*entities.OrderProduct{
				mockOrderProducts[2],
			},
			Total: 35.00,
		},
	}

	tests := []struct {
		description string
		startDate   time.Time
		endDate     time.Time
		setMocks    func(
			mor *mockorder.MockRepository,
			mopr *orderproducts.MockRepository,
			mur *user.MockRepository,
		)
		expectedPurchases []*order.Purchase
		expectedErr       error
	}{
		{
			description: "should return order products and user info",
			startDate:   startDate,
			endDate:     endDate,
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetByInterval(startDate, endDate).
					Return(mockOrders, nil)

				mur.
					EXPECT().
					Get(uint(10)).
					Return(mockUsers[0], nil)

				mur.
					EXPECT().
					Get(uint(20)).
					Return(mockUsers[1], nil)

				mopr.
					EXPECT().
					GetByOrderID(uint(1)).
					Return([]*entities.OrderProduct{mockOrderProducts[0], mockOrderProducts[1]}, nil)

				mopr.
					EXPECT().
					GetByOrderID(uint(2)).
					Return([]*entities.OrderProduct{mockOrderProducts[2]}, nil)
			},
			expectedPurchases: expectedPurchases,
			expectedErr:       nil,
		},
		{
			description: "should return error on invalid date",
			startDate:   endDate,
			endDate:     startDate,
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
			},
			expectedPurchases: nil,
			expectedErr:       errors.ErrInvalidDateInterval,
		},
		{
			description: "should return error on no orders found",
			startDate:   startDate.AddDate(0, 1, 0),
			endDate:     endDate.AddDate(0, 1, 0),
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetByInterval(startDate.AddDate(0, 1, 0), endDate.AddDate(0, 1, 0)).
					Return(nil, errors.ErrNoOrders)
			},
			expectedPurchases: nil,
			expectedErr:       errors.ErrNoOrders,
		},
		{
			description: "should return error on get orders",
			startDate:   startDate,
			endDate:     endDate,
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetByInterval(startDate, endDate).
					Return(nil, assert.AnError)
			},
			expectedPurchases: nil,
			expectedErr:       assert.AnError,
		},
		{
			description: "should return error on get user",
			startDate:   startDate,
			endDate:     endDate,
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetByInterval(startDate, endDate).
					Return(mockOrders, nil)

				mur.
					EXPECT().
					Get(uint(10)).
					Return(nil, errors.ErrUserNotFound)
			},
			expectedPurchases: nil,
			expectedErr:       errors.ErrUserNotFound,
		},
		{
			description: "should return error on get order products",
			startDate:   startDate,
			endDate:     endDate,
			setMocks: func(
				mor *mockorder.MockRepository,
				mopr *orderproducts.MockRepository,
				mur *user.MockRepository,
			) {
				mor.
					EXPECT().
					GetByInterval(startDate, endDate).
					Return(mockOrders, nil)

				mur.
					EXPECT().
					Get(uint(10)).
					Return(mockUsers[0], nil)

				mopr.
					EXPECT().
					GetByOrderID(uint(1)).
					Return(nil, assert.AnError)
			},
			expectedPurchases: nil,
			expectedErr:       assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mor := mockorder.NewMockRepository(ctrl)
			mopr := orderproducts.NewMockRepository(ctrl)
			mur := user.NewMockRepository(ctrl)

			tt.setMocks(mor, mopr, mur)

			orderService := services.NewOrderService(mor, mopr, mur)

			purchases, err := orderService.GetOrdersProductsByInterval(tt.startDate, tt.endDate)
			if err != nil {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedErr, err.Error())
				assert.Nil(t, purchases)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expectedPurchases), len(purchases))

			for i, purchase := range purchases {
				assert.Equal(t, tt.expectedPurchases[i].UserID, purchase.UserID)
				assert.Equal(t, tt.expectedPurchases[i].Name, purchase.Name)
				assert.Equal(t, tt.expectedPurchases[i].Order, purchase.Order)
				assert.Equal(t, len(tt.expectedPurchases[i].Products), len(purchase.Products))
			}
		})
	}
}
