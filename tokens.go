package main

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func generateNewToken(username string) (string, error) {
	// Create the JWT claims, which includes the username and expiry time
	var expirationTime *jwt.NumericDate

	if config.expireDuration > 0 {
		expirationTime = jwt.NewNumericDate(time.Now().Add(config.expireDuration))
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: expirationTime,
		Issuer:    config.JwtIssuer,
		Subject:   username,
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	secret := []byte(config.JwtSecret)
	return token.SignedString(secret)
}

func isValidToken(tknStr string) bool {
	tkn, err := parseTokenString(tknStr)
	return err == nil && tkn.Valid
}
