package main

import (
	"crypto/sha256"
	"fmt"
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
