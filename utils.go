package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// Create a SHA256 hash from a string
func createHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Anonymize a part of a sensitive string
func anonymize(s string) string {
	if len(s) > 4 {
		return s[:2] + "****" + s[len(s)-2:]
	} else {
		return "****"
	}
}

// identity of an authenticated user, including the user object and the authentication token.
type userIdentity struct {
	user  *User
	token string
}

// identityForUser verifies the given username and password,
// checks if the user is known, and compares the password
// with the hashed password in the database. It returns a userIdentity
// object that includes the user and an authentication token,
// or an error if any of the checks fail.
func identityForUser(username string, password string) (*userIdentity, error) {
	// Check if the username and password are empty
	if username == "" || password == "" {
		return nil, errors.New("username or password is empty")
	}

	// Check if we know the user
	user := GetUserByUsername(username)
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check that the password matches the hashed password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("wrong password")
	}

	tokenString, err := generateNewToken(username)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return &userIdentity{
		user:  user,
		token: tokenString,
	}, nil
}
