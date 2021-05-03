package app

import (
	"github.com/Komdosh/golang-microservices/oauth-api/src/api/controllers/oauth"
	"github.com/Komdosh/golang-microservices/src/api/controllers/health"
)

func mapUrls() {
	router.GET("/health", health.Health)
	router.GET("/oauth/access_token", oauth.CreateAccessToken)
	router.GET("/oauth/access_token/:token_id", oauth.GetAccessToken)
}
