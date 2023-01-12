package main

import (
	"fmt"

	"github.com/Similadayo/backend/db"
	"github.com/Similadayo/backend/router"
	"github.com/gin-gonic/gin"
)

func main() {
	//connect to database
	db.ConnectToDb()
	defer db.DB.Close()

	//
	r := gin.Default()
	router.SetupRouter(r)
	fmt.Println("Server is running...")
	r.Run(":8080")
}
