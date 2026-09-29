package router

import (
	"net/http"

	"github.com/Tasks2/hmis/internal/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.Health(db))

	mux.HandleFunc("/api/v1/auth", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	})

	mux.HandleFunc("/api/v1/practitioners", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	})

	mux.HandleFunc("/api/v1/schedules", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	})

	mux.HandleFunc("/api/v1/appointments", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	})

	return mux
}
