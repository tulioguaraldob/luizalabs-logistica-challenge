package entities

import "time"

type Order struct {
	ID     uint
	UserID uint
	Date   time.Time
}
