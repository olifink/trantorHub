package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const jwtIssuer = "trantor"

var jwtKey = []byte("my_secret_key")

type User struct {
	ID       int
	Username string
	Password string // This should be a hashed password
}

// GetUserByUsername fetches a user by username from the database
func GetUserByUsername(username string) (*User, error) {
	// TODO Database fetching logic here
	// For now, let's assume we get some user or nil if not found
	return &User{ID: 1, Username: "example", Password: "$2a$14$HNOQGnDpfyF/95TT6VToEuyS4NCYKXH1pVlcq9fx9JaC/zBW.cn0i"}, nil // bcrypt hash for "password"
}

// loginHandler handles the login request and generates a JWT token if the authentication is successful
// It uses BasicAuth to authenticate the user, checks if the user exists in the database,
// verifies the password, and generates a token with a 5-minute expiration time if the password is valid.
// The generated token is returned as a JSON response or an appropriate error response if any authentication step fails.
func loginHandler(c *gin.Context) {
	var username, password string

	// check if there is an authorization in the request
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {

		// if not, see if it's a json request with login data
		contentType := c.GetHeader("Content-Type")

		if contentType != "application/json" {
			c.JSON(415, gin.H{"error": "Content-Type must be application/json if no authorization header provided"})
			return
		}

		// parse login data
		var loginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid login data"})
			return
		}
		username = loginData.Username
		password = loginData.Password
	} else {
		// BasicAuth loginHandler
		ok := false
		username, password, ok = c.Request.BasicAuth()
		if !ok {
			c.JSON(401, gin.H{"error": "Authentication required"})
			return
		}
	}

	// Check if the username and password are empty
	if username == "" || password == "" {
		c.JSON(400, gin.H{"error": "Missing username or password"})
		return
	}

	// Check if we know the user
	user, err := GetUserByUsername(username)
	if err != nil || user == nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check that the password matches the hashed password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	// Set expiration time of the token
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Issuer:    jwtIssuer,
		Subject:   username,
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

func generateToken(c *gin.Context) {
	// Set expiration time of the token
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Issuer:    jwtIssuer,
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
