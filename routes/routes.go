package routes

import (
	"github.com/gorilla/mux"
	"myTodo/handler"
	"myTodo/middleware"
)

func SetupTodoRoutes() *mux.Router {
	srv := mux.NewRouter()

	// user routes
	srv.HandleFunc("/refresh", handler.Refresh).Methods("GET")
	srv.HandleFunc("/register", handler.Register).Methods("POST")
	srv.HandleFunc("/login", handler.Login).Methods("POST")

	protected := srv.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/logout", handler.Logout).Methods("POST")
	protected.HandleFunc("/profile", handler.GetProfile).Methods("GET")
	protected.HandleFunc("/delete-user", handler.DeleteUser).Methods("DELETE")

	// todo routes

	protected.HandleFunc("/create-todo", handler.CreateTodo).Methods("POST")
	protected.HandleFunc("/get-todos", handler.GetAllTodos).Methods("GET")
	protected.HandleFunc("/update-todo", handler.UpdateTodo).Methods("PUT")
	protected.HandleFunc("/delete-todo", handler.DeleteTodo).Methods("DELETE")

	return srv
}
