package main

type User struct {
	ID       int
	Username string
	Password string // This should be a hashed password
}

// GetUserByUsername fetches a user by username from the database
func GetUserByUsername(username string) (*User, error) {
	// TODO Database fetching logic here
	// For now, let's assume we get some user or nil if not found
	return &User{ID: 1, Username: "example", Password: "$2a$14$HNOQGnDpfyF/95TT6VToEuyS4NCYKXH1pVlcq9fx9JaC/zBW.cn0i"}, nil // bcrypt hash for "password"
}
