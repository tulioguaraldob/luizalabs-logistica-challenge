package errors

import "errors"

var (
	ErrEmptyProducts       error = errors.New("there are no products to be add to the order")
	ErrOrderNotFound       error = errors.New("order does not exist")
	ErrInvalidDateInterval error = errors.New("end date can not be smaller than start date")
	ErrNoOrders            error = errors.New("no orders were found")
)
