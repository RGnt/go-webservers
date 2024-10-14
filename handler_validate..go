package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type profaneReturnVals struct {
		Cleaned_Body string `json:"cleaned_body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanMsg := checkProfanity(params.Body)

	respondWithJSON(w, http.StatusOK, profaneReturnVals{
		Cleaned_Body: cleanMsg,
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
