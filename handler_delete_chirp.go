package main

import (
	"net/http"

	"github.com/JoeVinten/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	userID, ok := getUserID(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "You don't own that chirp", nil)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirp.ID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
