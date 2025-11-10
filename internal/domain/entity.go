package domain

import "time"

type User struct {
	ID           uint64     `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password_hash"`
	Email        string     `json:"email"`
	PhoneNumber  string     `json:"phone_number"`
	Role         string     `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	LastLogin    *time.Time `json:"last_login"`
}
