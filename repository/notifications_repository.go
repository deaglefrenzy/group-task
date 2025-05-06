package repository

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"time"

	"cloud.google.com/go/firestore"
)

type NotificationRepository struct {
	fs       *firestore.Client
	coll     *firestore.CollectionRef
	collName string
}

func NewNotificationRepository(fs *firestore.Client) *NotificationRepository {
	return &NotificationRepository{
		fs:       fs,
		coll:     fs.Collection("notifications"),
		collName: "notifications",
	}
}

func (r *NotificationRepository) WithCollection(collectionName string) *NotificationRepository {
	return &NotificationRepository{
		fs:       r.fs,
		coll:     r.fs.Collection(collectionName),
		collName: collectionName,
	}
}

func (r NotificationRepository) GetNotificationWithID(ctx context.Context, id string) (models.Notification, error) {
	snap, err := r.coll.Doc(id).Get(ctx)
	if err != nil {
		return models.Notification{}, fmt.Errorf("failed to get notification")
	}

	var notification models.Notification
	if err := snap.DataTo(&notification); err != nil {
		return models.Notification{}, fmt.Errorf("failed to return notification")
	}

	return notification, nil
}

func (r NotificationRepository) CreateNotification(ctx context.Context, user models.NotifiedUser, reference models.Reference, message string) (models.Notification, error) {
	notifRef := r.coll.NewDoc()
	notification := models.Notification{
		ID:        notifRef.ID,
		User:      user,
		Message:   message,
		Reference: reference,
		Read:      false,
		CreatedAt: time.Now(),
	}
	if _, err := notifRef.Set(ctx, notification); err != nil {
		return models.Notification{}, fmt.Errorf("failed to create notification: %w", err)
	}
	return notification, nil
}

func (r NotificationRepository) UpdateNotification(ctx context.Context, id string, notif *models.Notification) error {
	if id == "" {
		return fmt.Errorf("notification id cannot be empty")
	}
	if notif == nil {
		return fmt.Errorf("notification not found")
	}

	if _, err := r.coll.Doc(id).Set(ctx, notif, firestore.MergeAll); err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}
	return nil
}

func (r NotificationRepository) DeleteNotification(ctx context.Context, id string) error {

	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	updates := []firestore.Update{
		{Path: "deletedAt", Value: time.Now()},
	}

	_, err := r.coll.Doc(id).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

func (r NotificationRepository) WatchGroups(ctx context.Context) error {
	coll := r.WithCollection("groups").coll
	snap := coll.Snapshots(ctx)

	fmt.Println("Start watching Groups Collection...")

	for {
		qs, err := snap.Next()
		if err != nil {
			panic(err)
		}

		for _, v := range qs.Changes {
			groupID := v.Doc.Ref.ID

			groupRef := r.WithCollection("groups").coll.Doc(groupID)
			groupDoc, err := groupRef.Get(ctx)
			if err != nil {
				return fmt.Errorf("failed to get group: %w", err)
			}

			// REMEMBER: YOU MUST CONVERT THIS TO STRUCTURE DO NOT USE .DATA
			var groupStruct models.Group
			if err := groupDoc.DataTo(&groupStruct); err != nil {
				panic(err)
			}

			reference := models.Reference{
				Type: "group",
				ID:   groupID,
			}

			if v.Kind == firestore.DocumentAdded {
				for _, m := range groupStruct.Members {
					r.CreateNotification(ctx, models.NotifiedUser(m), reference, "New group have been created with you in it.")
				}
			} else if v.Kind == firestore.DocumentModified {
				for _, m := range groupStruct.Members {
					r.CreateNotification(ctx, models.NotifiedUser(m), reference, "There have been modifications in the group.")
				}
			} else if v.Kind == firestore.DocumentRemoved {
				for _, m := range groupStruct.Members {
					r.CreateNotification(ctx, models.NotifiedUser(m), reference, "Group have been deleted.")
				}
			}
		}
	}
}
