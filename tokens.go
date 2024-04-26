package main

import (
	"github.com/golang-jwt/jwt/v5"
	"strings"
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

func isValidToken(tkn *jwt.Token, err error) bool {
	return err == nil && tkn.Valid
}

func parseTokenString(tknStr string) (*jwt.Token, error) {
	tknStr = strings.TrimPrefix(tknStr, "Bearer ")
	claims := &jwt.RegisteredClaims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		secret := []byte(config.JwtSecret)
		return secret, nil
	})
	return tkn, err
}
