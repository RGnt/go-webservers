package main

import (
	"net/http"
)

func (cfg *apiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	data, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chiprs", err)
	}

	retChirps := []Chirp{}
	for _, chirp := range data {
		retChirps = append(retChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, retChirps)
}
