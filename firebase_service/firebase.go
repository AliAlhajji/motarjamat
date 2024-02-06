package firebaseservice

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
)

func NewFirebaseApp(pathToCreds string) (*firebase.App, error) {
	// opt := option.WithCredentialsFile(pathToCreds)

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not initialize Firebase app: %v", err)
	}

	return app, nil
}
