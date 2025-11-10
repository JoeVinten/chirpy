package main

import (
	"net/http"
	"time"

	"github.com/JoeVinten/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token found", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create JWT", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking the token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
