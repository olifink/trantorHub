package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
)

// loginHandler handles the login request and generates a JWT token if the authentication is successful
// It uses BasicAuth to authenticate the user, checks if the user exists in the database,
// verifies the password, and generates a token with a 5-minute expiration time if the password is valid.
// The generated token is returned as a JSON response or an appropriate error response if any authentication step fails.
func loginHandler(c *gin.Context) {
	var username, password, redirect string

	log.Println(c.FullPath())
	// check if there is an authorization in the request
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {

		// if not, see if it's a json request with login data
		contentType := c.GetHeader("Content-Type")

		if contentType == "application/json" {
			log.Println("login handling application/json")

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
		} else if contentType == "application/x-www-form-urlencoded" {
			log.Println("login handling application/x-www-form-urlencoded")

			// parse form data
			err := c.Request.ParseForm()
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid form data"})
				return
			}
			username = c.Request.Form.Get("username")
			password = c.Request.Form.Get("password")
			redirect = c.Request.Form.Get("redirect")
		}

	} else {
		log.Println("login handling Basic Auth")
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

	//c.Header("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	maxAge := int(config.expireDuration.Seconds())
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"authToken",
		tokenString,
		maxAge,
		"/",
		"",
		config.Release,
		true,
	)

	if redirect != "" {
		c.Redirect(http.StatusFound, redirect)
	} else {
		c.JSON(200, gin.H{"token": tokenString})
	}
}

func logoutHandler(c *gin.Context) {
	// Clear the auth token cookie
	c.SetCookie(
		"authToken",
		"",
		-1,
		"/",
		"",
		config.Release,
		true,
	)
	if c.Request.ParseForm() == nil {
		redirect := c.Request.Form.Get("redirect")
		c.Redirect(http.StatusFound, redirect)
	} else {
		c.JSON(200, gin.H{"token": nil})
	}
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

	// forward pseudonoimzied identity header for use in downstream system
	username := c.GetString("username")
	if username != "" {
		c.Header("X-Trantor-Identity", anonymize(username))
	}

	// Forward the status code and response body
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
