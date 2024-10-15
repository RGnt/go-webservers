package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	//chirpyId := r.PathValue("chirpy_id")
	chirpyId, err := uuid.Parse(r.PathValue("chirpy_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request", err)
	}
	chirpyData, err := cfg.db.GetChirpByID(r.Context(), chirpyId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request", err)
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirpyData.ID,
		CreatedAt: chirpyData.CreatedAt,
		UpdatedAt: chirpyData.UpdatedAt,
		Body:      chirpyData.Body,
		UserId:    chirpyData.UserID,
	})
}
