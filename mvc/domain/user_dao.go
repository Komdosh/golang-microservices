package domain

import (
	"fmt"
	"github.com/Komdosh/golang-microservices/mvc/utils"
	"net/http"
)

var (
	users = map[int64]*User{
		123: &User{Id: 123, FirstName: "Andrey", Lastname: "Tabakov", Email: "at@email.com"},
	}
)

func GetUser(userId int64) (*User, *utils.ApplicationError) {
	if user := users[userId]; user != nil {
		return user, nil
	}

	return nil, &utils.ApplicationError{
		Message:    fmt.Sprintf("user %v was not found", userId),
		StatusCode: http.StatusNotFound,
		Code:       "not_found",
	}
}
