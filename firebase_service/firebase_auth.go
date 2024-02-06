package firebaseservice

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/AliAlhajji/Motarjamat/models"
)

type AuthService struct {
	authClient *auth.Client
}

func NewAuthService(ctx context.Context, app *firebase.App) (*AuthService, error) {
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create auth client: %v", err)
	}

	return &AuthService{
		authClient: authClient,
	}, nil

}

// Create a new Firebase user and return the uuid
func (s *AuthService) CreateUser(email string, password string) (string, error) {
	user := &auth.UserToCreate{}
	user.Email(email)
	user.Password(password)

	firebaseUser, err := s.authClient.CreateUser(context.Background(), user)
	if err != nil {
		return "", err
	}

	return firebaseUser.UID, nil
}

func (s *AuthService) VerifyToken(token string) (*models.User, error) {
	claims, err := s.authClient.VerifyIDToken(context.Background(), token)
	if err != nil {
		return nil, err
	}

	if claims.Expires < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	user, err := s.authClient.GetUser(context.Background(), claims.UID)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Email: user.Email,
		UUID:  user.UID,
	}, nil

}
