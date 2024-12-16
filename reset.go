package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.status != "ADMIN" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset only allowed in admin environment"))
	}
	err := cfg.db.Reset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't reset users", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset users table, successfully completed"))
}