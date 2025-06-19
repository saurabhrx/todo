package routes

import (
	"github.com/gorilla/mux"
	"myTodo/handler"
)

func SetupTodoRoutes() *mux.Router {
	srv := mux.NewRouter()

	// user routes

	srv.HandleFunc("/register", handler.Register).Methods("POST")
	srv.HandleFunc("/login", handler.Login).Methods("POST")
	srv.HandleFunc("/logout", handler.Logout).Methods("POST")
	srv.HandleFunc("/profile", handler.GetProfile).Methods("GET")
	srv.HandleFunc("/delete-user", handler.DeleteUser).Methods("DELETE")

	// todo routes

	srv.HandleFunc("/create-todo", handler.CreateTodo).Methods("POST")
	srv.HandleFunc("/get-todos", handler.GetAllTodos).Methods("GET")
	srv.HandleFunc("/update-todo", handler.UpdateTodo).Methods("PUT")
	srv.HandleFunc("/delete-todo", handler.DeleteTodo).Methods("DELETE")

	return srv
}
