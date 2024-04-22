package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var config = struct {
	filename       string        // from commandline flag
	ServerPort     int           `json:"serverPort"`
	JwtSecret      string        `json:"jwtSecret"`
	JwtIssuer      string        `json:"jwtIssuer"`
	JwtExpire      string        `json:"jwtExpire"`
	ProxyPath      string        `json:"proxyPath"`
	Target         string        `json:"target"`
	targetUrl      url.URL       // parsed from Target
	expireDuration time.Duration // parsed from JwtExpire
}{
	JwtSecret: "my-secret-key",
	JwtIssuer: "trantor-hub",
	JwtExpire: "5m",
	ProxyPath: "/proxy",
}

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

	log.Println(c.FullPath())
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

	tokenString, err := generateNewToken(username)
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

func profileHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello World!"})
}

func generateTokenHandler(c *gin.Context) {
	tokenString, err := generateNewToken("example")
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(200, gin.H{"token": tokenString})
}

func validateTokenHandler(c *gin.Context) {
	// Parse the token
	tknStr := c.GetHeader("Authorization")
	tkn, err := parseTokenString(tknStr)

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

func parseTokenString(tknStr string) (*jwt.Token, error) {
	tknStr = strings.TrimPrefix(tknStr, "Bearer ")
	claims := &jwt.RegisteredClaims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		secret := []byte(config.JwtSecret)
		return secret, nil
	})
	return tkn, err
}

func proxyHandler(c *gin.Context) {
	// Determine the target URL (modify as needed)
	targetURL := config.targetUrl
	targetURL.Path = c.Param("path")

	log.Println("Proxying to", targetURL.String())

	// Create a new request to the target service, copying the method and the body
	proxyReq, err := http.NewRequest(c.Request.Method, targetURL.String(), c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers (optional, choose which headers to forward)
	for key, value := range c.Request.Header {
		proxyReq.Header[key] = value
	}

	// Create a client and send the request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	// also copy response headers to proxy response
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	// Forward the status code and response body
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func runServer() {

	r := gin.Default()

	// Set up routes for user management
	r.POST("/login", loginHandler)
	r.POST("/register", registerHandler)
	r.GET("/profile", authMiddleware, profileHandler)

	// Set up routes for token management
	r.POST("/token/generate", generateTokenHandler)
	r.GET("/token/validate", validateTokenHandler)

	// authenticated proxy handler
	path := fmt.Sprintf("%s/*path", config.ProxyPath)
	r.Any(path, authMiddleware, proxyHandler)

	// Run the server
	r.Run(fmt.Sprintf(":%d", config.ServerPort)) // listens and serves on defined port
}

func main() {
	commandlineConfig()
	readConfig()
	runServer()
}
