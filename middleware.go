package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func authenticateRequest(c *gin.Context) {

	// handle preflight requests
	if c.Request.Method == http.MethodOptions {
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	// allow public path access
	if config.PublicPath != "" && strings.HasPrefix(c.Request.URL.Path, config.PublicPath) {
		c.Next()
	}

	// see if public GET is allowed
	if config.AllowPublicGet && c.Request.Method == "GET" {
		c.Next()
		return
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

// getUsernameFromToken fetches the username from the given token.
func getUsernameFromToken(tkn *jwt.Token) (string, error) {
	// get the subject username from the token
	username, err := tkn.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("error getting subject from token: %w", err)
	}
	return username, nil
}

func authenticateUserFromToken(token string) (*User, error) {
	// parse the token string
	tkn, err := parseTokenString(token)
	if err != nil || !tkn.Valid {
		return nil, fmt.Errorf("error parsing token or token is invalid: %w", err)
	}

	// get the subject username from the token
	username, err := getUsernameFromToken(tkn)
	if err != nil {
		return nil, fmt.Errorf("error retrieving username from token: %w", err)
	}

	// find the corresponding user
	user := GetUserByUsername(username)
	if user == nil {
		return nil, fmt.Errorf("no user found corresponding to username %s", username)
	}
	return user, nil
}

func authenticateUsingJWT(c *gin.Context, token string) {
	user, err := authenticateUserFromToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	// store user in context if we need it here
	c.Set("username", user.Username)
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
