package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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

	// token creation
	r.POST("/token", tokenHandler)

	if config.AllowWebLogin {
		r.LoadHTMLGlob("templates/*.html")
		r.GET("/login", webFormHandler)
		r.POST("/login", loginHandler)
		r.POST("/logout", logoutHandler)
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
