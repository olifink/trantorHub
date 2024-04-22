package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

type User struct {
	//ID       int
	Username string
	Password string // This should be a hashed password
}

var Users []User

// GetUserByUsername fetches a user by username from the database
func GetUserByUsername(username string) (*User, error) {
	// For now, let's assume we get some user or nil if not found
	for _, user := range Users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, nil
}

func readUsers() {
	// default data is example/password
	csvUsers := `example "$2a$14$HNOQGnDpfyF/95TT6VToEuyS4NCYKXH1pVlcq9fx9JaC/zBW.cn0i"`

	// use file if it was given
	var data io.Reader
	if config.UserFile == "" {
		data = strings.NewReader(csvUsers)
	} else {
		file, err := os.Open(config.UserFile)
		if err != nil {
			log.Panicln("Error opening user file:", err)
		}
		defer file.Close()
		data = file
	}

	// parse file as 2 columns. space separated
	r := csv.NewReader(data)
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = 2
	r.Comma = ' '

	// read all entries into user array
	records, err := r.ReadAll()
	if err != nil {
		log.Panicln("error reading users", err)
	}

	Users = make([]User, len(records))
	for i, record := range records {
		Users[i] = User{
			Username: record[0],
			Password: record[1],
		}
	}
	log.Println("Users:", Users)
}
