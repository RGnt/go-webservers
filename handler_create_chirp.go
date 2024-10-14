package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type retVals = Chirp

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleanMessage := checkProfanity(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(),
		database.CreateChirpParams{
			Body:   cleanMessage,
			UserID: params.UserId})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request", err)
	}
	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, http.StatusOK, retVals{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserId:    chirp.UserID,
		Body:      chirp.Body,
	})
}

func checkProfanity(msg string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	splitMsg := strings.Split(msg, " ")
	for i, word := range splitMsg {
		if slices.Contains(profanity, strings.ToLower(word)) {
			splitMsg[i] = "****"
		}
	}
	return strings.Join(splitMsg, " ")
}
