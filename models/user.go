package models

import "time"

type User struct {
	ID        string     `firestore:"-" json:"id"` //omitempty versi firestore?
	Name      string     `firestore:"name" json:"name"`
	Email     string     `firestore:"email" json:"email"`
	CreatedAt time.Time  `firestore:"createdAt" json:"created_at"`
	DeletedAt *time.Time `firestore:"deletedAt" json:"deleted_at,omitempty"`
}
