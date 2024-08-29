package main

import (
	"example.com/rest-api/db"
	_ "example.com/rest-api/docs"
	"example.com/rest-api/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Events Booking Rest API
// @version 1.0
// @description This is a backend API for an event booking application. It allows users to register for events, view events, and cancel registrations and authenticate users.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.Run(":8000") // localhost:8000
}
