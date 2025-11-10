package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateAndVerifyJWT(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	token, err := MakeJWT(userID, string(secret), 5*time.Second)
	if err != nil {
		t.Fatalf("Error making token: %v", err)
	}

	if len(token) < 1 {
		t.Fatalf("Token output is empty")
	}

	id, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Error verfiying token %v", err)
	}

	if id != userID {
		t.Errorf("Wrong UserId in token. got %v, want %v", id, userID)
	}

}

func TestTokenExpires(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	token, err := MakeJWT(userID, string(secret), 1*time.Millisecond)
	if err != nil {
		t.Fatalf("Error making token: %v", err)
	}

	time.Sleep(2 * time.Millisecond)

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Token did not expire when it should have")
	} else {
		t.Log("Successfully caught expired token")
	}
}

func TestJWTInvalidSecret(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	invalidSecret := "invalidSecret"
	duration := 5 * time.Minute

	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("Error making token: %v", err)
	}

	_, err = ValidateJWT(token, invalidSecret)
	if err == nil {
		t.Error("Token validation passed with an invalid secret")
	}
}
