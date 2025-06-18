package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"myTodo/database/dbHelper"
	"net/http"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {

	sessionID := r.Header.Get("session-id")
	if sessionID == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
	}

	userID, err := dbHelper.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Invalid or Expired Session", http.StatusUnauthorized)
	}

	body := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      bool   `json:"status"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logrus.Panic("Failed to parse request body")
		return
	}
	CreateErr := dbHelper.CreateTodo(userID, body.Name, body.Description, body.Status)
	if CreateErr != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo created successfully",
	})
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("session-id")
	if sessionID == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
	}

	userID, err := dbHelper.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Invalid or Expired Session", http.StatusUnauthorized)

	}
	todos, GetTodoErr := dbHelper.GetTodos(userID)
	if GetTodoErr != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(todos)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("session-id")
	if sessionID == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
	}

	userID, err := dbHelper.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Invalid or Expired Session", http.StatusUnauthorized)

	}
	body := struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      bool   `json:"status"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logrus.Panic("Failed to parse request body")
		return
	}

	UpdateErr := dbHelper.UpdateTodo(body.ID, userID, body.Name, body.Description, body.Status)
	if UpdateErr != nil {
		http.Error(w, "Failed to Update", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo Updated Successfully",
	})

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("session-id")
	if sessionID == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
	}

	_, err := dbHelper.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Invalid or Expired Session", http.StatusUnauthorized)

	}
	body := struct {
		ID string `json:"id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logrus.Panic("Failed to parse request body")
		return
	}

	DeleteErr := dbHelper.DeleteTodo(body.ID)
	if DeleteErr != nil {
		http.Error(w, "Failed to Delete todo", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo deleted successfully",
	})
}
