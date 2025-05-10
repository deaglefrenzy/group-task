package main

import (
	"context"
	connectFirestore "learnfirestore/firestore"
	"learnfirestore/handlers"
	"learnfirestore/repository"
)

func main() {
	ctx := context.Background()

	fs, err := connectFirestore.ConnectToFirestore(ctx)
	if err != nil {
		panic(err)
	}

	notif_repo := repository.NewNotificationRepository(fs)
	group_repo := repository.NewGroupRepository(fs)

	handlers.NewGroupsWatcherHandler(ctx, group_repo, notif_repo)
	//user_repo := repository.NewUserRepository(fs)

	// var newUser []models.User
	// newUser, err = seeds.UserFactory(ctx, user_repo, 1)
	// if err != nil {
	// 	panic(err)
	// }

	// member := models.Member{ID: "yVVnNxl9GBF0ag2dWmd8", Name: "thick fast"}

	// newGroup, err := seeds.GroupFactory(ctx, group_repo, member, 1)
	// if err != nil {
	// 	panic(err)
	// }

	// err = seeds.CommentsFactory(ctx, group_repo, member, newGroup[0], 1)
	// if err != nil {
	// 	panic(err)
	// }

	// err = seeds.TaskFactory(ctx, group_repo, member, newGroup[0], 1)
	// if err != nil {
	// 	panic(err)
	// }

	<-ctx.Done()
}
