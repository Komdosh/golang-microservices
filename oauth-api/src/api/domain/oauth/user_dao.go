package oauth

import (
	"fmt"
	"github.com/Komdosh/golang-microservices/src/api/utils/errors"
)

const (
	queryGetUserByUsernameAndPassword = "SELECT id, username FROM users where username=? and password=?;"
)

var (
	users = map[string]*User{
		"kom": &User{
			Id:       123,
			Username: "kom",
		},
	}
)

func GetUserByUsernameAndPassword(username string, password string) (*User, errors.ApiError) {
	user := users[username]
	if user == nil {
		return nil, errors.NewNotFoundApiError(fmt.Sprintf("user with username: %v doesn't exists", username))
	}

	return user, nil
}
