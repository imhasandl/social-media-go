package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/imhasandl/go-restapi/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't create user", err)
		return
	}

	userCreateResponse := parameters{
		ID:    uuid.New(),
		Email: params.Email,
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams(userCreateResponse))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't create user", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}

func (cfg *apiConfig) handlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	userIDString := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get the user by id", err)
		return
	}

	respondWithJSON(w, http.StatusFound, User{
		ID:        userID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *apiConfig) handlerListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.db.ListAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't list all users", err)
		return
	}

	respondWithJSON(w, http.StatusAccepted, users)
}
