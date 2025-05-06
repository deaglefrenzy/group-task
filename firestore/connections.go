package connectFirestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func ConnectToFirestore() (context.Context, *firestore.Client, error) {
	ctx := context.Background() // digunakan untuk action yg memerlukan time

	opt := option.WithCredentialsFile("service.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	fs, err := app.Firestore(ctx) // requires connection to firestore database/server
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to initialize firestore client: %w", err)
	}

	return ctx, fs, nil
}
