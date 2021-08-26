package router

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	router.POST("/print", print)
	router.GET("/send", send)
	return router
}
