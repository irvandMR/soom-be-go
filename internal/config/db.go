package config

import (
	"database/sql"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, *sql.DB) {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: false})
	if err != nil {
		panic("Failed to connect to the database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ DB Connected")
	return db, sqlDB
}
