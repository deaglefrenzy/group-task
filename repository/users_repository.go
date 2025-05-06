package repository

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"time"

	"cloud.google.com/go/firestore"
)

type UserRepository struct {
	fs       *firestore.Client
	coll     *firestore.CollectionRef
	collName string
}

func NewUserRepository(fs *firestore.Client) *UserRepository {
	return &UserRepository{
		fs:       fs,
		coll:     fs.Collection("users"),
		collName: "users",
	}
}

func (r *UserRepository) WithCollection(collectionName string) *UserRepository {
	return &UserRepository{
		fs:       r.fs,
		coll:     r.fs.Collection(collectionName),
		collName: collectionName,
	}
}

func (r UserRepository) GetUserWithID(ctx context.Context, id string) (models.User, error) {
	snap, err := r.coll.Doc(id).Get(ctx)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	var User models.User
	if err := snap.DataTo(&User); err != nil {
		return models.User{}, fmt.Errorf("failed to return user: %w", err)
	}

	return User, nil
}

func (r UserRepository) CreateUser(ctx context.Context, name string, email string) (models.User, error) {
	userRef := r.coll.NewDoc()
	user := models.User{
		ID:        userRef.ID,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	if _, err := userRef.Set(ctx, user); err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (r UserRepository) UpdateUser(ctx context.Context, id string, user *models.User) error {
	if id == "" {
		return fmt.Errorf("id user cannot be empty")
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if _, err := r.coll.Doc(id).Set(ctx, user, firestore.MergeAll); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r UserRepository) DeleteUser(ctx context.Context, id string) error {

	if id == "" {
		return fmt.Errorf("user cannot be empty")
	}

	updates := []firestore.Update{
		{Path: "deletedAt", Value: time.Now()},
	}

	_, err := r.coll.Doc(id).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
