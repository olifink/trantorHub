package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func authMiddleware(c *gin.Context) {
	// allow unauthenticated GET requests if configured
	if config.AllowGet && c.Request.Method == "GET" {
		c.Next()
	}

	// otherwise validate JWT
	if !isValidToken(c.GetHeader("Authorization")) {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}
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
	readUsers()
	runServer()
}
