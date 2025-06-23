package interfaces

import (
	"backend/models"
	"context"
)

// DatabaseService defines the interface for database operations
type DatabaseService interface {
	// Users
	GetUser(ctx context.Context, userID string) (*models.User, error)
	SaveUser(ctx context.Context, user *models.User) error
	
	// Avatars
	GetAvatar(ctx context.Context, avatarID string) (*models.Avatar, error)
	GetUserAvatars(ctx context.Context, userID string) ([]*models.Avatar, error)
	GetPublicAvatars(ctx context.Context) ([]*models.Avatar, error)
	SaveAvatar(ctx context.Context, avatar *models.Avatar) error
	UpdateAvatar(ctx context.Context, avatar *models.Avatar) error
	DeleteAvatar(ctx context.Context, avatarID string) error
	
	// Chats
	GetChat(ctx context.Context, chatID string) (*models.Chat, error)
	GetUserChats(ctx context.Context, userID string) ([]*models.Chat, error)
	SaveChat(ctx context.Context, chat *models.Chat) error
	UpdateChat(ctx context.Context, chat *models.Chat) error
	DeleteChat(ctx context.Context, chatID string) error
	
	// Images
	GetImage(ctx context.Context, imageID string) (*models.Image, error)
	GetUserImages(ctx context.Context, userID string) ([]*models.Image, error)
	SaveImage(ctx context.Context, image *models.Image) error
	DeleteImage(ctx context.Context, imageID string) error
	
	// Settings
	GetUserSetting(ctx context.Context, userID, settingKey string) (map[string]interface{}, error)
	SaveUserSetting(ctx context.Context, userID, settingKey string, data map[string]interface{}) error
	DeleteUserSetting(ctx context.Context, userID, settingKey string) error
}

// StorageService defines the interface for file storage operations
type StorageService interface {
	UploadFromURL(sourceURL, destinationPath string) (string, error)
	UploadBase64Image(base64Image, destinationPath string) (string, error)
	DeleteFile(filePath string) error
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	VerifyToken(ctx context.Context, token string) (*AuthUser, error)
	ValidateToken(token string) bool
}

// AuthUser represents an authenticated user
type AuthUser struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
} 