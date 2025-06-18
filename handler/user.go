package handler

import (
	"encoding/json"
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
		http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
		return
	}
	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		http.Error(w, "Error while creating User", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	_, CreateErr := dbHelper.CreateUser(body.Name, body.Email, string(hashedPassword))
	if CreateErr != nil {
		http.Error(w, "Failed to create new user", http.StatusInternalServerError)
		return
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
		http.Error(w, "Failed to Decode json", http.StatusInternalServerError)
		return
	}
	userID, validateErr := dbHelper.ValidateUser(body.Email, body.Password)
	if validateErr != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid Credentials",
		})
	}
	sessionID, sessErr := dbHelper.CreateSession(userID)
	if sessErr != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "User Login",
		"session_token": sessionID,
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
		return
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
		return
	}
	json.NewEncoder(w).Encode(userDetails)

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	body := struct {
		ID string `json:"id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
		return
	}

	DeleteErr := dbHelper.DeleteUser(body.ID)
	if DeleteErr != nil {
		http.Error(w, "Failed to Delete User", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}
