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

const ENV_SECRET = "TRANTOR_JWT_SECRET"

var config = struct {
	configFile     string        // from commandline flag
	UserFile       string        `json:"userFile"`
	Release        bool          `json:"releaseMode"`
	ServerPort     int           `json:"serverPort"`
	JwtSecret      string        `json:"jwtSecret"`
	JwtIssuer      string        `json:"jwtIssuer"`
	JwtExpire      string        `json:"jwtExpire"`
	ProxyPath      string        `json:"proxyPath"`
	NoCache        bool          `json:"noCache"`
	AllowGet       bool          `json:"allowGet"`
	AllowCors      bool          `json:"allowCors"`
	Target         string        `json:"target"`
	targetUrl      url.URL       // parsed from Target
	expireDuration time.Duration // parsed from JwtExpire
}{
	JwtSecret: "my-secret-key",
	JwtIssuer: "localhost",
	JwtExpire: "0s",
	ProxyPath: "/proxy",
	NoCache:   true,
	AllowGet:  false,
	AllowCors: true,
}

// parseFlags parses command line flags and updates the `config` variable accordingly.
func parseFlags() {
	flag.StringVar(&config.configFile, "config", "config.json", "Configuration file")
	flag.StringVar(&config.UserFile, "users", "", "File with list of users and passwords, empty creates an 'example' user")
	flag.IntVar(&config.ServerPort, "port", 8080, "Port for server")
	flag.StringVar(&config.Target, "target", "http://example.com/", "Target URL for proxying requests")
	flag.BoolVar(&config.Release, "release", false, "Enable release mode")
	flag.Parse()
}

// Check environment variables if secret is not configured yet
func readEnv() {
	if config.JwtSecret == "" {
		config.JwtSecret = os.Getenv(ENV_SECRET)
	}
}

// Open and read the configuration file, decode its contents into the `config` variable,
func readConfig() {
	file, err := os.Open(config.configFile)
	if err == nil {
		fmt.Println("Using configuration file", config.configFile)
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
		if err != nil {
			fmt.Println("Error decoding config", config.configFile, err)
			return
		}
	} else {
		fmt.Println("Using default configuration")
	}

	u, err := url.Parse(config.Target)
	if err != nil {
		log.Fatalln("Error parsing target URL:", err)
	}
	config.targetUrl = *u

	if !strings.HasPrefix(config.ProxyPath, "/") {
		log.Fatalln("ProxyPath must start with /", config.ProxyPath)
	}

	if len(config.ProxyPath) < 2 {
		log.Fatalln("ProxyPath must have a name")
	}

	expireDuration, err := time.ParseDuration(config.JwtExpire)
	if err != nil {
		log.Fatalln("Error parsing JWT expire duration:", err)
	}
	config.expireDuration = expireDuration

	log.Println("Server Port:", config.ServerPort)
	log.Println("Target URL:", config.targetUrl.String())
	log.Println("Proxy Path:", config.ProxyPath)
	if config.expireDuration.Seconds() > 0 {
		log.Println("JWT Expire:", config.expireDuration.String())
	} else {
		log.Println("JWT Expire:", "never")
	}
	log.Println("JWT Secret:", anonymize(config.JwtSecret))
	log.Println("JWT Issuer:", config.JwtIssuer)
}
