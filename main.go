package main

import (
	"e-learn/internal/config"
	"e-learn/internal/database"
	"e-learn/internal/migrations"
	"e-learn/internal/routes"
)

func main() {
	config.LoadConfig()

	database.InitDB()
	migrations.SeedDb(database.DB)

	r := routes.SetupRouter()
	r.Run(":8080")

}
