package services

import (
	"backend/interfaces"
	"context"

	"firebase.google.com/go/v4/auth"
)

type FirebaseAuth struct {
	client *auth.Client
}

// NewFirebaseAuth creates a new Firebase authentication service
func NewFirebaseAuth(client *auth.Client) interfaces.AuthService {
	return &FirebaseAuth{
		client: client,
	}
}

func (a *FirebaseAuth) VerifyToken(ctx context.Context, token string) (*interfaces.AuthUser, error) {
	authToken, err := a.client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}
	
	return &interfaces.AuthUser{
		UID:   authToken.UID,
		Email: authToken.Claims["email"].(string),
	}, nil
}

func (a *FirebaseAuth) ValidateToken(token string) bool {
	_, err := a.VerifyToken(context.Background(), token)
	return err == nil
} 