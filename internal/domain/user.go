package domain

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // hashed password
	CreatedAt time.Time `json:"created_at"`
}
