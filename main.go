package main

import (
	"e-learn/internal/config"
	"e-learn/internal/database"
	"e-learn/internal/routes"
)

func main() {
	config.LoadConfig()
	database.InitDB()
	r := routes.SetupRouter()
	r.Run(":8080")
}
