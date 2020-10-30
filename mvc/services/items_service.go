package services

import (
	"github.com/Komdosh/golang-microservices/mvc/domain"
	"github.com/Komdosh/golang-microservices/mvc/utils"
	"net/http"
)

type itemService struct {
}

var (
	ItemService itemService
)

func (i *itemService) GetItem(userId int64) (*domain.Item, *utils.ApplicationError) {
	return nil, &utils.ApplicationError{
		Message:    "implement me",
		StatusCode: http.StatusInternalServerError,
		Code:       "internal_server_url",
	}
}
