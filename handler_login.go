package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/imhasandl/go-restapi/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ExpiresInSecond int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't decode the params", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get user by email", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong password", err)
		return
	}

	expirationTime := time.Hour
	if params.ExpiresInSecond > 0 && params.ExpiresInSecond < 3600 {
		expirationTime = time.Duration(params.ExpiresInSecond) * time.Second
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusAccepted, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: accessToken,
	})
}
