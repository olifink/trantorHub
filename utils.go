package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

// commandlineConfig parses command line flags and updates the `config` variable accordingly.
func commandlineConfig() {
	flag.StringVar(&config.filename, "config", "config.json", "Configuration file")
	flag.IntVar(&config.ServerPort, "port", 8080, "Port for server")
	flag.StringVar(&config.Target, "target", "http://example.com/", "Target URL for proxying requests")
	flag.Parse()
}

// Open and read the configuration file, decode its contents into the `config` variable,
func readConfig() {
	file, err := os.Open(config.filename)
	if err == nil {
		fmt.Println("Using configuration file", config.filename)
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
		if err != nil {
			fmt.Println("Error decoding config", config.filename, err)
			return
		}
	} else {
		fmt.Println("Using default configuration")
	}

	u, err := url.Parse(config.Target)
	if err != nil {
		log.Panicln("Error parsing target URL:", err)
	}
	config.targetUrl = *u

	if !strings.HasPrefix(config.ProxyPath, "/") {
		log.Panicln("ProxyPath must start with /", config.ProxyPath)
	}

	if len(config.ProxyPath) < 2 {
		log.Panicln("ProxyPath must have a name")
	}

	expireDuration, err := time.ParseDuration(config.JwtExpire)
	if err != nil {
		log.Panicln("Error parsing JWT expire duration:", err)
	}
	config.expireDuration = expireDuration

	log.Println("Server Port:", config.ServerPort)
	log.Println("Target URL:", config.targetUrl.String())
	log.Println("Proxy Path:", config.ProxyPath)
	log.Println("JWT Expire:", config.expireDuration.String())
	log.Println("JWT Secret:", anonymize(config.JwtSecret))
	log.Println("JWT Issuer:", config.JwtIssuer)
}

// Anonymize a part of a sensitive string
func anonymize(s string) string {
	if len(s) > 4 {
		return s[:2] + "****" + s[len(s)-2:]
	} else {
		return "****"
	}
}