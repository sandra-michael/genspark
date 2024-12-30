package handlers

import (
	"net/http"
	"user-service/internal/users"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Signup(c *gin.Context) {

	var newUser users.NewUser

	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	err = h.validate.Struct(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide values in correct format"})
		return
	}

	user, err := h.c.InsertUser(c, newUser)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "problem inserting user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
