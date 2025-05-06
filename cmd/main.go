package main

import (
	connectFirestore "learnfirestore/firestore"
	"learnfirestore/repository"
	"learnfirestore/seeds"
)

func main() {
	ctx, fs, err := connectFirestore.ConnectToFirestore()
	if err != nil {
		panic(err)
	}

	notif_repo := repository.NewNotificationRepository(fs)
	seeds.SeedNotification(ctx, notif_repo)
}
