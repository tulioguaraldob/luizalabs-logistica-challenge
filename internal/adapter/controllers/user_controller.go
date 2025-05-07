package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/errors"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/user"
)

type userController struct {
	service user.Service
}

func NewUserController(service user.Service) *userController {
	return &userController{
		service: service,
	}
}

func (c *userController) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	u, err := c.service.GetUserByID(uint(id))
	if err != nil {
		if err == errors.ErrUserNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := json.Marshal(user.FromUserToResponse(u))
	if err != nil {
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (c *userController) PostUsersData(w http.ResponseWriter, r *http.Request) {
	// Max Memory up to 5 MB (10 * 1024 * 1024)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	file, _, err := r.FormFile("users_data")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	processedLines := c.service.LoadUsersDataFile(file)
	userFileRes := &user.UserFileResponse{
		Message:        "file processed successfully!",
		ProcessedLines: processedLines,
	}

	res, err := json.Marshal(userFileRes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
