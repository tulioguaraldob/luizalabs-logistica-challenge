package errors

import "errors"

var (
	ErrEmptyProducts error = errors.New("there are no products to be add to the order")
)
