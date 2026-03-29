package model

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string // bcrypt hash
	Role      string // "user" or "admin"
	CreatedAt time.Time
	UpdatedAt time.Time
}