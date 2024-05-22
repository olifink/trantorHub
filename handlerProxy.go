package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

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
		c.Header("X-Trantor-Identity", createHash(username))
	}

	// Forward the status code and response body
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
