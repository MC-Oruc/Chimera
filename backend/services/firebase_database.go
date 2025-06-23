package services

import (
	"backend/interfaces"
	"backend/models"
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirebaseDatabase struct {
	client *firestore.Client
}

// NewFirebaseDatabase creates a new Firebase database instance
func NewFirebaseDatabase(client *firestore.Client) interfaces.DatabaseService {
	return &FirebaseDatabase{
		client: client,
	}
}

// User operations
func (db *FirebaseDatabase) GetUser(ctx context.Context, userID string) (*models.User, error) {
	doc, err := db.client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}
	
	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *FirebaseDatabase) SaveUser(ctx context.Context, user *models.User) error {
	_, err := db.client.Collection("users").Doc(user.ID).Set(ctx, user)
	return err
}

// Avatar operations
func (db *FirebaseDatabase) GetAvatar(ctx context.Context, avatarID string) (*models.Avatar, error) {
	doc, err := db.client.Collection("avatars").Doc(avatarID).Get(ctx)
	if err != nil {
		return nil, err
	}
	
	var avatar models.Avatar
	if err := doc.DataTo(&avatar); err != nil {
		return nil, err
	}
	return &avatar, nil
}

func (db *FirebaseDatabase) GetUserAvatars(ctx context.Context, userID string) ([]*models.Avatar, error) {
	iter := db.client.Collection("avatars").Where("ownerId", "==", userID).Documents(ctx)
	defer iter.Stop()
	
	var avatars []*models.Avatar
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		
		var avatar models.Avatar
		if err := doc.DataTo(&avatar); err != nil {
			continue
		}
		avatars = append(avatars, &avatar)
	}
	return avatars, nil
}

func (db *FirebaseDatabase) GetPublicAvatars(ctx context.Context) ([]*models.Avatar, error) {
	iter := db.client.Collection("avatars").Where("isPublic", "==", true).Documents(ctx)
	defer iter.Stop()
	
	var avatars []*models.Avatar
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		
		var avatar models.Avatar
		if err := doc.DataTo(&avatar); err != nil {
			continue
		}
		avatars = append(avatars, &avatar)
	}
	return avatars, nil
}

func (db *FirebaseDatabase) SaveAvatar(ctx context.Context, avatar *models.Avatar) error {
	_, err := db.client.Collection("avatars").Doc(avatar.ID).Set(ctx, avatar)
	return err
}

func (db *FirebaseDatabase) UpdateAvatar(ctx context.Context, avatar *models.Avatar) error {
	return db.SaveAvatar(ctx, avatar)
}

func (db *FirebaseDatabase) DeleteAvatar(ctx context.Context, avatarID string) error {
	_, err := db.client.Collection("avatars").Doc(avatarID).Delete(ctx)
	return err
}

// Chat operations
func (db *FirebaseDatabase) GetChat(ctx context.Context, chatID string) (*models.Chat, error) {
	doc, err := db.client.Collection("chats").Doc(chatID).Get(ctx)
	if err != nil {
		return nil, err
	}
	
	var chat models.Chat
	if err := doc.DataTo(&chat); err != nil {
		return nil, err
	}
	return &chat, nil
}

func (db *FirebaseDatabase) GetUserChats(ctx context.Context, userID string) ([]*models.Chat, error) {
	iter := db.client.Collection("chats").Where("userId", "==", userID).Documents(ctx)
	defer iter.Stop()
	
	var chats []*models.Chat
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		
		var chat models.Chat
		if err := doc.DataTo(&chat); err != nil {
			continue
		}
		chats = append(chats, &chat)
	}
	return chats, nil
}

func (db *FirebaseDatabase) SaveChat(ctx context.Context, chat *models.Chat) error {
	_, err := db.client.Collection("chats").Doc(chat.ID).Set(ctx, chat)
	return err
}

func (db *FirebaseDatabase) UpdateChat(ctx context.Context, chat *models.Chat) error {
	return db.SaveChat(ctx, chat)
}

func (db *FirebaseDatabase) DeleteChat(ctx context.Context, chatID string) error {
	_, err := db.client.Collection("chats").Doc(chatID).Delete(ctx)
	return err
}

// Image operations
func (db *FirebaseDatabase) GetImage(ctx context.Context, imageID string) (*models.Image, error) {
	doc, err := db.client.Collection("gallery").Doc(imageID).Get(ctx)
	if err != nil {
		return nil, err
	}
	
	var image models.Image
	if err := doc.DataTo(&image); err != nil {
		return nil, err
	}
	return &image, nil
}

func (db *FirebaseDatabase) GetUserImages(ctx context.Context, userID string) ([]*models.Image, error) {
	iter := db.client.Collection("gallery").Where("UserID", "==", userID).Documents(ctx)
	defer iter.Stop()
	
	var images []*models.Image
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		
		var image models.Image
		if err := doc.DataTo(&image); err != nil {
			continue
		}
		images = append(images, &image)
	}
	return images, nil
}

func (db *FirebaseDatabase) SaveImage(ctx context.Context, image *models.Image) error {
	_, err := db.client.Collection("gallery").Doc(image.ID).Set(ctx, image)
	return err
}

func (db *FirebaseDatabase) DeleteImage(ctx context.Context, imageID string) error {
	_, err := db.client.Collection("gallery").Doc(imageID).Delete(ctx)
	return err
}

// Settings operations
func (db *FirebaseDatabase) GetUserSetting(ctx context.Context, userID, settingKey string) (map[string]interface{}, error) {
	doc, err := db.client.Collection("users").Doc(userID).Collection("settings").Doc(settingKey).Get(ctx)
	if err != nil {
		return nil, err
	}
	
	data := make(map[string]interface{})
	if err := doc.DataTo(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (db *FirebaseDatabase) SaveUserSetting(ctx context.Context, userID, settingKey string, data map[string]interface{}) error {
	_, err := db.client.Collection("users").Doc(userID).Collection("settings").Doc(settingKey).Set(ctx, data)
	return err
}

func (db *FirebaseDatabase) DeleteUserSetting(ctx context.Context, userID, settingKey string) error {
	_, err := db.client.Collection("users").Doc(userID).Collection("settings").Doc(settingKey).Delete(ctx)
	return err
} 