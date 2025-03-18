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
	Secret         string        `json:"jwtSecret"`
	Issuer         string        `json:"jwtIssuer"`
	Expire         string        `json:"jwtExpire"`
	ProxyPath      string        `json:"proxyPath"`
	PublicPath     string        `json:"publicPath"`
	Target         string        `json:"target"`
	NoCacheHeaders bool          `json:"noCacheHeaders"`
	NoAuth         bool          `json:"noAuth"`
	AllowPublicGet bool          `json:"allowPublicGet"`
	AllowCors      bool          `json:"allowCors"`
	AllowWebLogin  bool          `json:"AllowWebLogin"`
	targetUrl      url.URL       // parsed from Target
	expireDuration time.Duration // parsed from Expire
}{
	Secret:         "my-secret-key",
	Issuer:         "localhost",
	Expire:         "0s",
	ProxyPath:      "/proxy",
	NoCacheHeaders: true,
	NoAuth:         false,
	AllowPublicGet: false,
	AllowCors:      false,
	AllowWebLogin:  false,
}

// parseFlags parses command line flags and updates the `config` variable accordingly.
func parseFlags() {
	flag.StringVar(&config.configFile, "config", "config.json", "Configuration file")
	flag.StringVar(&config.UserFile, "users", "", "File with list of users and passwords, empty creates an 'example' user")
	flag.IntVar(&config.ServerPort, "port", 8080, "Port for server")
	flag.StringVar(&config.ProxyPath, "path", "/proxy", "Path name for proxy server")
	flag.StringVar(&config.Target, "target", "http://localhost:3000/", "Target URL for proxying requests")
	flag.BoolVar(&config.Release, "release", false, "Enable release mode")
	flag.BoolVar(&config.AllowWebLogin, "web", false, "Allow web login")
	flag.Parse()
}

// Check environment variables if secret is not configured yet
func readEnv() {
	if config.Secret == "" {
		config.Secret = os.Getenv(ENV_SECRET)
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

	expireDuration, err := time.ParseDuration(config.Expire)
	if err != nil {
		log.Fatalln("Error parsing JWT expire duration:", err)
	}
	config.expireDuration = expireDuration

	log.Println("Server Port:", config.ServerPort)
	log.Println("Target URL:", config.targetUrl.String())
	log.Println("Proxy Path:", config.ProxyPath)
	log.Println("Public Path:", config.PublicPath)
	if config.expireDuration.Seconds() > 0 {
		log.Println("JWT Expire:", config.expireDuration.String())
	} else {
		log.Println("JWT Expire:", "never")
	}
	log.Println("JWT Secret:", anonymize(config.Secret))
	log.Println("JWT Issuer:", config.Issuer)
	log.Println("Allow CORS:", config.AllowCors)
}
