package seeds

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"learnfirestore/repository"
	"learnfirestore/utils"
	"math/rand"
	"time"

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

func SeedTask(createdBy models.Member) models.Tasks {

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

func SeedComment(createdBy models.Member) models.Comments {

	return models.Comments{
		UUID:      uuid.NewString(),
		Text:      utils.GenerateString(10),
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
}

func GroupFactory(ctx context.Context, groupRepo *repository.GroupRepository, member models.Member, count int) ([]models.Group, error) {

	var groups []models.Group
	for range count {
		g := SeedGroup()
		newGroup, err := groupRepo.CreateGroup(ctx, g.Name, g.Description, member)
		if err != nil {
			return nil, err
		}
		groups = append(groups, newGroup)
		fmt.Printf("Group %s (%s) created\n", g.Name, newGroup.ID)
	}
	return groups, nil
}

func TaskFactory(ctx context.Context, groupRepo *repository.GroupRepository, createdBy models.Member, group models.Group, count int) error {

	for range count {
		t := SeedTask(createdBy)
		_, err := groupRepo.CreateTask(ctx, group.ID, t.Title, t.Description, t.Priority, t.DueDate, t.CreatedBy)
		if err != nil {
			return err
		}
		fmt.Printf("Task %s created\n", t.Title)
	}

	return nil
}

func CommentsFactory(ctx context.Context, groupRepo *repository.GroupRepository, createdBy models.Member, group models.Group, count int) error {

	for range count {
		c := SeedComment(createdBy)
		_, err := groupRepo.CreateComment(ctx, group.ID, c.Text, c.CreatedBy)
		if err != nil {
			return err
		}
		fmt.Printf("Comment %s created\n", c.Text)
	}

	return nil
}
