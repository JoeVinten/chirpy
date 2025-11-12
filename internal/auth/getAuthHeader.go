package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	return GetAuthHeader(headers, "Bearer")
}

func GetAPIKey(headers http.Header) (string, error) {
	return GetAuthHeader(headers, "ApiKey")
}

func GetAuthHeader(headers http.Header, key string) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	splitAuth := strings.Fields(authHeader)
	if len(splitAuth) < 2 || splitAuth[0] != key {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
