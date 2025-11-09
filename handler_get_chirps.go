package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	chirpString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse given chirpID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to find chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps from db", err)
	}

	var chirpsArr []Chirp

	for _, chirp := range chirps {
		chirpsArr = append(chirpsArr, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsArr)
}
