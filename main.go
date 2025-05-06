package main

import (
	connectFirestore "learnfirestore/firestore"
	"learnfirestore/repository"
	"log"
)

func main() {

	ctx, fs, err := connectFirestore.ConnectToFirestore()
	if err != nil {
		panic(err)
	}

	notif_repo := repository.NewNotificationRepository(fs)
	//user_repo := repository.NewUserRepository(fs)
	//group_repo := repository.NewGroupRepository(fs)

	//notif_repo.GetDocumentWithID(ctx, "3NdE2nmCcgO0cHIBOMdH")

	go func() {
		if err := notif_repo.WatchGroups(ctx); err != nil {
			log.Printf("Watch error: %v", err)
		}
	}()

	select {}
}

// newUser, err := users.UserFactory(ctx, fs, 1)
// if err != nil {
// 	panic(err)
// }

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
