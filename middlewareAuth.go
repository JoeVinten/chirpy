package main

import (
	"context"
	"net/http"

	"github.com/JoeVinten/chirpy/internal/auth"
	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

func (cfg *apiConfig) middlewareAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "No token ", err)
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token ", err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		handler(w, r.WithContext(ctx))
	}
}

func getUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}
