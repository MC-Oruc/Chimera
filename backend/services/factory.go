package services

import (
	"backend/interfaces"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
)

// ServiceFactory creates and manages all services based on configuration
type ServiceFactory struct {
	firebaseEnabled bool
	
	// Services
	DatabaseService interfaces.DatabaseService
	StorageService  interfaces.StorageService
	AuthService     interfaces.AuthService
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(firestoreClient *firestore.Client, authClient *auth.Client) *ServiceFactory {
	factory := &ServiceFactory{}
	
	// Check if Firebase is enabled
	firebaseEnable := os.Getenv("FIREBASE_ENABLE")
	factory.firebaseEnabled = firebaseEnable != "false" && firebaseEnable != "FALSE"
	
	if factory.firebaseEnabled {
		log.Println("🔥 Firebase mode: ENABLED")
		factory.initFirebaseServices(firestoreClient, authClient)
	} else {
		log.Println("📁 Local mode: ENABLED (Firebase disabled)")
		factory.initLocalServices()
	}
	
	return factory
}

func (f *ServiceFactory) initFirebaseServices(firestoreClient *firestore.Client, authClient *auth.Client) {
	if firestoreClient != nil {
		f.DatabaseService = NewFirebaseDatabase(firestoreClient)
		log.Println("✅ Firebase Database service initialized")
	} else {
		log.Println("⚠️  Warning: Firestore client is nil, falling back to local database")
		f.DatabaseService = NewLocalDatabase()
	}
	
	if authClient != nil {
		f.AuthService = NewFirebaseAuth(authClient)
		log.Println("✅ Firebase Auth service initialized")
	} else {
		log.Println("⚠️  Warning: Firebase Auth client is nil, falling back to local auth")
		f.AuthService = NewLocalAuth()
	}
	
	// Firebase Storage
	f.StorageService = NewFirebaseService()
	log.Println("✅ Firebase Storage service initialized")
}

func (f *ServiceFactory) initLocalServices() {
	f.DatabaseService = NewLocalDatabase()
	f.StorageService = NewLocalStorage()
	f.AuthService = NewLocalAuth()
	
	log.Println("✅ Local Database service initialized")
	log.Println("✅ Local Storage service initialized")
	log.Println("✅ Local Auth service initialized")
}

// IsFirebaseEnabled returns whether Firebase is enabled
func (f *ServiceFactory) IsFirebaseEnabled() bool {
	return f.firebaseEnabled
}

// GetDatabaseService returns the database service
func (f *ServiceFactory) GetDatabaseService() interfaces.DatabaseService {
	return f.DatabaseService
}

// GetStorageService returns the storage service
func (f *ServiceFactory) GetStorageService() interfaces.StorageService {
	return f.StorageService
}

// GetAuthService returns the auth service
func (f *ServiceFactory) GetAuthService() interfaces.AuthService {
	return f.AuthService
}

// GetMode returns a string describing the current mode
func (f *ServiceFactory) GetMode() string {
	if f.firebaseEnabled {
		return "Firebase"
	}
	return "Local"
} 