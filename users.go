package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type User struct {
	//ID       int
	Username string
	Password string // This is a bcyrpt hashed password
	Identity string // Hashed username for downstream systems
}

var Users []User

// GetUserByUsername fetches a user by username from the database
func GetUserByUsername(username string) *User {
	// For now, let's assume we get some user or nil if not found
	for _, user := range Users {
		if user.Username == username {
			return &user
		}
	}
	return nil
}

func readUsers() {
	// default data is example/password
	csvUsers := `example:$2a$14$HNOQGnDpfyF/95TT6VToEuyS4NCYKXH1pVlcq9fx9JaC/zBW.cn0i`

	// use file if it was given
	var data io.Reader
	if config.UserFile == "" {
		data = strings.NewReader(csvUsers)
	} else {
		file, err := os.Open(config.UserFile)
		if err != nil {
			log.Fatalln("Error opening user file:", err)
		}
		defer file.Close()
		data = file
	}

	// parse file as 2 columns. colon separated
	r := csv.NewReader(data)
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = 2
	r.Comma = ':'

	// read all entries into user array
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln("error reading users", err)
	}

	Users = make([]User, len(records))
	for i, record := range records {
		log.Println(fmt.Sprintf("User %d: %s %s", i, record[0], anonymize(record[1])))
		Users[i] = User{
			Username: record[0],
			Password: record[1],
			Identity: createHash(record[0]),
		}
	}
}
