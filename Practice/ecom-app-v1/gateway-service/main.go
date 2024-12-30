package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// this should be the first thing to load all the environment data
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error loading .env file")
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "80"
	}
	router := gin.New()
	router.Any("/*path", Handler)

	err = router.Run(":" + appPort)
	if err != nil {
		log.Panic(err)
	}

}

func Handler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}
