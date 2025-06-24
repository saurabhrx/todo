package dbHelper

import (
	"myTodo/database"
	"myTodo/models"
)

func CreateTodo(userID, name, description string, status bool) error {
	SQL := `INSERT INTO todo(user_id, name, description, status) VALUES ($1,$2,$3,$4)`
	_, err := database.Todo.Exec(SQL, userID, name, description, status)
	return err

}

func GetTodos(userID string) ([]models.TodoResponse, error) {
	SQL := `SELECT id , user_id , name , description , status , created_at FROM todo WHERE user_id=$1 AND archived_at is NULL`
	var todos []models.TodoResponse
	err := database.Todo.Select(&todos, SQL, userID)
	return todos, err
}

func UpdateTodo(todoID, userID, name, description string, status bool) (bool, error) {
	SQL := `UPDATE todo SET name=$1 ,description = $2 , status=$3 WHERE id = $4 AND user_id=$5 AND archived_at IS NULL `
	result, err := database.Todo.Exec(SQL, name, description, status, todoID, userID)
	if err != nil {
		return false, err
	}
	res, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return res > 0, err
}
func DeleteTodo(todoID string, userID string) (bool, error) {
	SQL := `UPDATE todo SET archived_at=now() WHERE id=$1 AND user_id=$2`
	result, err := database.Todo.Exec(SQL, todoID, userID)
	if err != nil {
		return false, err
	}
	res, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return res > 0, err
}
