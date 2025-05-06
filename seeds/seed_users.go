package seeds

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"learnfirestore/utils"
	"time"

	"cloud.google.com/go/firestore"
)

func SeedUser() models.User {
	return models.User{
		Name:      utils.GenerateString(2),
		Email:     utils.GenerateString(1) + "@" + utils.GenerateString(1) + ".com",
		CreatedAt: time.Now(),
	}
}

func UserFactory(ctx context.Context, fs *firestore.Client, count int) ([]models.User, error) {

	var user []models.User
	for range count {
		u := SeedUser()
		userColl := fs.Collection("users").NewDoc()
		u.ID = userColl.ID
		_, err := userColl.Set(ctx, u)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		user = append(user, u)
		fmt.Printf("User %s (%s) created\n", u.Name, u.ID)
	}
	return user, nil
}
