package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func runServer() {

	if config.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	if config.AllowCors {
		r.Use(corsMiddleware())
	}
	if config.NoCacheHeaders {
		r.Use(nocacheMiddleware)
	}

	r.POST("/login", loginHandler)
	r.POST("/logout", logoutHandler)

	// Set up routes for token development
	if !config.Release {
		r.POST("/token/generate", generateTokenHandler)
		r.GET("/token/validate", validateTokenHandler)
	}

	// TODO replace with tempated forms
	if config.AllowWebLogin {
		r.StaticFS("/web", http.Dir("web"))
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
