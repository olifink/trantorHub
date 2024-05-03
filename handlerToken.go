package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func tokenHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// see if it's a json or form request with login data
	contentType := c.GetHeader("Content-Type")

	// json api login
	if contentType != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	// parse login data
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := identityForUser(loginData.Username, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": user.token})
}
