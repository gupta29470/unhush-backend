package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"unhush-backend/routes"
)

func main() {
	envError := godotenv.Load(".env")
	if envError != nil {
		log.Fatal("Something went wrong while loading env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	engine := gin.New()
	engine.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // or set specific origins
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}

	engine.Use(cors.New(config))

	routes.AppRoutes(engine)

	runError := engine.Run(":" + port)
	if runError != nil {
		log.Fatal("Something went wrong while running server")
	}
}
