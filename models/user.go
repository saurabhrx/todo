package models

import "time"

type User struct {
	ID         string     `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Email      string     `json:"email" db:"email"`
	Password   string     `json:"-" db:"password"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ArchivedAt *time.Time `json:"archived_at" db:"archived_at"`
}

type Todo struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ArchivedAt  *time.Time `json:"archived_at" db:"archived_at"`
}

type Session struct {
	ID         string `json:"id" db:"id"`
	UserID     string `json:"user_id" db:"user_id"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	ArchivedAt string `json:"archived_at" db:"archived_at"`
}
