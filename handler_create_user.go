package main

import (
	"encoding/json"
	"net/http"

	"github.com/JoeVinten/chirpy/internal/auth"
	"github.com/JoeVinten/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "issue decoding params", err)
		return
	}

	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "no email given", nil)
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing the password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPW,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
