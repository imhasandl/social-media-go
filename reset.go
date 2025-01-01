package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.status != "ADMIN" {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("Reset users only allowed in admin environment"))
		if err != nil {
			log.Printf("Error writing JSON: %s", err)
		}
		return
	}
	
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't reset users", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Reset users table, completed"))
	if err != nil {
		log.Printf("Error writing JSON: %s", err)
	}
}

func (cfg *apiConfig) handlerResetPosts(w http.ResponseWriter, r *http.Request) {
	if cfg.status != "ADMIN" {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("Reset posts only allowed in admin environment"))
		if err != nil {
			log.Printf("Error writing JSON: %s", err)
		}
		return
	}

	err := cfg.db.ResetPosts(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't reset posts", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Reset posts table, completed"))
	if err != nil {
		log.Printf("Error writing JSON: %s", err)
	}
}

func (cfg *apiConfig) handlerResetReports(w http.ResponseWriter, r *http.Request) {
	if cfg.status != "ADMIN" {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("Reset posts only allowed in admin environment"))
		if err != nil {
			log.Printf("Error writing JSON: %s", err)
		}
		return
	}

	err := cfg.db.ResetReports(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't reset posts", err)
		return
	}

	err = cfg.db.ResetReports(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't reset posts", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Reset reports table, completed"))
	if err != nil {
		log.Printf("Error writing JSON: %s", err)
	}
}