package main

import (
	connectFirestore "learnfirestore/firestore"
	"learnfirestore/models"
	"learnfirestore/repository"
	"learnfirestore/seeds"
)

func main() {

	ctx, fs, err := connectFirestore.ConnectToFirestore()
	if err != nil {
		panic(err)
	}

	//notif_repo := repository.NewNotificationRepository(fs)
	//user_repo := repository.NewUserRepository(fs)
	group_repo := repository.NewGroupRepository(fs)

	// var newUser []models.User
	// newUser, err = seeds.UserFactory(ctx, user_repo, 1)
	// if err != nil {
	// 	panic(err)
	// }

	member := models.Member{ID: "yVVnNxl9GBF0ag2dWmd8", Name: "thick fast"}

	newGroup, err := seeds.GroupFactory(ctx, group_repo, member, 1)
	if err != nil {
		panic(err)
	}

	err = seeds.CommentsFactory(ctx, group_repo, member, newGroup[0], 1)
	if err != nil {
		panic(err)
	}

	err = seeds.TaskFactory(ctx, group_repo, member, newGroup[0], 1)
	if err != nil {
		panic(err)
	}

	// if err := notif_repo.WatchGroups(ctx); err != nil {
	// 	log.Printf("Watch error: %v", err)
	// }

}

// member := models.User{ID: newUser[0].ID, Name: newUser[0].Name}

// newGroup, err := groups.GroupFactory(ctx, fs, 1, member)
// if err != nil {
// 	panic(err)
// }

// err = groups.CommentsFactory(ctx, fs, member, newGroup[0], 1)
// if err != nil {
// 	panic(err)
// }

// err = groups.TaskFactory(ctx, fs, member, newGroup[0], 1)
// if err != nil {
// 	panic(err)
// }
