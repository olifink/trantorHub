package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func loginHandler(c *gin.Context) {
	// TODO: Implement loginHandler
}

func registerHandler(c *gin.Context) {
	// TODO: Implement registerHandler
}

func authMiddleware(c *gin.Context) {
	// Example: Validate JWT
	if !isValidToken(c.GetHeader("Authorization")) {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}
	c.Next()
}

func isValidToken(header string) bool {
	return true
}

func profileHandler(c *gin.Context) {
	// TODO: Implement profileHandler
}

var jwtKey = []byte("my_secret_key")

func generateToken(c *gin.Context) {
	// Set expiration time of the token
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Issuer:    "trantor",
		Subject:   "someuser",
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(200, gin.H{"token": tokenString})
}

func validateToken(c *gin.Context) {
	// Parse the token
	tknStr := c.GetHeader("Authorization")
	println(tknStr)

	claims := &jwt.RegisteredClaims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			c.JSON(403, gin.H{"error": "Invalid token signature"})
			return
		}
		c.JSON(400, gin.H{"error": "Invalid token"})
		return
	}
	if !tkn.Valid {
		c.JSON(403, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(200, gin.H{"status": "Token is valid"})
}

func main() {
	r := gin.Default()

	// Set up routes for user management
	r.POST("/login", loginHandler)
	r.POST("/register", registerHandler)
	r.GET("/profile", authMiddleware, profileHandler)

	// Set up routes for token management
	r.POST("/token/generate", generateToken)
	r.GET("/token/validate", validateToken)

	// Run the server
	r.Run() // listens and serves on 0.0.0.0:8080 by default
}
