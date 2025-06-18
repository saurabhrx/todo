package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"myTodo/database/dbHelper"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logrus.Panic("Failed to parse request body")
		return
	}
	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		fmt.Println("Error while creating user", existsErr)
		return
	}
	if exists {
		fmt.Println("User Already exists")
		return
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println(hashErr)
		return
	}
	_, saveErr := dbHelper.CreateUser(body.Name, body.Email, string(hashedPassword))
	if saveErr != nil {
		fmt.Println(saveErr)
	}
	//sessErr := dbHelper.CreateSession(userId)
	//if sessErr != nil {
	//	fmt.Println(sessErr)
	//}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logrus.Panic("Failed to parse request body")
		return
	}
	userID, validateErr := dbHelper.ValidateUser(body.Email, body.Password)
	if validateErr != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid Credentials",
		})
	}
	sessErr := dbHelper.CreateSession(userID)
	if sessErr != nil {
		fmt.Println(sessErr)
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User Login",
	})

}
func Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("session-id")
	if sessionID == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return

	}

	_, err := dbHelper.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Invalid or Expired Session", http.StatusUnauthorized)
		return
	}

	LogoutErr := dbHelper.Logout(sessionID)
	if LogoutErr != nil {
		http.Error(w, "Logout Failed", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successfully",
	})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("session-id")
	if sessionID == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return

	}

	userID, err := dbHelper.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Invalid or Expired Session", http.StatusUnauthorized)
		return
	}
	userDetails, GetProfileErr := dbHelper.GetProfile(userID)
	if GetProfileErr != nil {
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(userDetails)

}
