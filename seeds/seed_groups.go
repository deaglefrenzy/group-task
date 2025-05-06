package seeds

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"learnfirestore/repository"
	"learnfirestore/utils"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

func SeedGroup() models.Group {
	return models.Group{
		Name:        utils.GenerateString(2),
		Description: utils.GenerateString(4),
		Members:     map[string]models.Member{},
		Tasks:       []models.Tasks{},
		Comments:    []models.Comments{},
		CreatedAt:   time.Now(),
	}
}

func GroupFactory(ctx context.Context, group_repo *repository.GroupRepository, notif_repo *repository.NotificationRepository, member models.Member, count int) ([]models.Group, error) {

	var group []models.Group
	for range count {
		g := SeedGroup()

		newGroup, err := group_repo.CreateGroup(ctx, g.Name, g.Description, member)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Group %s (%s) created\n", g.Name, newGroup.ID)

		reference := models.Reference{
			Type: "group",
			ID:   newGroup.ID,
		}

		user := models.NotifiedUser{
			ID:   member.ID,
			Name: member.Name,
		}

		_, err = notif_repo.CreateNotification(ctx, user, reference, "You are added into a group!")
		if err != nil {
			return []models.Group{}, err
		}
	}
	return group, nil
}

func SeedComment(createdBy models.Member) models.Comments {

	return models.Comments{
		UUID:      uuid.NewString(),
		Text:      utils.GenerateString(10),
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
}

func CommentsFactory(ctx context.Context, fs *firestore.Client, createdBy models.Member, group models.Group, count int) error {

	groupRef := fs.Collection("groups").Doc(group.ID)
	groupDoc, err := groupRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	membersData := groupDoc.Data()["members"].([]interface{})
	var membersID []string
	for _, m := range membersData {
		memberMap, ok := m.(map[string]interface{})
		if !ok {
			return nil
		}

		if id, exists := memberMap["id"].(string); exists {
			membersID = append(membersID, id)
		}
	}

	var comments []models.Comments
	if err := groupDoc.DataTo(&comments); err != nil {
		comments = []models.Comments{}
	}

	for range count {
		c := SeedComment(createdBy)
		comments = append(comments, c)
		fmt.Printf("Comment %s created\n", c.Text)
	}
	_, err = groupRef.Update(ctx, []firestore.Update{
		{
			Path:  "comments",
			Value: comments,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update group comments: %w", err)
	}

	for _, m := range membersID {
		for _, c := range comments {
			reference := models.Reference{
				Type: "comment",
				ID:   c.UUID,
			}
			err = repository.CreateNotification(ctx, fs, m, reference, "There's a new comment in the group.")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SeedTask(createdBy models.User) models.Tasks {

	return models.Tasks{
		UUID:        uuid.NewString(),
		Title:       utils.GenerateString(3),
		Description: utils.GenerateString(6),
		Priority:    rand.Intn(2) == 1,
		Done:        false,
		DueDate:     time.Now().AddDate(0, 0, 1),
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}
}

func TaskFactory(ctx context.Context, fs *firestore.Client, createdBy models.User, group models.Group, count int) error {

	groupRef := fs.Collection("groups").Doc(group.ID)
	groupDoc, err := groupRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	membersData := groupDoc.Data()["members"].([]interface{})
	var membersID []string
	for _, m := range membersData {
		memberMap, ok := m.(map[string]interface{})
		if !ok {
			return nil
		}

		if id, exists := memberMap["id"].(string); exists {
			membersID = append(membersID, id)
		}
	}

	var tasks []models.Tasks
	if err := groupDoc.DataTo(&tasks); err != nil {
		tasks = []models.Tasks{}
	}

	for range count {
		t := SeedTask(createdBy)
		tasks = append(tasks, t)
		fmt.Printf("Task %s created\n", t.Title)
	}
	_, err = groupRef.Update(ctx, []firestore.Update{
		{
			Path:  "tasks",
			Value: tasks,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update group tasks: %w", err)
	}

	for _, m := range membersID {
		for _, c := range tasks {
			reference := models.Reference{
				Type: "tasks",
				ID:   c.UUID,
			}
			err = notifications.CreateNotification(ctx, fs, m, reference, "There's a new task in the group.")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
