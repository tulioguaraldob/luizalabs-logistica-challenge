package errors

import "errors"

var (
	ErrEmptyOrders  error = errors.New("there are no orders to be add to user")
	ErrUserNotFound error = errors.New("user does not exist")
)
