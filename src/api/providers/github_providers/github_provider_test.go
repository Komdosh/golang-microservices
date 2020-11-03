package github_providers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAuthorizationHeader(t *testing.T) {
	header := getAuthorizationHeader("abc123")

	assert.EqualValues(t, "token abc123", header)
}

func TestDefer(t *testing.T) {
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")

	fmt.Println("function's body")
}
