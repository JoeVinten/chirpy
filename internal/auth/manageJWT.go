package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	secretKey := []byte(tokenSecret)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return ss, err
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}

	uID, err := uuid.Parse(claims.Subject)

	if err != nil {
		return uuid.UUID{}, err
	}

	return uID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authH := headers.Get("Authorization")

	if authH == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.Fields(authH)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return strings.TrimSpace(parts[1]), nil
}
