package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func authMiddleware(c *gin.Context) {
	// get authorization header
	tknStr := c.GetHeader("Authorization")

	// allow unauthenticated GET requests if configured
	if tknStr == "" && config.AllowGet && c.Request.Method == "GET" {
		c.Next()
	}

	// otherwise validate JWT
	tkn, err := parseTokenString(tknStr)
	if !isValidToken(tkn, err) {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
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

func runServer() {

	if config.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Set up routes for user management
	r.POST("/login", loginHandler)

	// Set up routes for token development
	if !config.Release {
		r.POST("/token/generate", generateTokenHandler)
		r.GET("/token/validate", validateTokenHandler)
	}

	// authenticated proxy handler
	path := fmt.Sprintf("%s/*path", config.ProxyPath)
	r.Any(path, authMiddleware, proxyHandler)

	// Run the server
	err := r.Run(fmt.Sprintf(":%d", config.ServerPort))
	if err != nil {
		log.Fatal(err)
	} // listens and serves on defined port
}

func main() {
	parseFlags()
	readConfig()
	readEnv()
	readUsers()
	runServer()
}
