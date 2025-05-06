package seeds

import (
	"context"
	"learnfirestore/repository"
)

func SeedNotification(ctx context.Context, notification_repo repository.NotificationRepository) {

	for i := range 10 {
		notification_repo.CreateNotification(ctx)
	}
}
