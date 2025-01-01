package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/imhasandl/go-restapi/internal/auth"
	"github.com/imhasandl/go-restapi/internal/database"
)

type ReportPost struct {
	ReportID  uuid.UUID `json:"report_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	Reason    string    `json:"reason"`
}

func (cfg *apiConfig) handlerReportPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		PostID uuid.UUID `json:"post_id"`
		Reason string    `json:"reason"`
	}
	type response struct {
		ReportPost
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse body - handlerReportPost", err)
		return
	}

	post, err := cfg.db.GetPostByID(r.Context(), params.PostID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get post by id - handlerReportPost", err)
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
		PostID:   post.ID,
		UserID:   userID,
		Reason:   params.Reason,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't report the post - handlerReportPost", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ReportPost: ReportPost{
			ReportID:  report.ReportID,
			CreatedAt: report.CreatedAt,
			UpdatedAt: report.UpdatedAt,
			PostID:    report.PostID,
			UserID:    report.UserID,
			Reason:    report.Reason,
		},
	})
}

func (cfg *apiConfig) handlerListAllReports(w http.ResponseWriter, r *http.Request) {
	reports, err := cfg.db.ListAllReports(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get all reports from db - handlerListAllReports", err)
		return
	}

	respondWithJSON(w, http.StatusOK, reports)
}

func (cfg *apiConfig) handlerDeleteReportByID(w http.ResponseWriter, r *http.Request) {
	reportIDString := r.PathValue("report_id")
	reportID, err := uuid.Parse(reportIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get path value - handlerDeleteReportByID", err)
		return
	}

	err = cfg.db.DeleteReportByID(r.Context(), reportID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't delete report by id - handlerDeleteReportByID", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerGetReportByID(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ReportPost
	}

	reportIDString := r.PathValue("report_id")
	reportID, err := uuid.Parse(reportIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't parse report id - handlerGetReportByID", err)
		return
	}

	report, err := cfg.db.GetReportByID(r.Context(), reportID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't get post by id - handlerGetReportByID", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ReportPost: ReportPost{
			ReportID:  report.ReportID,
			CreatedAt: report.CreatedAt,
			UpdatedAt: report.UpdatedAt,
			PostID:    report.PostID,
			UserID:    report.UserID,
			Reason:    report.Reason,
		},
	})
}
