package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupGINRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", Ping)
	return r
}

func HandleBookStore() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("v1/books")
	{
		v1.POST("/createTable", CreateTable)
		v1.POST("/", Insertfunc)
	}
	return r
}
