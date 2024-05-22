[![Go Release, Build for All Platforms](https://github.com/olifink/trantorHub/actions/workflows/release.yml/badge.svg)](https://github.com/olifink/trantorHub/actions/workflows/release.yml)

# TrantorHub Gateway

TrantorHub is a local API Gateway with token-based authentication, file-based user management and direct Http forwarding written with Golang and Gin.

## Brief Introduction

TrantorHub is a useful tool for developing applications utilizing JSON Web Tokens (JWT) and is designed to function as a secure gateway proxy. This project is under active development, and the features, as well as configuration, are subject to change.

As a word of caution, if you plan to utilize it for production eventually, it is strongly recommended to enforce HTTPS via a reverse proxy like Caddy due to its initial design to cater to development setups.

## Functionality

TrantorHub is used as a REST API utilizing JWT token generation for user credentials. You can use the generated JWTs in bearer authentication headers to incorporate authentication to a proxied endpoint API. TrantorHub works for development and using its default configuration as standalone executable directly for a local setup, to more advanced setups via configuration and user files when developing clients using bearer authentication.

## How to Use

We've deliberately kept it simple - start the `trantorHub` without command-line arguments, and it will proxy against `localhost:3000` on `localhost:8080/proxy`. All request methods forwarded are located under the `proxy` subpath. That said, these can only be accessed with a valid token set in the `Authentication: Bearer request` header.


```
./trantorHub
```
```
Server Port: 8080
Target URL: http://localhost:3000/
Proxy Path: /proxy
JWT Expire: never
JWT Secret: my****ey
JWT Issuer: localhost
User 0: example $2****0i
```

## Customizing your use

The TrantorHub comes with various pre-configured values for the JWT parameters, and a test user (`example`, with password: `password`) is also available to begin with. If you'd like to customize your experience, you can use any of the following commands:

| flag             | description                                                            |
|------------------|------------------------------------------------------------------------|
| `-port <int>`    | Port for server (default 8080)                                         |
| `-path <string>` | Path name for proxy server (default "/proxy")                          |
| `-target <url>`  | Target URL for proxying requests (default "http://localhost:3000/")    |
| `-config <file>` | Configuration file (default "config.json")                             |
| `-users <file>`  | File with list of users and passwords, empty creates an 'example' user |

Use these commands as such:

```
./trantorHub -port 9090 -path /api -target http://test.local/rest/api
```

## Enhanced Configuration

For more complex use-cases, the JSON configuration file allows for specializedJWT parameter setting (`jwtSecret`, `jwtIssuer`, `jwtExpire`), as well as interaction preferences like `allowCors` and `noCacheHeaders`. This also facilitates a basic authorization option with `allowPublicGet`.

## Token Generation

To get a token, make a POST request to the `token` endpoint with your username and password. See the example below:

```
POST http://localhost:8080/token
Content-Type: application/json

{
"username": "example",
"password": "password"
}

HTTP/1.1 200 OK

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJsb2NhbGhvc3QiLC..."
}
```

### User Management

Manage your own users by providing a text file to the `-users` option in the format

```
<username>    <bcrypted password>
```

### Downstream identity

In cases where the downstream endpoint needs to work with the authenticated user, TrantorHub sends a new `X-Trantor-Identity` header in the forwarded request with a hashed username as a value.

## Experimental Features

We also have early experimental flags like `-web` for interactive web authentication and proxy to a web application. With this, the `template/login.html` and `template/logout.html` can be utilized and the token will be returned and stored in the browser as an HttpOnly Cookie.