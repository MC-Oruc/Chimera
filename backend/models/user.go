package models

// User represents a user in the system
type User struct {
	ID          string `json:"id" firestore:"id"`
	Email       string `json:"email" firestore:"email"`
	DisplayName string `json:"displayName" firestore:"displayName"`
	PhotoURL    string `json:"photoURL" firestore:"photoURL"`
	CreatedAt   int64  `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt" firestore:"updatedAt"`
}

// UserRequest is used for creating or updating user data
type UserRequest struct {
	DisplayName string `json:"displayName"`
	PhotoURL    string `json:"photoURL"`
} 