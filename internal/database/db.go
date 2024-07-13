package database

import (
	"database/sql"
	"e-learn/internal/config"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := config.Cfg.DATABASE_URI
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connection established")
	DB = db
}
