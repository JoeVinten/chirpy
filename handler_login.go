package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/JoeVinten/chirpy/internal/auth"
	"github.com/JoeVinten/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},

		Token:        accessToken,
		RefreshToken: refreshToken,
	})

}
