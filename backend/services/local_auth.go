package services

import (
	"backend/interfaces"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type LocalAuth struct {
	dataDir string
	tokens  map[string]*TokenInfo
	mu      sync.RWMutex
}

type TokenInfo struct {
	UserID    string    `json:"userId"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type LocalUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"` // In real implementation, this should be hashed
	Name     string `json:"name"`
}

// NewLocalAuth creates a new local authentication service
func NewLocalAuth() interfaces.AuthService {
	dataDir := "local_auth"
	
	// Ensure auth directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Warning: Failed to create auth directory %s: %v", dataDir, err)
	}
	
	auth := &LocalAuth{
		dataDir: dataDir,
		tokens:  make(map[string]*TokenInfo),
	}
	
	// Load existing tokens
	auth.loadTokens()
	
	// Start cleanup routine for expired tokens
	go auth.cleanupExpiredTokens()
	
	log.Println("Local authentication initialized with directory:", dataDir)
	return auth
}

func (a *LocalAuth) VerifyToken(ctx context.Context, token string) (*interfaces.AuthUser, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	// Handle different token formats
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)
	
	// Check for local token
	if tokenInfo, exists := a.tokens[token]; exists {
		if time.Now().Before(tokenInfo.ExpiresAt) {
			return &interfaces.AuthUser{
				UID:   tokenInfo.UserID,
				Email: tokenInfo.Email,
			}, nil
		} else {
			// Token expired, remove it
			delete(a.tokens, token)
			a.saveTokens()
			return nil, fmt.Errorf("token expired")
		}
	}
	
	// For development purposes, accept simple tokens like "user1", "user2", etc.
	if strings.HasPrefix(token, "user") {
		userID := token
		email := fmt.Sprintf("%s@localhost.com", userID)
		
		// Create a new token for this user
		newToken := a.generateToken()
		a.tokens[newToken] = &TokenInfo{
			UserID:    userID,
			Email:     email,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		a.saveTokens()
		
		return &interfaces.AuthUser{
			UID:   userID,
			Email: email,
		}, nil
	}
	
	return nil, fmt.Errorf("invalid token")
}

func (a *LocalAuth) ValidateToken(token string) bool {
	_, err := a.VerifyToken(context.Background(), token)
	return err == nil
}

// CreateUser creates a new user (for development/testing)
func (a *LocalAuth) CreateUser(email, password, name string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	userID := a.generateUserID()
	user := LocalUser{
		ID:       userID,
		Email:    email,
		Password: password, // In production, this should be hashed
		Name:     name,
	}
	
	filePath := filepath.Join(a.dataDir, "users", fmt.Sprintf("%s.json", userID))
	
	// Ensure users directory exists
	if err := os.MkdirAll(filepath.Join(a.dataDir, "users"), 0755); err != nil {
		return "", fmt.Errorf("failed to create users directory: %v", err)
	}
	
	// Save user to file
	userData, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal user data: %v", err)
	}
	
	if err := ioutil.WriteFile(filePath, userData, 0644); err != nil {
		return "", fmt.Errorf("failed to save user: %v", err)
	}
	
	return userID, nil
}

// Login authenticates a user and returns a token
func (a *LocalAuth) Login(email, password string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Load all users and find matching email/password
	usersDir := filepath.Join(a.dataDir, "users")
	files, err := ioutil.ReadDir(usersDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to read users directory: %v", err)
	}
	
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		
		filePath := filepath.Join(usersDir, file.Name())
		userData, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}
		
		var user LocalUser
		if err := json.Unmarshal(userData, &user); err != nil {
			continue
		}
		
		if user.Email == email && user.Password == password {
			// Generate token
			token := a.generateToken()
			a.tokens[token] = &TokenInfo{
				UserID:    user.ID,
				Email:     user.Email,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}
			a.saveTokens()
			
			return token, nil
		}
	}
	
	return "", fmt.Errorf("invalid credentials")
}

func (a *LocalAuth) generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (a *LocalAuth) generateUserID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (a *LocalAuth) loadTokens() {
	filePath := filepath.Join(a.dataDir, "tokens.json")
	
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Warning: Failed to load tokens: %v", err)
		}
		return
	}
	
	if err := json.Unmarshal(data, &a.tokens); err != nil {
		log.Printf("Warning: Failed to unmarshal tokens: %v", err)
		a.tokens = make(map[string]*TokenInfo)
	}
}

func (a *LocalAuth) saveTokens() {
	filePath := filepath.Join(a.dataDir, "tokens.json")
	
	data, err := json.MarshalIndent(a.tokens, "", "  ")
	if err != nil {
		log.Printf("Warning: Failed to marshal tokens: %v", err)
		return
	}
	
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		log.Printf("Warning: Failed to save tokens: %v", err)
	}
}

func (a *LocalAuth) cleanupExpiredTokens() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		a.mu.Lock()
		now := time.Now()
		for token, info := range a.tokens {
			if now.After(info.ExpiresAt) {
				delete(a.tokens, token)
			}
		}
		a.saveTokens()
		a.mu.Unlock()
	}
} 