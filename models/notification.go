package models

import "time"

type Notification struct {
	ID        string       `firestore:"-" json:"id"`
	User      NotifiedUser `firestore:"user" json:"user"`
	Message   string       `firestore:"message" json:"message"`
	Reference Reference    `firestore:"reference,omitempty" json:"reference,omitempty"`
	Read      bool         `firestore:"read" json:"read"`
	CreatedAt time.Time    `firestore:"createdAt" json:"created_at"`
}

type Reference struct {
	Type string `firestore:"type,omitempty" json:"type,omitempty"`
	ID   string `firestore:"id,omitempty" json:"id,omitempty"`
}

type NotifiedUser struct {
	ID   string `firestore:"id" json:"id"`
	Name string `firestore:"name" json:"name"`
}
