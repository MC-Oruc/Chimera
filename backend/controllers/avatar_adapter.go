package controllers

import (
	"backend/interfaces"
	"backend/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AvatarController struct {
	db interfaces.DatabaseService
}

func NewAvatarController(db interfaces.DatabaseService) *AvatarController {
	return &AvatarController{
		db: db,
	}
}

// CreateAvatar creates a new avatar
func (ac *AvatarController) CreateAvatar(c *gin.Context) {
	var req models.AvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	now := time.Now().Unix()
	avatar := models.Avatar{
		ID:              uuid.New().String(),
		Name:            req.Name,
		Description:     req.Description,
		Story:           req.Story,
		Persona:         req.Persona,
		ProfileImageURL: req.ProfileImageURL,
		IsPublic:        req.IsPublic,
		OwnerID:         userID,
		CreatorNickname: req.CreatorNickname,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err := ac.db.SaveAvatar(context.Background(), &avatar)
	if err != nil {
		log.Printf("Error creating avatar: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create avatar"})
		return
	}

	c.JSON(http.StatusCreated, avatar)
}

// GetAvatar gets a single avatar by ID
func (ac *AvatarController) GetAvatar(c *gin.Context) {
	avatarID := c.Param("id")
	if avatarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar ID is required"})
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get the avatar from database
	avatar, err := ac.db.GetAvatar(context.Background(), avatarID)
	if err != nil {
		log.Printf("Error getting avatar: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	// Check if the avatar is public or belongs to the user
	if !avatar.IsPublic && avatar.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this avatar"})
		return
	}

	// If creator nickname is empty, try to get the user's display name
	if avatar.CreatorNickname == "" {
		// Get the user from database
		user, err := ac.db.GetUser(context.Background(), avatar.OwnerID)
		if err == nil && user.DisplayName != "" {
			avatar.CreatorNickname = user.DisplayName
		}
		
		// If we still don't have a nickname, use a default
		if avatar.CreatorNickname == "" {
			avatar.CreatorNickname = "Anonymous Creator"
		}
	}

	// Return the avatar with ownership information
	avatarResponse := models.AvatarResponse{
		Avatar:  *avatar,
		IsOwner: avatar.OwnerID == userID,
	}

	c.JSON(http.StatusOK, avatarResponse)
}

// GetUserAvatars gets all avatars owned by the current user
func (ac *AvatarController) GetUserAvatars(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	avatars, err := ac.db.GetUserAvatars(context.Background(), userID)
	if err != nil {
		log.Printf("Error getting user avatars: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve avatars"})
		return
	}

	var avatarResponses []models.AvatarResponse
	for _, avatar := range avatars {
		avatarResponse := models.AvatarResponse{
			Avatar:  *avatar,
			IsOwner: true, // All are owned by the user
		}
		avatarResponses = append(avatarResponses, avatarResponse)
	}

	c.JSON(http.StatusOK, gin.H{"avatars": avatarResponses})
}

// GetPublicAvatars gets all public avatars for the marketplace
func (ac *AvatarController) GetPublicAvatars(c *gin.Context) {
	userID := c.GetString("userId")
	// For public avatars, we allow access even without authentication

	avatars, err := ac.db.GetPublicAvatars(context.Background())
	if err != nil {
		log.Printf("Error getting public avatars: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve avatars"})
		return
	}

	var avatarResponses []models.AvatarResponse
	for _, avatar := range avatars {
		avatarResponse := models.AvatarResponse{
			Avatar:  *avatar,
			IsOwner: userID != "" && avatar.OwnerID == userID,
		}
		avatarResponses = append(avatarResponses, avatarResponse)
	}

	c.JSON(http.StatusOK, gin.H{"avatars": avatarResponses})
}

// UpdateAvatar updates an existing avatar
func (ac *AvatarController) UpdateAvatar(c *gin.Context) {
	avatarID := c.Param("id")
	if avatarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar ID is required"})
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.AvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the existing avatar
	avatar, err := ac.db.GetAvatar(context.Background(), avatarID)
	if err != nil {
		log.Printf("Error getting avatar: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	// Check ownership
	if avatar.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this avatar"})
		return
	}

	// Update avatar fields
	avatar.Name = req.Name
	avatar.Description = req.Description
	avatar.Story = req.Story
	avatar.Persona = req.Persona
	avatar.ProfileImageURL = req.ProfileImageURL
	avatar.IsPublic = req.IsPublic
	avatar.CreatorNickname = req.CreatorNickname
	avatar.UpdatedAt = time.Now().Unix()

	err = ac.db.UpdateAvatar(context.Background(), avatar)
	if err != nil {
		log.Printf("Error updating avatar: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}

	c.JSON(http.StatusOK, avatar)
}

// DeleteAvatar deletes an avatar
func (ac *AvatarController) DeleteAvatar(c *gin.Context) {
	avatarID := c.Param("id")
	if avatarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar ID is required"})
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get the existing avatar
	avatar, err := ac.db.GetAvatar(context.Background(), avatarID)
	if err != nil {
		log.Printf("Error getting avatar: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	// Check ownership
	if avatar.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this avatar"})
		return
	}

	err = ac.db.DeleteAvatar(context.Background(), avatarID)
	if err != nil {
		log.Printf("Error deleting avatar: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete avatar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar deleted successfully"})
} 