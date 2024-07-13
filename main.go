package main

import (
	"e-learn/internal/config"
	"e-learn/internal/database"
	"e-learn/internal/routes"
	"fmt"
)

func main() {
	config.LoadConfig()

	database.InitDB()
	//migrations.SeedDb(database.DB)

	fmt.Println("Hihhihihi.asddfsd. sdf.")

	r := routes.SetupRouter()
	r.Run(":8080")

	//fmt.Println("Hello World")
}
