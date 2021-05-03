package oauth

import (
	"fmt"
	"github.com/Komdosh/golang-microservices/src/api/utils/errors"
)

var (
	tokens = make(map[string]*AccessToken, 0)
)

func (at *AccessToken) Save() errors.ApiError {
	at.AccessToken = fmt.Sprintf("USER_%d", at.UserId)
	tokens[at.AccessToken] = at

	return nil
}

func GetAccessTokenByToken(accessToken string) (*AccessToken, errors.ApiError) {
	at := tokens[accessToken]
	if at == nil {
		return nil, errors.NewNotFoundApiError("no access token was found by given parameters")
	}

	return at, nil
}
