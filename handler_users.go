package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/imhasandl/go-restapi/internal/auth"
	"github.com/imhasandl/go-restapi/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	IsPremium bool      `json:"is_premium"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't decode body", err)
		return
	}

	_, err = cfg.db.CheckIfUsernameOrEmailTaken(r.Context(), database.CheckIfUsernameOrEmailTakenParams{
		Username: params.Username,
		Email: params.Email,
	})
	if err == nil {
		respondWithError(w, http.StatusUnauthorized, "username or email alredy taken", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't hash the password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:       uuid.New(),
		Email:    params.Email,
		Username: params.Username,
		Password: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't create user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.Username,
			IsPremium: user.IsPremium,
		},
	})
}

func (cfg *apiConfig) handlerGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't decode body", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get user by email", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.Username,
			IsPremium: user.IsPremium,
		},
	})
}

func (cfg *apiConfig) handlerGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "can't get user bu username - handlerGetUserByUsername", err)
		return
	}

	user, err := cfg.db.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "can't get user by username - handlerGetUserByUsername", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.Username,
			IsPremium: user.IsPremium,
		},
	})
}

func (cfg *apiConfig) handlerUserChange(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't decode body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't hash the password", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't get header token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't validate jwt", err)
		return
	}

	user, err := cfg.db.ChangeUser(r.Context(), database.ChangeUserParams{
		Email:    params.Email,
		Password: hashedPassword,
		ID:       userID,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't change user", err)
		return
	}

	respondWithJSON(w, http.StatusAccepted, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.Username,
			IsPremium: user.IsPremium,
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

	respondWithJSON(w, http.StatusOK, User{
		ID:        userID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Username:  user.Username,
		IsPremium: user.IsPremium,
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
