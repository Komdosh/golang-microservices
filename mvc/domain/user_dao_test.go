package domain

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetUserNoUserFound(t *testing.T) {
	user, err := GetUser(0)

	assert.Nil(t, user, "We were not expecting a user with id 0")
	assert.NotNil(t, err, "We were expecting an error when user id is 0")
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode, "We were expecting an error when user id is 0")

	assert.EqualValues(t, "not_found", err.Code)
	assert.EqualValues(t, "user 0 doesn't exists", err.Message)
}

func TestGetUserNoError(t *testing.T) {
	user, err := GetUser(123)

	assert.NotNil(t, user, "We were expecting a user with id 123")
	assert.Nil(t, err, "We were not expecting any errors here")

	assert.EqualValues(t, 123, user.Id)
	assert.EqualValues(t, "Andrey", user.FirstName)
	assert.EqualValues(t, "Tabakov", user.Lastname)
	assert.EqualValues(t, "at@email.com", user.Email)
}
