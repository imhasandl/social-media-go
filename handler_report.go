package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/imhasandl/go-restapi/internal/auth"
	"github.com/imhasandl/go-restapi/internal/database"
)

type ReportPostParams struct {
	ReportID  uuid.UUID `json:"report_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	Reason    string    `json:"reason"`
}

func (cfg *apiConfig) handlerReportPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		PostID uuid.UUID
		Reason string `json:"reason"`
	}
	type response struct {
		ReportPostParams
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse body - handlerReportPost", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid header bearer - GetBearerToken", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't validate JWT", err)
		return
	}

	report, err := cfg.db.ReportPost(r.Context(), database.ReportPostParams{
		ReportID: uuid.New(),
		PostID: params.PostID,
		UserID: userID,
		Reason: params.Reason,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't report the post - handlerReportPost", err)
		return
	}
	
	respondWithJSON(w, http.StatusOK, response{
		ReportPostParams: ReportPostParams{
			ReportID: report.ReportID,
			CreatedAt: report.CreatedAt,
			UpdatedAt: report.UpdatedAt,
			PostID: report.PostID,
			UserID: report.UserID,
			Reason: report.Reason,
		},
	})
}

func (cfg *apiConfig) handlerListAllReports(w http.ResponseWriter, r *http.Request) {
	
}

func (cfg *apiConfig) handlerDeleteReportByID(w http.ResponseWriter, r *http.Request) {

}
