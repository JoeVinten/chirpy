package main

import (
	"encoding/json"
	"net/http"

	"github.com/JoeVinten/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil || cfg.polkaKey != apiKey {
		respondWithError(w, http.StatusUnauthorized, "Incorrect API key", err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	uID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	err = cfg.db.UpgradeUser(r.Context(), uID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
