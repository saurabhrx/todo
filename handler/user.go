package handler

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"myTodo/database/dbHelper"
	"myTodo/middleware"
	"myTodo/models"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var body models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
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
	userID, createErr := dbHelper.CreateUser(body.Name, body.Email, string(hashedPassword))
	if createErr != nil {
		http.Error(w, "failed to create new user", http.StatusInternalServerError)
		return
	}
	//_, sessErr := dbHelper.CreateSession(userID)
	//if sessErr != nil {
	//	fmt.Println(sessErr)
	//	return
	//}

	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "user created successfully",
		"user_id": userID,
	})
	if err != nil {
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "failed to Decode json", http.StatusBadRequest)
		return
	}
	userID, validateErr := dbHelper.ValidateUser(body.Email, body.Password)
	if validateErr != nil {
		err := json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid credentials",
		})
		if err != nil {
			return
		}
		return
	}

	accessToken, accessErr := middleware.GenerateAccessToken(userID)
	refreshToken, refreshErr := middleware.GenerateRefreshToken(userID)
	if accessErr != nil || refreshErr != nil {
		http.Error(w, "could not generate jwt token", http.StatusInternalServerError)
		return
	}

	_, sessErr := dbHelper.CreateSession(userID, refreshToken)
	if sessErr != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":       "user logged in successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	if EncodeErr != nil {
		return
	}

}
func Logout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID := middleware.UserContext(r)
	if err := dbHelper.Logout(userID, body.RefreshToken); err != nil {
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}

	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successfully",
	})
	if EncodeErr != nil {
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

func Refresh(w http.ResponseWriter, r *http.Request) {
	var body models.RefreshToken
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "failed to parse request body", http.StatusInternalServerError)
		return
	}
	if body.UserID == "" || body.Token == "" {
		http.Error(w, "missing user id or refresh token", http.StatusBadRequest)
		return
	}
	if dbHelper.ValidateSession(body.UserID, body.Token) {
		accessToken, err := middleware.GenerateAccessToken(body.UserID)
		if err != nil {
			http.Error(w, "could not generate jwt token", http.StatusInternalServerError)
			return
		}
		EncodeErr := json.NewEncoder(w).Encode(map[string]string{
			"message":      "new access token generated successfully",
			"access_token": accessToken,
		})
		if EncodeErr != nil {
			return
		}
	} else {
		err := dbHelper.Logout(body.UserID, body.Token)
		if err != nil {
			http.Error(w, "failed to delete the session", http.StatusInternalServerError)
		}
		http.Error(w, "session expired login again", http.StatusUnauthorized)
	}

}
