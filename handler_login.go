package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/JoeVinten/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "issue decoding params", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user was not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "database error getting user", err)
		return
	}

	isCorrectPW, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error with password hashing", err)
		return
	}

	if !isCorrectPW {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	const defaultExpiryInSeconds = 3600

	expSecs := defaultExpiryInSeconds
	if params.ExpiresInSeconds != nil {
		if *params.ExpiresInSeconds <= 0 {
			expSecs = defaultExpiryInSeconds
		} else if *params.ExpiresInSeconds > defaultExpiryInSeconds {
			expSecs = defaultExpiryInSeconds
		} else {
			expSecs = *params.ExpiresInSeconds
		}
	}

	expires := time.Duration(expSecs) * time.Second

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expires)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})

}
