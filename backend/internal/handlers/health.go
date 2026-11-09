package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Tasks2/hmis/internal/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func Health(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			response.JSON(w, http.StatusServiceUnavailable, HealthResponse{
				Status:   "unhealthy",
				Database: "unavailable",
			})
			return
		}

		response.JSON(w, http.StatusOK, HealthResponse{
			Status:   "ok",
			Database: "ok",
		})
	}
}
