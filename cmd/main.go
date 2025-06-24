package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"myTodo/database"
	"myTodo/routes"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
		return
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	if err := database.ConnectToDB(host, port, user, password, databaseName); err != nil {
		logrus.Panic("failed to connect database")
		return
	}

	fmt.Println("database connected")

	srv := routes.SetupTodoRoutes()

	SrvErr := http.ListenAndServe(":8080", srv)

	if SrvErr != nil {
		logrus.Panic("failed to connect to server")
		return
	}

	DBCloseErr := database.CloseDBConnection()
	if DBCloseErr != nil {
		logrus.Panic("failed to close database")
		return
	}

}
