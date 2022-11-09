package main

import (
	"fmt"

	"github.com/letenk/use_deal_user/config"
	"github.com/letenk/use_deal_user/router"
)

func main() {
	fmt.Println("Server is starting...")
	db := config.SetupDB()
	router := router.SetupRouter(db)
	router.Run(":8080")
}
