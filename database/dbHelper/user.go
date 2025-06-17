package dbHelper

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"myTodo/database"
	"time"
)

func IsUserExists(email string) (bool, error) {
	SQL := `SELECT id FROM users WHERE email=$1 AND archived_at IS NULL `
	var id string
	err := database.Todo.Get(&id, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func CreateUser(name string, email string, password string) (string, error) {
	SQL := `INSERT INTO users(name , email , password) VALUES ($1,$2,$3) RETURNING id`
	var userID string
	if err := database.Todo.QueryRowx(SQL, name, email, password).Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func CreateSession(userID string) error {
	expireDate := time.Now().Add(7 * 24 * time.Hour)
	SQL := `INSERT INTO user_session(user_id,expired_at) VALUES($1,$2)`
	_, err := database.Todo.Exec(SQL, userID, expireDate)
	return err
}
func ValidateUser(email, password string) (string, error) {
	SQL := `Select id , password from users where archived_at IS NULL and email=$1`
	var userId string
	var hashPassword string
	err := database.Todo.QueryRowx(SQL, email).Scan(&userId, &hashPassword)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		return "", nil
	}
	passwordErr := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if passwordErr != nil {
		return "", passwordErr
	}
	return userId, nil
}
