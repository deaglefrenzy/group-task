package handlers

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"learnfirestore/repository"
)

type GroupsWatcherHandler struct {
	group_repo *repository.GroupRepository
	notif_repo *repository.NotificationRepository
}

func NewGroupsWatcherHandler(ctx context.Context, group_repo *repository.GroupRepository, notif_repo *repository.NotificationRepository) *GroupsWatcherHandler {
	ctx2, cancel := context.WithCancel(ctx)

	h := &GroupsWatcherHandler{
		group_repo,
		notif_repo,
	}

	documentAdded := make(chan models.Group, 1)
	documentModified := make(chan models.GroupChanges, 1)
	documentRemoved := make(chan models.Group, 1)

	go func() {
		if err := h.group_repo.WatchGroups(ctx2, documentAdded, documentModified, documentRemoved); err != nil {
			cancel()
		}
	}()

	go func() {
		for {
			select {
			case <-ctx2.Done():
				return
			case added := <-documentAdded:
				NewGroupAdded(ctx, notif_repo, added)
			case modified := <-documentModified:
				MemberChanged(ctx, notif_repo, modified)
				TaskAdded(ctx, notif_repo, modified)
				CommentAdded(ctx, notif_repo, modified)
			case removed := <-documentRemoved:
				GroupDeleted(ctx, notif_repo, removed)
			}
		}
	}()

	return h
}

func NewGroupAdded(ctx context.Context, notif_repo *repository.NotificationRepository, group models.Group) {
	for _, member := range group.Members {
		_, err := notif_repo.CreateNotification(ctx, models.NotifiedUser(member), group.ID, "User "+member.Name+" have been added to the new group "+group.Name+".")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GroupDeleted(ctx context.Context, notif_repo *repository.NotificationRepository, group models.Group) {
	for _, member := range group.Members {
		_, err := notif_repo.CreateNotification(ctx, models.NotifiedUser(member), group.ID, "The group "+group.Name+" have been deleted.")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func MemberChanged(ctx context.Context, notif_repo *repository.NotificationRepository, changes models.GroupChanges) {
	addedMembers := make(map[string]models.Member)
	removedMembers := make(map[string]models.Member)

	for id, member := range changes.After.Members {
		if _, exists := changes.Before.Members[id]; !exists {
			addedMembers[id] = member
		}
	}

	for id, member := range changes.Before.Members {
		if _, exists := changes.After.Members[id]; !exists {
			removedMembers[id] = member
		}
	}
	for _, member := range addedMembers {
		for _, m := range changes.After.Members {
			notif_repo.CreateNotification(ctx, models.NotifiedUser(m), changes.After.ID, "Member "+member.Name+" have joined group "+changes.After.Name+".")
		}
	}

	for _, member := range removedMembers {
		for _, m := range changes.Before.Members {
			notif_repo.CreateNotification(ctx, models.NotifiedUser(m), changes.After.ID, "Member "+member.Name+" have been removed from group "+changes.Before.Name+".")
		}
	}
}

func TaskAdded(ctx context.Context, notif_repo *repository.NotificationRepository, changes models.GroupChanges) {
	var addedTasks []models.Tasks

	for _, afterTask := range changes.After.Tasks {
		found := false
		for _, task := range changes.Before.Tasks {
			if task.UUID == afterTask.UUID {
				found = true
			}
		}
		if !found {
			addedTasks = append(addedTasks, afterTask)
		}
	}
	for _, task := range addedTasks {
		for _, m := range changes.After.Members {
			notif_repo.CreateNotification(ctx, models.NotifiedUser(m), changes.After.ID, "Task "+task.Title+" have been added to the group "+changes.After.Name+".")
		}
	}
}

func CommentAdded(ctx context.Context, notif_repo *repository.NotificationRepository, changes models.GroupChanges) {
	var addedComments []models.Comments

	for _, afterComment := range changes.After.Comments {
		found := false
		for _, comment := range changes.Before.Comments {
			if comment.UUID == afterComment.UUID {
				found = true
			}
		}
		if !found {
			addedComments = append(addedComments, afterComment)
		}
	}
	for _, comment := range addedComments {
		for _, m := range changes.After.Members {
			notif_repo.CreateNotification(ctx, models.NotifiedUser(m), changes.After.ID, "New comment from "+comment.CreatedBy.Name+" in the group "+changes.After.Name+".")
		}
	}
}
