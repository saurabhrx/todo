package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"myTodo/database"
	"myTodo/handler"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	if err := database.ConnectToDB(host, port, user, password, databaseName); err != nil {
		fmt.Println("Failed to connect database")
		return
	}

	fmt.Println("Database connected")
	r := mux.NewRouter()
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/logout", handler.Logout).Methods("POST")
	r.HandleFunc("/profile", handler.GetProfile).Methods("GET")
	r.HandleFunc("/delete-user", handler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/create-todo", handler.CreateTodo).Methods("POST")
	r.HandleFunc("/get-todos", handler.GetAllTodos).Methods("GET")
	r.HandleFunc("/update-todo", handler.UpdateTodo).Methods("PUT")
	r.HandleFunc("/delete-todo", handler.DeleteTodo).Methods("DELETE")

	SrvErr := http.ListenAndServe(":8080", r)
	if SrvErr != nil {
		return
	}

	DBCloseErr := database.CloseDBConnection()
	if DBCloseErr != nil {
		return
	}

}
