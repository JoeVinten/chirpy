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

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}

//     We already have a bit of authorization built into Chirpy: authenticated users can only create chirps for themselves, not for others.

// Add a PUT /api/users endpoint so that users can update their own (but not others') email and password. It requires:
// An access token in the header
// A new password and email in the request body
// Hash the password, then update the hashed password and the email for the authenticated user in the database. Respond with a 200 if everything is successful and the newly updated User resource (omitting the password of course).
// If the access token is malformed or missing, respond with a 401 status code.
// }
