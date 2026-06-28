package services

import (
	"context"
	"errors"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type FirebaseAuthService struct {
	client *auth.Client
}

type FirebaseTokenClaims struct {
	UID      string
	Email    string
	Name     string
	Verified bool
}

func NewFirebaseAuthService(app *firebase.App) *FirebaseAuthService {
	if app == nil {
		return &FirebaseAuthService{}
	}
	c, err := app.Auth(context.Background())
	if err != nil {
		return &FirebaseAuthService{}
	}
	return &FirebaseAuthService{client: c}
}

func (s *FirebaseAuthService) VerifyIDToken(idToken string) (*FirebaseTokenClaims, error) {
	if s.client == nil {
		// Fallback untuk dev: ekstrak UID dari JWT secara sederhana
		// (jangan pakai di production!)
		parts := strings.Split(idToken, ".")
		if len(parts) < 2 {
			return nil, errors.New("invalid firebase token (dev mode)")
		}
		return &FirebaseTokenClaims{
			UID:      "dev_" + parts[1][:min(20, len(parts[1]))],
			Email:    "dev@local.test",
			Name:     "Dev User",
			Verified: true,
		}, nil
	}
	tok, err := s.client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, err
	}
	out := &FirebaseTokenClaims{UID: tok.UID, Verified: true}
	if tok.Claims["email"] != nil {
		out.Email = tok.Claims["email"].(string)
	}
	if tok.Claims["email_verified"] != nil {
		out.Verified = tok.Claims["email_verified"].(bool)
	}
	if tok.Claims["name"] != nil {
		out.Name = tok.Claims["name"].(string)
	}
	return out, nil
}

func (s *FirebaseAuthService) SetEmailVerified(ctx context.Context, uid string) error {
	if s.client == nil {
		return errors.New("firebase client tidak tersedia")
	}
	_, err := s.client.UpdateUser(ctx, uid, (&auth.UserToUpdate{}).EmailVerified(true))
	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
