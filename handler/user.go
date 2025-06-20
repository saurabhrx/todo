package handler

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"myTodo/database/dbHelper"
	"myTodo/middleware"
	"myTodo/models"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var body models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "failed to parse request body", http.StatusInternalServerError)
		return
	}
	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		http.Error(w, "error while creating User", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	userID, CreateErr := dbHelper.CreateUser(body.Name, body.Email, string(hashedPassword))
	if CreateErr != nil {
		http.Error(w, "failed to create new user", http.StatusInternalServerError)
		return
	}
	_, sessErr := dbHelper.CreateSession(userID)
	if sessErr != nil {
		fmt.Println(sessErr)
		return
	}

	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "user created successfully",
	})
	if err != nil {
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "failed to Decode json", http.StatusInternalServerError)
		return
	}
	userID, validateErr := dbHelper.ValidateUser(body.Email, body.Password)
	if validateErr != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid credentials",
		})
		return
	}
	sessionID, sessErr := dbHelper.CreateSession(userID)
	if sessErr != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(map[string]string{
		"message":       "user logged in successfully",
		"session_token": sessionID,
	})
	if err != nil {
		return
	}

}
func Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("session-id")
	LogoutErr := dbHelper.Logout(sessionID)
	if LogoutErr != nil {
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successfully",
	})
	if err != nil {
		return
	}
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserContext(r)
	userDetails, GetProfileErr := dbHelper.GetProfile(userID)
	if GetProfileErr != nil {
		http.Error(w, "failed to get user profile", http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(map[string]string{
		"id":    userDetails.ID,
		"name":  userDetails.Name,
		"email": userDetails.Email,
	})
	if err != nil {
		return
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	userID := middleware.UserContext(r)

	DeleteErr := dbHelper.DeleteUser(userID)
	if DeleteErr != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "user deleted successfully",
	})
	if err != nil {
		return
	}
}
