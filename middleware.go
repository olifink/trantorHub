package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func authenticateRequest(c *gin.Context) {
	// see if public GET is allowed
	if config.AllowPublicGet && c.Request.Method == "GET" {
		c.Next()
	}

	// Check the authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// Remove "Bearer " from the beginning (if it exists)
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Proceed to authenticate with this token
		authenticateUsingJWT(c, token)
	} else {
		// If authorization header doesn't exist, check the cookies
		if cookieToken, err := c.Cookie("authToken"); err == nil {
			// Proceed to authenticate with the token from the cookie
			authenticateUsingJWT(c, cookieToken)
		} else {
			// No token in cookies either, respond with error
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization token was found"})
		}
	}
}

func authenticateUsingJWT(c *gin.Context, token string) {
	tkn, err := parseTokenString(token)
	if !isValidToken(tkn, err) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// get the subject username from the token
	username, err := tkn.Claims.GetSubject()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
		return
	}

	// find the corresponding user
	user, err := GetUserByUsername(username)
	if user == nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	// store user in context if we need it here
	c.Set("username", username)
	c.Next()
}

func authMiddleware(c *gin.Context) {
	authenticateRequest(c)
}

func nocacheMiddleware(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1
	c.Header("Pragma", "no-cache")                                   // HTTP 1.0
	c.Header("Expires", "0")                                         // Proxies
}
