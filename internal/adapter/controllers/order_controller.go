package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/errors"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/order"
)

type orderController struct {
	service order.Service
}

func NewOrderController(service order.Service) *orderController {
	return &orderController{
		service: service,
	}
}

func (c *orderController) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	// If the ID is passed will always consider the ID first
	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		purchases, err := c.service.GetOrdersProductsByOrderId(uint(id))
		if err != nil {
			if err == errors.ErrOrderNotFound {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		purchaseRes := make([]*order.PurchaseResponse, 0)
		for _, purchase := range purchases {
			purchaseRes = append(purchaseRes, order.FromPurchaseToResponse(purchase))
		}

		res, err := json.Marshal(purchaseRes)
		if err != nil {
			w.WriteHeader(http.StatusFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	if startDateStr != "" {
		startDate, err := time.Parse(time.DateOnly, startDateStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid startDate format. Expected 2006-01-02"))
			return
		}

		endDate := time.Now()
		if endDateStr != "" {
			parsedEndDate, err := time.Parse(time.DateOnly, endDateStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid endDate format. Expected 2006-01-02"))
				return
			}

			endDate = parsedEndDate
		}

		purchases, err := c.service.GetOrdersProductsByInterval(startDate, endDate)
		if err != nil {
			if err == errors.ErrInvalidDateInterval {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			if err == errors.ErrNoOrders {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		purchaseRes := make([]*order.PurchaseResponse, 0)
		for _, purchase := range purchases {
			purchaseRes = append(purchaseRes, order.FromPurchaseToResponse(purchase))
		}

		res, err := json.Marshal(purchaseRes)
		if err != nil {
			w.WriteHeader(http.StatusFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	purchases, err := c.service.GetAllOrdersProducts()
	if err != nil {
		if err == errors.ErrInvalidDateInterval {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if err == errors.ErrNoOrders {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	purchaseRes := make([]*order.PurchaseResponse, 0)
	for _, purchase := range purchases {
		purchaseRes = append(purchaseRes, order.FromPurchaseToResponse(purchase))
	}

	res, err := json.Marshal(purchaseRes)
	if err != nil {
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (c *orderController) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	o, err := c.service.GetOrderById(uint(id))
	if err != nil {
		if err == errors.ErrOrderNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := json.Marshal(order.FromOrderToResponse(o))
	if err != nil {
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
