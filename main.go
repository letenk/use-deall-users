package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/letenk/use_deal_user/config"
	"github.com/letenk/use_deal_user/router"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Failed load .env: %v", err)
	}

	fmt.Println("Server is starting...")
	db := config.SetupDB()
	router := router.SetupRouter(db)
	router.Run(":8080")
}
