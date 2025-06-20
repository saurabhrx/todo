package models

import "time"

type TodoResponse struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Status      bool      `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type TodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TodoUpdate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

type TodoDelete struct {
	ID string `json:"id"`
}
