package main

import (
	"encoding/json"
	"net/http"

	"github.com/JoeVinten/chirpy/internal/auth"
	"github.com/JoeVinten/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateAccount(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	userID, ok := getUserID(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.db.UpdateUsernamePassword(r.Context(), database.UpdateUsernamePasswordParams{
		Email:          params.Email,
		HashedPassword: hashedPW,
		ID:             userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating the user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})

}
