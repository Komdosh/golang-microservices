package app

import (
	"github.com/Komdosh/golang-microservices/src/api/log/logger-a"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func StartApp() {
	logger_a.Info("about to map the urls", logger_a.Field("step", "1"), logger_a.Field("status", "pending"))
	mapUrls()
	logger_a.Info("urls successfully mapped", logger_a.Field("step", "2"), logger_a.Field("status", "pending"))

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
