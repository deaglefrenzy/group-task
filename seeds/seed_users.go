package seeds

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"learnfirestore/repository"
	"learnfirestore/utils"
	"time"
)

func SeedUser() models.User {
	return models.User{
		Name:      utils.GenerateString(2),
		Email:     utils.GenerateString(1) + "@" + utils.GenerateString(1) + ".com",
		CreatedAt: time.Now(),
	}
}

func UserFactory(ctx context.Context, userRepo *repository.UserRepository, count int) ([]models.User, error) {

	var users []models.User
	for range count {
		u := SeedUser()
		newUser, err := userRepo.CreateUser(ctx, u.Name, u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, newUser)
		fmt.Printf("User %s created\n", u.Name)
	}

	return users, nil
}
