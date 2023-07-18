package main

import (
	"log"
	"os"

	"dataimpact-backend/database"
	"dataimpact-backend/handlers"
	"dataimpact-backend/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger()

	// I set this env Variable in docker compose file.
	// If you are not running with docker, app will
	// connect trough localhost
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		mongodbURI = "mongodb://localhost:27017"
	}
	// Connect to MongoDB
	db, err := database.Connect(mongodbURI, "userdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect()

	// Create a new Gin router
	r := gin.Default()

	// Initialize user handler
	userHandler := handlers.NewUserHandler(db)

	// Define API routes
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.POST("/login", userHandler.Login)

			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUserByID)

			users.DELETE("/:id", userHandler.DeleteUser)
			users.PUT("/:id", userHandler.UpdateUser)
		}
	}

	// Run the server
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
