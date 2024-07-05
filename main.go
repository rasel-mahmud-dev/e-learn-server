package main

import (
	"e-learn/internal/config"
	"e-learn/internal/database"
	"e-learn/internal/routes"
)

var db = make(map[string]string)

func main() {
	config.LoadConfig()

	database.InitDB()
	database.SeedDb()

	r := routes.SetupRouter()
	r.Run(":8080")

}
