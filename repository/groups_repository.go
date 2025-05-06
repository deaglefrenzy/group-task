package repository

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"time"

	"cloud.google.com/go/firestore"
)

type GroupRepository struct {
	fs       *firestore.Client
	coll     *firestore.CollectionRef
	collName string
}

func NewGroupRepository(fs *firestore.Client) *GroupRepository {
	return &GroupRepository{
		fs:       fs,
		coll:     fs.Collection("groups"),
		collName: "groups",
	}
}

func (r *GroupRepository) WithCollection(collectionName string) *GroupRepository {
	return &GroupRepository{
		fs:       r.fs,
		coll:     r.fs.Collection(collectionName),
		collName: collectionName,
	}
}

func (r GroupRepository) GetGroupWithID(ctx context.Context, id string) (models.Group, error) {
	snap, err := r.coll.Doc(id).Get(ctx)
	if err != nil {
		return models.Group{}, fmt.Errorf("failed to get group: %w", err)
	}

	var group models.Group
	if err := snap.DataTo(&group); err != nil {
		return models.Group{}, fmt.Errorf("failed to return group: %w", err)
	}

	return group, nil
}

func (r GroupRepository) CreateGroup(ctx context.Context, name string, description string, member models.Member) (models.Group, error) {
	groupRef := r.coll.NewDoc()
	group := models.Group{
		ID:          groupRef.ID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
	newMember := models.Member{
		ID:   member.ID,
		Name: member.Name,
	}

	group.Members[member.ID] = newMember

	if _, err := groupRef.Set(ctx, group); err != nil {
		return models.Group{}, fmt.Errorf("failed to create Group: %w", err)
	}
	return group, nil
}

func (r GroupRepository) UpdateGroup(ctx context.Context, id string, group *models.Group) error {
	if id == "" {
		return fmt.Errorf("id group cannot be empty")
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	if _, err := r.coll.Doc(id).Set(ctx, group, firestore.MergeAll); err != nil {
		return fmt.Errorf("failed to update Group: %w", err)
	}
	return nil
}

func (r GroupRepository) DeleteGroup(ctx context.Context, id string) error {

	if id == "" {
		return fmt.Errorf("group cannot be empty")
	}

	updates := []firestore.Update{
		{Path: "deletedAt", Value: time.Now()},
	}

	_, err := r.coll.Doc(id).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}
