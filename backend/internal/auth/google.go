package auth

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"
)

type GoogleVerifier struct {
	clientID string
}

func NewGoogleVerifier(clientID string) *GoogleVerifier {
	return &GoogleVerifier{clientID: clientID}
}

func (v *GoogleVerifier) Verify(ctx context.Context, idTokenValue string) (GoogleClaims, error) {
	if v.clientID == "" {
		return GoogleClaims{}, fmt.Errorf("google client id is not configured")
	}

	payload, err := idtoken.Validate(ctx, idTokenValue, v.clientID)
	if err != nil {
		return GoogleClaims{}, fmt.Errorf("validate google id token: %w", err)
	}

	claims := GoogleClaims{Subject: payload.Subject}
	if email, ok := payload.Claims["email"].(string); ok {
		claims.Email = email
	}
	if name, ok := payload.Claims["name"].(string); ok {
		claims.Name = name
	}
	return claims, nil
}
