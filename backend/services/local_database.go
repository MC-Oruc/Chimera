package services

import (
	"backend/interfaces"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LocalDatabase struct {
	dataDir string
	mu      sync.RWMutex
}

// NewLocalDatabase creates a new local database instance
func NewLocalDatabase() interfaces.DatabaseService {
	dataDir := "local_data"
	
	// Ensure data directory exists
	dirs := []string{
		filepath.Join(dataDir, "users"),
		filepath.Join(dataDir, "avatars"),
		filepath.Join(dataDir, "chats"),
		filepath.Join(dataDir, "images"),
		filepath.Join(dataDir, "settings"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Failed to create directory %s: %v", dir, err)
		}
	}
	
	log.Println("Local database initialized with directory:", dataDir)
	return &LocalDatabase{
		dataDir: dataDir,
	}
}

// Helper methods for file operations
func (db *LocalDatabase) saveToFile(filePath string, data interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	
	return ioutil.WriteFile(filePath, jsonData, 0644)
}

func (db *LocalDatabase) loadFromFile(filePath string, data interface{}) error {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}
	
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	
	return json.Unmarshal(jsonData, data)
}

func (db *LocalDatabase) listFiles(dir string) ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	
	var result []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			result = append(result, file.Name())
		}
	}
	return result, nil
}

// User operations
func (db *LocalDatabase) GetUser(ctx context.Context, userID string) (*models.User, error) {
	filePath := filepath.Join(db.dataDir, "users", fmt.Sprintf("%s.json", userID))
	var user models.User
	if err := db.loadFromFile(filePath, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *LocalDatabase) SaveUser(ctx context.Context, user *models.User) error {
	filePath := filepath.Join(db.dataDir, "users", fmt.Sprintf("%s.json", user.ID))
	return db.saveToFile(filePath, user)
}

// Avatar operations
func (db *LocalDatabase) GetAvatar(ctx context.Context, avatarID string) (*models.Avatar, error) {
	filePath := filepath.Join(db.dataDir, "avatars", fmt.Sprintf("%s.json", avatarID))
	var avatar models.Avatar
	if err := db.loadFromFile(filePath, &avatar); err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (db *LocalDatabase) GetUserAvatars(ctx context.Context, userID string) ([]*models.Avatar, error) {
	dir := filepath.Join(db.dataDir, "avatars")
	files, err := db.listFiles(dir)
	if err != nil {
		return nil, err
	}
	
	var avatars []*models.Avatar
	for _, file := range files {
		filePath := filepath.Join(dir, file)
		var avatar models.Avatar
		if err := db.loadFromFile(filePath, &avatar); err != nil {
			log.Printf("Warning: Failed to load avatar file %s: %v", file, err)
			continue
		}
		if avatar.OwnerID == userID {
			avatars = append(avatars, &avatar)
		}
	}
	return avatars, nil
}

func (db *LocalDatabase) GetPublicAvatars(ctx context.Context) ([]*models.Avatar, error) {
	dir := filepath.Join(db.dataDir, "avatars")
	files, err := db.listFiles(dir)
	if err != nil {
		return nil, err
	}
	
	var avatars []*models.Avatar
	for _, file := range files {
		filePath := filepath.Join(dir, file)
		var avatar models.Avatar
		if err := db.loadFromFile(filePath, &avatar); err != nil {
			log.Printf("Warning: Failed to load avatar file %s: %v", file, err)
			continue
		}
		if avatar.IsPublic {
			avatars = append(avatars, &avatar)
		}
	}
	return avatars, nil
}

func (db *LocalDatabase) SaveAvatar(ctx context.Context, avatar *models.Avatar) error {
	filePath := filepath.Join(db.dataDir, "avatars", fmt.Sprintf("%s.json", avatar.ID))
	return db.saveToFile(filePath, avatar)
}

func (db *LocalDatabase) UpdateAvatar(ctx context.Context, avatar *models.Avatar) error {
	avatar.UpdatedAt = time.Now().Unix()
	return db.SaveAvatar(ctx, avatar)
}

func (db *LocalDatabase) DeleteAvatar(ctx context.Context, avatarID string) error {
	filePath := filepath.Join(db.dataDir, "avatars", fmt.Sprintf("%s.json", avatarID))
	db.mu.Lock()
	defer db.mu.Unlock()
	return os.Remove(filePath)
}

// Chat operations
func (db *LocalDatabase) GetChat(ctx context.Context, chatID string) (*models.Chat, error) {
	filePath := filepath.Join(db.dataDir, "chats", fmt.Sprintf("%s.json", chatID))
	var chat models.Chat
	if err := db.loadFromFile(filePath, &chat); err != nil {
		return nil, err
	}
	return &chat, nil
}

func (db *LocalDatabase) GetUserChats(ctx context.Context, userID string) ([]*models.Chat, error) {
	dir := filepath.Join(db.dataDir, "chats")
	files, err := db.listFiles(dir)
	if err != nil {
		return nil, err
	}
	
	var chats []*models.Chat
	for _, file := range files {
		filePath := filepath.Join(dir, file)
		var chat models.Chat
		if err := db.loadFromFile(filePath, &chat); err != nil {
			log.Printf("Warning: Failed to load chat file %s: %v", file, err)
			continue
		}
		if chat.UserID == userID {
			chats = append(chats, &chat)
		}
	}
	return chats, nil
}

func (db *LocalDatabase) SaveChat(ctx context.Context, chat *models.Chat) error {
	filePath := filepath.Join(db.dataDir, "chats", fmt.Sprintf("%s.json", chat.ID))
	return db.saveToFile(filePath, chat)
}

func (db *LocalDatabase) UpdateChat(ctx context.Context, chat *models.Chat) error {
	chat.UpdatedAt = time.Now().Unix()
	return db.SaveChat(ctx, chat)
}

func (db *LocalDatabase) DeleteChat(ctx context.Context, chatID string) error {
	filePath := filepath.Join(db.dataDir, "chats", fmt.Sprintf("%s.json", chatID))
	db.mu.Lock()
	defer db.mu.Unlock()
	return os.Remove(filePath)
}

// Image operations
func (db *LocalDatabase) GetImage(ctx context.Context, imageID string) (*models.Image, error) {
	filePath := filepath.Join(db.dataDir, "images", fmt.Sprintf("%s.json", imageID))
	var image models.Image
	if err := db.loadFromFile(filePath, &image); err != nil {
		return nil, err
	}
	return &image, nil
}

func (db *LocalDatabase) GetUserImages(ctx context.Context, userID string) ([]*models.Image, error) {
	dir := filepath.Join(db.dataDir, "images")
	files, err := db.listFiles(dir)
	if err != nil {
		return nil, err
	}
	
	var images []*models.Image
	for _, file := range files {
		filePath := filepath.Join(dir, file)
		var image models.Image
		if err := db.loadFromFile(filePath, &image); err != nil {
			log.Printf("Warning: Failed to load image file %s: %v", file, err)
			continue
		}
		if image.UserID == userID {
			images = append(images, &image)
		}
	}
	return images, nil
}

func (db *LocalDatabase) SaveImage(ctx context.Context, image *models.Image) error {
	filePath := filepath.Join(db.dataDir, "images", fmt.Sprintf("%s.json", image.ID))
	return db.saveToFile(filePath, image)
}

func (db *LocalDatabase) DeleteImage(ctx context.Context, imageID string) error {
	filePath := filepath.Join(db.dataDir, "images", fmt.Sprintf("%s.json", imageID))
	db.mu.Lock()
	defer db.mu.Unlock()
	return os.Remove(filePath)
}

// Settings operations
func (db *LocalDatabase) GetUserSetting(ctx context.Context, userID, settingKey string) (map[string]interface{}, error) {
	filePath := filepath.Join(db.dataDir, "settings", fmt.Sprintf("%s_%s.json", userID, settingKey))
	var setting map[string]interface{}
	if err := db.loadFromFile(filePath, &setting); err != nil {
		return nil, err
	}
	return setting, nil
}

func (db *LocalDatabase) SaveUserSetting(ctx context.Context, userID, settingKey string, data map[string]interface{}) error {
	filePath := filepath.Join(db.dataDir, "settings", fmt.Sprintf("%s_%s.json", userID, settingKey))
	return db.saveToFile(filePath, data)
}

func (db *LocalDatabase) DeleteUserSetting(ctx context.Context, userID, settingKey string) error {
	filePath := filepath.Join(db.dataDir, "settings", fmt.Sprintf("%s_%s.json", userID, settingKey))
	db.mu.Lock()
	defer db.mu.Unlock()
	return os.Remove(filePath)
} 