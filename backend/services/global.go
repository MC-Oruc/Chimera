package services

import (
	"backend/interfaces"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
)

var (
	serviceFactory *ServiceFactory
	once           sync.Once
)

// InitializeServices initializes the global service factory
func InitializeServices(firestoreClient *firestore.Client, authClient *auth.Client) {
	once.Do(func() {
		serviceFactory = NewServiceFactory(firestoreClient, authClient)
		log.Println("🚀 Global services initialized successfully")
	})
}

// GetDatabaseService returns the global database service
func GetDatabaseService() interfaces.DatabaseService {
	if serviceFactory == nil {
		log.Fatal("❌ Services not initialized. Call InitializeServices first.")
	}
	return serviceFactory.GetDatabaseService()
}

// GetStorageService returns the global storage service
func GetStorageService() interfaces.StorageService {
	if serviceFactory == nil {
		log.Fatal("❌ Services not initialized. Call InitializeServices first.")
	}
	return serviceFactory.GetStorageService()
}

// GetAuthService returns the global auth service
func GetAuthService() interfaces.AuthService {
	if serviceFactory == nil {
		log.Fatal("❌ Services not initialized. Call InitializeServices first.")
	}
	return serviceFactory.GetAuthService()
}

// IsFirebaseEnabled returns whether Firebase is enabled
func IsFirebaseEnabled() bool {
	if serviceFactory == nil {
		log.Fatal("❌ Services not initialized. Call InitializeServices first.")
	}
	return serviceFactory.IsFirebaseEnabled()
}

// GetServiceMode returns the current service mode
func GetServiceMode() string {
	if serviceFactory == nil {
		log.Fatal("❌ Services not initialized. Call InitializeServices first.")
	}
	return serviceFactory.GetMode()
} 