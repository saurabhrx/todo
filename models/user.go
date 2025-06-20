package models

type UserResponse struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}
type UserRequest struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}
