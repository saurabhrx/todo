package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var Todo *sqlx.DB

func ConnectToDB(host, port, user, password, dbname string) error {
	plsqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	DB, err := sqlx.Open("postgres", plsqlInfo)
	if err != nil {
		return err
	}
	fmt.Println("abc....", plsqlInfo)
	err = DB.Ping()
	fmt.Println("m...")
	if err != nil {
		return err
	}
	fmt.Println("connection")
	Todo = DB
	return migrateUp(Todo)

}

func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://database/migrations", "postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	fmt.Println("Migration completed")
	return nil
}

func CloseDBConnection() error {
	return Todo.Close()
}
