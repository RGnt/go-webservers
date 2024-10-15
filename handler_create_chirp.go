package main

import (
	chirpauth "chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	secret := os.Getenv("SECRET")
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
	}

	token, err := chirpauth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error acquiring token", nil)
		return
	}
	id, err := chirpauth.ValidateJWT(token, secret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleanMessage := checkProfanity(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanMessage,
		UserID: id,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
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
