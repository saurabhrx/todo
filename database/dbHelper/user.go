package dbHelper

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"myTodo/database"
	"myTodo/models"
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

func ValidateSession(sessionID string) (string, error) {
	SQL := `SELECT user_id from user_session WHERE id=$1 AND expired_at > now() AND archived_at IS NULL `
	var userID string
	err := database.Todo.Get(&userID, SQL, sessionID)
	return userID, err
}

func Logout(sessionID string) error {
	SQL := `UPDATE user_session SET archived_at=now(),expired_at=now() WHERE id=$1`
	_, err := database.Todo.Exec(SQL, sessionID)
	return err

}

func GetProfile(userID string) (models.User, error) {
	SQL := `SELECT * FROM users WHERE id=$1 AND archived_at IS NULL`
	var userDetails models.User
	err := database.Todo.Get(&userDetails, SQL, userID)
	return userDetails, err
}

func CreateTodo(userID, name, description string, status bool) error {
	SQL := `INSERT INTO todo(user_id, name, description, status) VALUES ($1,$2,$3,$4)`
	_, err := database.Todo.Exec(SQL, userID, name, description, status)
	return err

}

func GetTodos(userID string) ([]models.Todo, error) {
	SQL := `SELECT * FROM todo WHERE user_id=$1 AND archived_at is NULL`
	var todos []models.Todo
	err := database.Todo.Select(&todos, SQL, userID)
	return todos, err
}

func UpdateTodo(todoID, userID, name, description string, status bool) error {
	SQL := `UPDATE todo SET name=$1 ,description = $2 , status=$3 WHERE id = $4 AND user_id=$5 AND archived_at IS NULL `
	_, err := database.Todo.Exec(SQL, name, description, status, todoID, userID)
	return err
}
func DeleteTodo(todoID string) error {
	SQL := `UPDATE todo SET archived_at=now() WHERE id=$1`
	_, err := database.Todo.Exec(SQL, todoID)
	return err
}

func DeleteUser(userID string) error {
	SQL := `UPDATE users SET archived_at=now() WHERE id=$1`
	_, err := database.Todo.Exec(SQL, userID)
	return err
}
