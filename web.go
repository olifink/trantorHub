package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func webFormHandler(c *gin.Context) {
	if _, err := c.Cookie("authToken"); err == nil {
		// allow logout if user if there is an authToken
		c.HTML(http.StatusOK, "logout.html", gin.H{
			"redirect": "/login",
		})
		return
	}

	// otherwise show login
	c.HTML(http.StatusOK, "login.html", gin.H{
		"redirect": config.ProxyPath,
	})
}
