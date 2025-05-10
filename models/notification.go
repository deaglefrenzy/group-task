package models

import "time"

type Notification struct {
	ID        string       `firestore:"-" json:"id"`
	User      NotifiedUser `firestore:"user" json:"user"`
	Message   string       `firestore:"message" json:"message"`
	GroupID   string       `firestore:"groupID,omitempty" json:"group_id,omitempty"`
	Read      bool         `firestore:"read" json:"read"`
	CreatedAt time.Time    `firestore:"createdAt" json:"created_at"`
}

type NotifiedUser struct {
	ID   string `firestore:"id" json:"id"`
	Name string `firestore:"name" json:"name"`
}
