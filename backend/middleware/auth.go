package middleware

import (
	"backend/interfaces"
	"context"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the token using the appropriate auth service
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("No Authorization header found for path: %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == authHeader {
			log.Printf("Authorization header does not have Bearer prefix for path: %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in the format 'Bearer {token}'"})
			c.Abort()
			return
		}

		// Get the auth service from context
		authService, exists := c.Get("authService")
		if exists {
			// Use the interface-based auth service
			if authSvc, ok := authService.(interfaces.AuthService); ok {
				user, err := authSvc.VerifyToken(context.Background(), idToken)
				if err != nil {
					log.Printf("Error verifying token with auth service for path %s: %v", c.Request.URL.Path, err)
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
					c.Abort()
					return
				}

				// Set the user ID in the context
				c.Set("userId", user.UID)
				log.Printf("User authenticated via auth service: %s for path: %s", user.UID, c.Request.URL.Path)
				c.Next()
				return
			}
		}

		// Fallback to Firebase Auth for backward compatibility
		firebaseAuth, exists := c.Get("firebaseAuth")
		if !exists {
			log.Printf("No auth service found in context for path: %s", c.Request.URL.Path)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication service not initialized"})
			c.Abort()
			return
		}

		// Verify the ID token with Firebase
		token, err := firebaseAuth.(*auth.Client).VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Printf("Error verifying Firebase ID token for path %s: %v", c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("userId", token.UID)
		log.Printf("User authenticated via Firebase: %s for path: %s", token.UID, c.Request.URL.Path)

		c.Next()
	}
}

// OptionalAuthMiddleware tries to authenticate the user but allows the request to proceed even if authentication fails
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("No Authorization header found for path: %s (optional auth)", c.Request.URL.Path)
			c.Next()
			return
		}

		// Check if the header has the Bearer prefix
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == authHeader {
			log.Printf("Authorization header does not have Bearer prefix for path: %s (optional auth)", c.Request.URL.Path)
			c.Next()
			return
		}

		// Try to get the auth service from context
		authService, exists := c.Get("authService")
		if exists {
			// Use the interface-based auth service
			if authSvc, ok := authService.(interfaces.AuthService); ok {
				user, err := authSvc.VerifyToken(context.Background(), idToken)
				if err != nil {
					log.Printf("Error verifying token with auth service for path %s (optional auth): %v", c.Request.URL.Path, err)
					c.Next()
					return
				}

				// Set the user ID in the context
				c.Set("userId", user.UID)
				log.Printf("User authenticated via auth service: %s for path: %s (optional auth)", user.UID, c.Request.URL.Path)
				c.Next()
				return
			}
		}

		// Fallback to Firebase Auth for backward compatibility
		firebaseAuth, exists := c.Get("firebaseAuth")
		if !exists {
			log.Printf("No auth service found in context for path: %s (optional auth)", c.Request.URL.Path)
			c.Next()
			return
		}

		// Verify the ID token with Firebase
		token, err := firebaseAuth.(*auth.Client).VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Printf("Error verifying Firebase ID token for path %s (optional auth): %v", c.Request.URL.Path, err)
			c.Next()
			return
		}

		// Set the user ID in the context
		c.Set("userId", token.UID)
		log.Printf("User authenticated via Firebase: %s for path: %s (optional auth)", token.UID, c.Request.URL.Path)

		c.Next()
	}
}
