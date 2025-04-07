package main

import (
	"github.com/PedroNetto404/marmitech-backend/internal/config"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.LoadEnvs()

	db, err := database.New()
	if err != nil {
		panic(err)
	}

	// TODO: fix this. Current error is: marmitech_user:marmitech_strong_password1235r43@!$@tcp(localhost:3306)/marmitech_db?charset=utf8mb4&parseTime=True&loc=Local
	if err := db.Migrate(); err != nil {
		panic(err)
	}
	if err := db.Close(); err != nil {
		panic(err)
	}
}
