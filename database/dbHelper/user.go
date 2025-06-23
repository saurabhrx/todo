package dbHelper

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"myTodo/database"
	"myTodo/models"
	"time"
)

func IsUserExists(email string) (bool, error) {
	SQL := `SELECT id FROM users WHERE email=$1 AND archived_at IS NULL `
	var id string
	err := database.Todo.Get(&id, SQL, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if errors.Is(err, sql.ErrNoRows) {
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

func CreateSession(userID string, refreshToken string) (string, error) {
	expireDate := time.Now().Add(48 * time.Hour)
	var sessionID string
	SQL := `INSERT INTO user_session(user_id,refresh_token,expired_at) VALUES($1,$2,$3) RETURNING id`
	if err := database.Todo.QueryRowx(SQL, userID, refreshToken, expireDate).Scan(&sessionID); err != nil {
		return "", err
	}
	return sessionID, nil
}
func ValidateUser(email, password string) (string, error) {
	SQL := `Select id , password from users where archived_at IS NULL and email=$1`
	var userId string
	var hashPassword string
	err := database.Todo.QueryRowx(SQL, email).Scan(&userId, &hashPassword)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	passwordErr := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if passwordErr != nil {
		return "", passwordErr
	}
	return userId, nil
}

func ValidateSession(userID string, refreshToken string) bool {
	SQL := `SELECT user_id from user_session WHERE refresh_token=$1 AND user_id = $2 AND expired_at > now() `
	var user string
	err := database.Todo.Get(&user, SQL, refreshToken, userID)
	if err != nil {
		return false
	}
	return true
}

func Logout(userID string, refreshToken string) error {
	SQL := `DELETE FROM user_session WHERE user_id=$1 AND refresh_token=$2`
	_, err := database.Todo.Exec(SQL, userID, refreshToken)
	return err

}

func GetProfile(userID string) (models.UserResponse, error) {
	SQL := `SELECT id,name , email FROM users WHERE id=$1 AND archived_at IS NULL`
	var userDetails models.UserResponse
	err := database.Todo.QueryRowx(SQL, userID).Scan(&userDetails.ID, &userDetails.Name, &userDetails.Email)
	if err != nil {
		return models.UserResponse{}, err
	}
	return userDetails, nil
}

func DeleteUser(userID string) error {
	SQL := `UPDATE users SET archived_at=now() WHERE id=$1`
	_, err := database.Todo.Exec(SQL, userID)
	return err
}

func DeleteSession(userID string, refreshToken string) error {
	SQL := `DELETE FROM user_session WHERE user_id=$1 AND refresh_token=$2`
	_, err := database.Todo.Exec(SQL, userID, refreshToken)
	return err
}
