package services

import (
	"bufio"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/errors"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/order"
	orderproducts "github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/order_products"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/product"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/user"
)

type userService struct {
	repository              user.Repository
	orderRepository         order.Repository
	productRepository       product.Repository
	orderProductsRepository orderproducts.Repository
}

func NewUserService(
	repository user.Repository,
	orderRepository order.Repository,
	productRepository product.Repository,
	orderProductsRepository orderproducts.Repository,
) *userService {
	return &userService{
		repository:              repository,
		orderRepository:         orderRepository,
		productRepository:       productRepository,
		orderProductsRepository: orderProductsRepository,
	}
}

func (s *userService) GetUserByID(userId uint) (*entities.User, error) {
	user, err := s.repository.Get(userId)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 && user.Name == "" {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (s *userService) LoadUsersDataFile(file multipart.File) int {
	scanner := bufio.NewScanner(file)
	usersData := make([]*user.UserFileData, 0)
	processedLines := 0

	for scanner.Scan() {
		line := scanner.Text()
		processedLines++
		userData := parseUserDataFromLine(line)
		usersData = append(usersData, userData)
	}

	for _, userData := range usersData {
		// Users
		user := &entities.User{
			ID:   userData.UserID,
			Name: userData.UserName,
		}

		if err := s.repository.Add(user); err != nil {
			log.Printf("Failed to create user. Details: %s\n", err.Error())
		}

		// Products
		product := &entities.Product{
			ID: userData.ProductID,
		}

		if err := s.productRepository.Add(product); err != nil {
			log.Printf("Failed to create product. Details: %s\n", err.Error())
		}

		// Orders
		order := &entities.Order{
			ID:     userData.OrderID,
			UserID: user.ID,
			Date:   userData.OrderDate,
		}

		if err := s.orderRepository.Add(order); err != nil {
			log.Printf("Failed to create order. Details: %s\n", err.Error())
		}

		// Orders and products related
		orderProduct := &entities.OrderProduct{
			OrderID:   order.ID,
			ProductID: product.ID,
			Value:     userData.ProductValue,
		}

		if err := s.orderProductsRepository.Add(orderProduct); err != nil {
			log.Printf("Failed to create product to order. Details: %s\n", err.Error())
		}
	}

	return processedLines
}

func parseUserDataFromLine(line string) *user.UserFileData {
	userID := strings.TrimSpace(line[0:10])
	userName := strings.TrimSpace(line[10:55])
	orderID := strings.TrimSpace(line[55:65])
	productID := strings.TrimSpace(line[65:75])
	productValue := strings.TrimSpace(line[75:87])
	orderDate := strings.TrimSpace(line[87:95])

	orderYear := orderDate[0:4]
	orderMonth := orderDate[4:6]
	orderDay := orderDate[6:8]
	date := fmt.Sprintf("%s-%s-%s", orderYear, orderMonth, orderDay)

	parsedUserId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Printf("Failed to parse user id [%s]. Details: %s\n", userID, err.Error())
	}

	parsedOrderId, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		log.Printf("Failed to parse order id [%s]. Details: %s\n", orderID, err.Error())
	}

	parsedProductId, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		log.Printf("Failed to parse product id [%s]. Details: %s\n", productID, err.Error())
	}

	parsedProductValue, err := strconv.ParseFloat(productValue, 64)
	if err != nil {
		log.Printf("Failed to parse product value [%s]. Details: %s\n", productID, err.Error())
	}

	parsedOrderDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		log.Printf("Failed to parse order date [%s]. Details: %s\n", date, err.Error())
	}

	return &user.UserFileData{
		UserID:       uint(parsedUserId),
		UserName:     userName,
		OrderID:      uint(parsedOrderId),
		ProductID:    uint(parsedProductId),
		ProductValue: parsedProductValue,
		OrderDate:    parsedOrderDate,
	}
}
