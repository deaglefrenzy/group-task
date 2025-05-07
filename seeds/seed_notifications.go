package seeds

import (
	"learnfirestore/models"
	"learnfirestore/utils"
	"time"
)

func SeedNotification(user models.NotifiedUser) models.Notification {
	return models.Notification{
		User:      user,
		Message:   utils.GenerateString(6),
		Read:      false,
		CreatedAt: time.Now(),
	}
}

// func SeedNotification(ctx context.Context, notification_repo repository.NotificationRepository) {

// 	for i := range 10 {
// 		notification_repo.CreateNotification(ctx)
// 	}
// }
