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
	userId, saveErr := dbHelper.CreateUser(body.Name, body.Email, string(hashedPassword))
	if saveErr != nil {
		fmt.Println(saveErr)
	}
	sessErr := dbHelper.CreateSession(userId)
	if sessErr != nil {
		fmt.Println(sessErr)
	}

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

	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "User Login",
	})
	if err != nil {
		return
	}
}
