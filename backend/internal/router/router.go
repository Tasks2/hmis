package router

import (
	"net/http"

	"github.com/Tasks2/hmis/internal/handlers"
	"github.com/Tasks2/hmis/internal/middleware"
	"github.com/Tasks2/hmis/internal/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool,
	authHandler *handlers.AuthHandler,
	pracHandler *handlers.PractitionerHandler,
	appHandler *handlers.AppointmentHandler,
	jwtSecret string,
) http.Handler {
	mux := http.NewServeMux()
	//Health
	mux.HandleFunc("/health", handlers.Health(db))

	//Authentication
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)

	//Practitioners

	mux.HandleFunc("/api/v1/practitioners", pracHandler.List)
	mux.HandleFunc("/api/v1/practitioners/", pracHandler.Availability)

	//Schedules
	mux.HandleFunc("/api/v1/schedules", func(w http.ResponseWriter, r *http.Request) {

		response.JSONError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Schedules endpoint is not implemented yet")
	})

	//Appointments
	protected := middleware.JWTAuth(jwtSecret)
	patientOnly := middleware.RequireRole("PATIENT")

	mux.Handle(
		"/api/v1/appointments",
		protected(
			patientOnly(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodPost:
						appHandler.Create(w, r)

					case http.MethodGet:
						appHandler.List(w, r)

					default:
						http.NotFound(w, r)
					}
				}),
			),
		),
	)

	mux.Handle(
		"/api/v1/appointments/",
		protected(
			patientOnly(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					switch r.Method {
					case http.MethodDelete:
						appHandler.Cancel(w, r)

					case http.MethodPatch:
						appHandler.Reschedule(w, r)

					default:
						http.NotFound(w, r)
					}
				}),
			),
		),
	)

	return mux
}
