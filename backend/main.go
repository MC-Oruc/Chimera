package main

import (
	"log"
	"os"

	"backend/firebase"
	"backend/routes"
	"backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Firebase (will be skipped if FIREBASE_ENABLE=false)
	if err := firebase.InitFirebase(); err != nil {
		log.Printf("Firebase initialization returned error: %v", err)
	}

	// Initialize services (Firebase or Local based on configuration)
	services.InitializeServices(firebase.GetFirestoreClient(), firebase.GetAuthClient())

	// Log the current mode
	log.Printf("🎯 Running in %s mode", services.GetServiceMode())

	// Set up Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add auth service to context (for both Firebase and Local auth)
	router.Use(func(c *gin.Context) {
		c.Set("authService", services.GetAuthService())
		
		// For backward compatibility, also set firebaseAuth if Firebase is enabled
		if services.IsFirebaseEnabled() && firebase.GetAuthClient() != nil {
			c.Set("firebaseAuth", firebase.GetAuthClient())
		}
		c.Next()
	})

	// Set up routes
	routes.SetupRouter(router)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
