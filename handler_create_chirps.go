package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/JoeVinten/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140

	if len(params.Body) >= maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanParams := database.CreateChirpParams{
		Body:   profanityFilter(params.Body),
		UserID: userID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), cleanParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp in database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	},
	)
}

func profanityFilter(t string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	const filter = "****"

	cleanBody := []string{}

	for _, word := range strings.Split(t, " ") {
		if slices.Contains(profanity, strings.ToLower(word)) {
			cleanBody = append(cleanBody, filter)
		} else {
			cleanBody = append(cleanBody, word)
		}
	}

	return strings.Join(cleanBody, " ")

}
