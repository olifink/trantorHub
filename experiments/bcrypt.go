//package experiments
//
//import (
//	"fmt"
//	"golang.org/x/crypto/bcrypt"
//)

package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is the cost factor
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	password := "password"
	hash, _ := hashPassword(password)

	fmt.Println("Password:", password)
	fmt.Println("Hash:", hash)

	// Check password against hash
	match := checkPasswordHash(password, hash)
	fmt.Println("Match:", match)
}
