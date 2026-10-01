package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Tasks2/hmis/internal/middleware"
	"github.com/Tasks2/hmis/internal/services"
)

type AppointmentHandler struct {
	service *services.AppointmentService
}

func NewAppointmentHandler(
	service *services.AppointmentService,
) *AppointmentHandler {
	return &AppointmentHandler{service: service}
}

type CreateAppointmentRequest struct {
	PractitionerID  string `json:"practitioner_id"`
	AppointmentDate string `json:"appointment_date"`
	StartTime       string `json:"start_time"`
}

type RescheduleAppointmentRequest struct {
	AppointmentDate string `json:"appointment_date"`
	StartTime       string `json:"start_time"`
}

//Create Appointments

func getUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)

	if !ok || userID == "" {
		return "", false
	}

	return userID, true
}

func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateAppointmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.PractitionerID == "" ||
		req.AppointmentDate == "" ||
		req.StartTime == "" {

		http.Error(
			w,
			"practitioner_id, appointment_date and start_time are required",
			http.StatusBadRequest,
		)
		return
	}

	date, err := time.Parse("2006-01-02", req.AppointmentDate)

	if err != nil {
		http.Error(w, "invalid appointment date", http.StatusBadRequest)
		return
	}

	err = h.service.Create(
		r.Context(),
		userID,
		req.PractitionerID,
		date,
		req.StartTime,
	)

	if err != nil {
		switch {
		case errors.Is(err, services.ErrPatientNotFound):
			http.Error(w, "patient not found", http.StatusNotFound)

		case errors.Is(err, services.ErrSlotUnavailable):
			http.Error(w, "appointment slot unavailable", http.StatusConflict)

		default:
			http.Error(w, "could not create appointment", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "appointment booked successfully",
	})
}

// Get Appointments
func (h *AppointmentHandler) List(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	appointments, err := h.service.GetPatientAppointments(
		r.Context(),
		userID,
	)

	if err != nil {
		if errors.Is(err, services.ErrPatientNotFound) {
			http.Error(w, "patient not found", http.StatusNotFound)
			return
		}

		http.Error(
			w,
			"could not retrieve appointments",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(appointments)
}

// Cancel Appointments
func (h *AppointmentHandler) Cancel(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	appointmentID := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/appointments/",
	)

	if appointmentID == "" {
		http.Error(w, "appointment id is required", http.StatusBadRequest)
		return
	}

	err := h.service.Cancel(
		r.Context(),
		userID,
		appointmentID,
	)

	if err != nil {
		if errors.Is(err, services.ErrAppointmentNotFound) {
			http.Error(w, "appointment not found", http.StatusNotFound)
			return
		}

		http.Error(
			w,
			"could not cancel appointment",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "appointment cancelled successfully",
	})
}

// Reschedule appointment
func (h *AppointmentHandler) Reschedule(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	appointmentID := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/appointments/",
	)

	appointmentID = strings.TrimSuffix(
		appointmentID,
		"/reschedule",
	)

	if appointmentID == "" {
		http.Error(w, "appointment id is required", http.StatusBadRequest)
		return
	}

	var req RescheduleAppointmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(
		"2006-01-02",
		req.AppointmentDate,
	)

	if err != nil {
		http.Error(w, "invalid appointment date", http.StatusBadRequest)
		return
	}

	err = h.service.Reschedule(
		r.Context(),
		userID,
		appointmentID,
		date,
		req.StartTime,
	)

	if err != nil {
		switch {
		case errors.Is(err, services.ErrAppointmentNotFound):
			http.Error(w, "appointment not found", http.StatusNotFound)

		case errors.Is(err, services.ErrSlotUnavailable):
			http.Error(w, "appointment slot unavailable", http.StatusConflict)

		default:
			http.Error(
				w,
				"could not reschedule appointment",
				http.StatusInternalServerError,
			)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "appointment rescheduled successfully",
	})
}
