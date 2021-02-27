package model

import "time"

// Session
type Session struct {
	Username string    `json:"username"`
	Expires  time.Time `json:"expires"`
}

// Credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User
type User struct {
	Password string
}
