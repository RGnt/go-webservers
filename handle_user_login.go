package main

import (
	chirpauth "chirpy/internal/auth"
	"encoding/json"
	"net/http"
	"time"
)

func (cfg *apiConfig) hanlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}

	type retVals struct {
		Token string `json:"token"`
		User
	}

	expriationDuration := 3600

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "User not found", err)
		return
	}

	if params.ExpiresInSeconds != 0 && params.ExpiresInSeconds <= 3600 {
		expriationDuration = params.ExpiresInSeconds
	}

	err = chirpauth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	token, err := chirpauth.MakeJSWT(user.ID, cfg.secret, time.Duration(expriationDuration*1000*1000))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't generate token", err)
		return
	}

	userVals := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	returnBody := retVals{
		User:  userVals,
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, returnBody)
}
