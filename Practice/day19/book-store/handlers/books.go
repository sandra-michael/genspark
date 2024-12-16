package handlers

import (
	"book-store/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	time.Sleep(10 * time.Second)
	c.JSON(200, gin.H{"message": "pong"})
}

func CreateTable(c *gin.Context) {

	conn, err := models.NewConn()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = conn.CreateBookTable(c.Request.Context())
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusCreated, "Book Table")

}

func Insertfunc(c *gin.Context) {
	var newBook models.NewBook

	// Call BindJSON to bind the received JSON to
	// newBooks.
	if err := c.BindJSON(&newBook); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Converting json to struct error "})
		return
	}
	conn, err := models.NewConn()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error forming a new connection"})
		return
	}

	book, err := conn.InsertBook(c.Request.Context(), newBook)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error during insert "})
		return
	}
	c.IndentedJSON(http.StatusCreated, book)
	return

}
