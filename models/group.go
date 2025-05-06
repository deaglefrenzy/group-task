package models

import "time"

type Group struct {
	ID          string            `firestore:"-" json:"id"`
	Name        string            `firestore:"name" json:"name"`
	Description string            `firestore:"description" json:"description"`
	Members     map[string]Member `firestore:"members" json:"members"`
	MembersID   []string          `firestore:"membersID" json:"members_id"`
	Tasks       []Tasks           `firestore:"tasks" json:"tasks,omitempty"`
	Comments    []Comments        `firestore:"comments" json:"comments,omitempty"`
	CreatedAt   time.Time         `firestore:"createdAt" json:"created_at"`
	DeletedAt   *time.Time        `firestore:"deletedAt" json:"deleted_at,omitempty"`
}

type Member struct {
	ID   string `firestore:"id" json:"id"`
	Name string `firestore:"name" json:"name"`
}

type Tasks struct {
	UUID        string     `firestore:"id" json:"id"`
	Title       string     `firestore:"title" json:"title"`
	Description string     `firestore:"description" json:"description"`
	Priority    bool       `firestore:"priority" json:"priority"`
	Done        bool       `firestore:"done" json:"done"`
	DueDate     time.Time  `firestore:"dueDate" json:"due_date"`
	CreatedBy   Member     `firestore:"createdBy" json:"created_by"`
	CreatedAt   time.Time  `firestore:"createdAt" json:"created_at"`
	DeletedAt   *time.Time `firestore:"deletedAt" json:"deleted_at,omitempty"`
}

type Comments struct {
	UUID      string     `firestore:"id" json:"id"`
	Text      string     `firestore:"text" json:"text"`
	CreatedBy Member     `firestore:"createdBy" json:"created_by"`
	CreatedAt time.Time  `firestore:"createdAt" json:"created_at"`
	DeletedAt *time.Time `firestore:"deletedAt" json:"deleted_at,omitempty"`
}
