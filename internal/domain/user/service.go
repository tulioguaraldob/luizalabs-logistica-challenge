package user

import (
	"mime/multipart"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/entities"
)

type Service interface {
	GetUserByID(userId uint) (*entities.User, error)
	LoadUsersDataFile(file multipart.File) int
}
