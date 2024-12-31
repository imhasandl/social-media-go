package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/imhasandl/go-restapi/internal/auth"
	"github.com/imhasandl/go-restapi/internal/database"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type responce struct {
		Post
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid header bearer", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't decode the body", err)
		return
	}

	post, err := cfg.db.CreatePost(r.Context(), database.CreatePostParams{
		ID:     uuid.New(),
		UserID: userID,
		Body:   params.Body,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't post it", err)
		return
	}

	respondWithJSON(w, http.StatusOK, responce{
		Post: Post{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			UserID:    post.UserID,
			Body:      post.Body,
		},
	})
}

func (cfg *apiConfig) handlerListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := cfg.db.GetPosts(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't list the posts", err)
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func (cfg *apiConfig) handlerGetPostByID(w http.ResponseWriter, r *http.Request) {
	postIDString := r.PathValue("post_id")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse the post id provided", err)
		return
	}

	post, err := cfg.db.GetPostByID(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get post from db", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Post{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		UserID:    post.UserID,
		Body:      post.Body,
	})
}

func (cfg *apiConfig) handlerChangePostByID(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	postIDString := r.PathValue("post_id")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse the post id", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "can't decode body")
		return
	}

	err = cfg.db.ChangePostByID(r.Context(), database.ChangePostByIDParams{
		Body: params.Body,
		ID:   postID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't change the post by id", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerDeletePostByID(w http.ResponseWriter, r *http.Request) {
	postIDString := r.PathValue("post_id")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse the post id", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't get header token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "jwt token is not correct one", err)
		return
	}

	post, err := cfg.db.GetPostByID(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't get the post by id", err)
		return
	}

	if post.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you can't delete this chirp", err)
		return
	}

	err = cfg.db.DeletePostByID(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't delete the post", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}