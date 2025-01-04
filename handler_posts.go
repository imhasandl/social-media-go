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
	Likes     int32     `json:"likes"`
}

type PostsLike struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
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
			Likes:     post.Likes,
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
		respondWithError(w, http.StatusBadRequest, "can't parse the post id - handlerChangePostByID", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "can't decode body - handlerChangePostByID")
		return
	}

	err = cfg.db.ChangePostByID(r.Context(), database.ChangePostByIDParams{
		Body: params.Body,
		ID:   postID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't change the post by id - handlerChangePostByID", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerDeletePostByID(w http.ResponseWriter, r *http.Request) {
	postIDString := r.PathValue("post_id")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse the post id - handlerDeletePostByID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't get header token - handlerDeletePostByID", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "jwt token is not correct one - handlerDeletePostByID", err)
		return
	}

	post, err := cfg.db.GetPostByID(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't get the post by id - handlerDeletePostByID", err)
		return
	}

	if post.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you can't delete this chirp - handlerDeletePostByID", err)
		return
	}

	err = cfg.db.DeletePostByID(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't delete the post - handlerDeletePostByID", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerLikePost(w http.ResponseWriter, r *http.Request) {
	type response struct {
		PostsLike
	}

	postIDString := r.PathValue("post_id")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse post id - handlerLikePost", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get header bearer - handlerLikePost", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't  validate jwt token - handlerLikePost", err)
		return
	}

	// Checks if the user already liked the post
	err = cfg.db.CheckIfUserLikeAlready(r.Context(), userID)
	if err == nil {
		respondWithError(w, http.StatusBadRequest, "you can only like this post once - handlerLikePost", err)
		return
	}

	postLike, err := cfg.db.LikePost(r.Context(), database.LikePostParams{
		ID:     uuid.New(),
		PostID: postID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't like post - handlerLikePost", err)
		return
	}

	// Increments the likes column in posts table
	err = cfg.db.IncrementPostLike(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't increment post's likes - handlerLikePost", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		PostsLike: PostsLike{
			ID:        postLike.ID,
			PostID:    postLike.PostID,
			UserID:    postLike.UserID,
			CreatedAt: postLike.CreatedAt,
		},
	})
}

func (cfg *apiConfig) handlerDislikePost(w http.ResponseWriter, r *http.Request) {
	postIDString := r.PathValue("likepost_id")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse post id - handlerDislikePost", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't get bearer from header - handlerDislikePost", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't validate jwt - handlerDislikePost", err)
		return
	}

	err = cfg.db.DislikePost(r.Context(), database.DislikePostParams{
		UserID: userID,
		PostID: postID,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "can't dislike the post - handlerDislikePost", err)
		return
	}

	// Decrements the likes column in posts table
	err = cfg.db.DecrementPostLike(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't decrement post's likes - handlerDislikePost", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerListLikePost(w http.ResponseWriter, r *http.Request) {
	likePosts, err := cfg.db.ListLikePost(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't list LikePost table - handlerListLikePost", err)
		return
	}

	respondWithJSON(w, http.StatusOK, likePosts)
}

func (cfg *apiConfig) handlerGetPostLikes(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Likes int32 `json:"likes"`
	}

	postIDString := r.PathValue("post_id")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse path value - handlerGetPostLikes", err)
		return
	}

	likes, err := cfg.db.GetPostLikes(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get post's likes - handlerGetPostLikes", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Likes: likes,
	})
}
