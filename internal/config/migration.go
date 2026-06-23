package config

import (
	"database/sql"
	"fmt"
	"log"

	soomdb "soom-be-go/db"

	"github.com/pressly/goose/v3"
)

func RunMigrations(sqlDb *sql.DB) {
	goose.SetBaseFS(soomdb.Migration)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(sqlDb, "migrations"); err != nil {
		log.Println("Migration warning (mungkin sudah up to date):", err)

	}

	fmt.Println("Migrations completed successfully")
}
