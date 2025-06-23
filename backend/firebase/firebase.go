package firebase

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	App           *firebase.App
	FirebaseAuth  *auth.Client
	firestoreClient *firestore.Client
)

// InitFirebase initializes the Firebase app and services
func InitFirebase() error {
	// Check if Firebase is enabled
	firebaseEnable := os.Getenv("FIREBASE_ENABLE")
	if firebaseEnable == "false" || firebaseEnable == "FALSE" {
		log.Println("🔥 Firebase is disabled via FIREBASE_ENABLE=false")
		log.Println("📁 Using local services instead")
		return nil
	}

	// Get the path to the service account key from environment variable or use default
	serviceAccountKeyPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")
	if serviceAccountKeyPath == "" {
		serviceAccountKeyPath = "serviceAccountKey.json"
	}

	// Check if service account key file exists
	if _, err := os.Stat(serviceAccountKeyPath); os.IsNotExist(err) {
		log.Printf("⚠️  Firebase service account key not found at %s", serviceAccountKeyPath)
		log.Println("📁 Firebase will be disabled, using local services")
		return nil
	}

	// Initialize Firebase app with service account
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	var err error
	App, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("⚠️  Error initializing Firebase app: %v", err)
		log.Println("📁 Firebase will be disabled, using local services")
		return nil
	}

	// Initialize Firebase Auth
	FirebaseAuth, err = App.Auth(context.Background())
	if err != nil {
		log.Printf("⚠️  Error initializing Firebase Auth: %v", err)
		log.Println("📁 Firebase Auth will be disabled, using local auth")
		FirebaseAuth = nil
	} else {
		log.Println("✅ Firebase Auth initialized successfully")
	}

	// Initialize Firestore
	firestoreClient, err = App.Firestore(context.Background())
	if err != nil {
		log.Printf("⚠️  Error initializing Firestore: %v", err)
		log.Println("📁 Firestore will be disabled, using local database")
		firestoreClient = nil
	} else {
		log.Println("✅ Firestore initialized successfully")
	}

	log.Println("🔥 Firebase initialization completed")
	return nil
}

// GetFirestoreClient returns the Firestore client
func GetFirestoreClient() *firestore.Client {
	return firestoreClient
}

// GetAuthClient returns the Firebase Auth client
func GetAuthClient() *auth.Client {
	return FirebaseAuth
}

// IsFirebaseEnabled returns whether Firebase is properly initialized
func IsFirebaseEnabled() bool {
	firebaseEnable := os.Getenv("FIREBASE_ENABLE")
	if firebaseEnable == "false" || firebaseEnable == "FALSE" {
		return false
	}
	return App != nil && (FirebaseAuth != nil || firestoreClient != nil)
}
