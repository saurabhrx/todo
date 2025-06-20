package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"myTodo/database/dbHelper"
	"myTodo/middleware"
	"myTodo/models"
	"net/http"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {

	userID := middleware.UserContext(r)
	var body models.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "failed to parse request body", http.StatusInternalServerError)
		return
	}
	CreateErr := dbHelper.CreateTodo(userID, body.Name, body.Description, false)
	if CreateErr != nil {
		http.Error(w, "failed to create todo", http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "todo created successfully",
	})
	if err != nil {
		return
	}
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserContext(r)
	todos, GetTodoErr := dbHelper.GetTodos(userID)
	if GetTodoErr != nil {
		http.Error(w, "failed to fetch todos", http.StatusInternalServerError)
		return
	}

	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		return
	}
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserContext(r)
	var body models.TodoUpdate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logrus.Panic("failed to parse request body")
		return
	}

	UpdateErr := dbHelper.UpdateTodo(body.ID, userID, body.Name, body.Description, body.Status)
	if UpdateErr != nil {
		http.Error(w, "failed to Update", http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "todo updated successfully",
	})
	if err != nil {
		return
	}

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {

	var body models.TodoDelete
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logrus.Panic("failed to parse request body")
		return
	}

	DeleteErr := dbHelper.DeleteTodo(body.ID)
	if DeleteErr != nil {
		http.Error(w, "failed to delete todo", http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": "todo deleted successfully",
	})
	if err != nil {
		return
	}
}
