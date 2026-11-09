package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Tasks2/hmis/internal/middleware"
	"github.com/Tasks2/hmis/internal/response"
	"github.com/Tasks2/hmis/internal/services"
	"github.com/Tasks2/hmis/internal/validation"
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

func getUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)

	if !ok || userID == "" {
		return "", false
	}

	return userID, true
}

// writeAppointmentError turns an error from the appointment service into a
// JSON error response. Every appointment handler uses it, so the same
// error always gets the same status code and message.
func writeAppointmentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrPatientNotFound):
		response.JSONError(w, http.StatusNotFound, "PATIENT_NOT_FOUND",
			"No patient profile is linked to this account.")

	case errors.Is(err, services.ErrAppointmentNotFound):
		response.JSONError(w, http.StatusNotFound, "APPOINTMENT_NOT_FOUND",
			"Appointment not found, or it is no longer booked.")

	case errors.Is(err, services.ErrPractitionerNotFound):
		response.JSONError(w, http.StatusNotFound, "PRACTITIONER_NOT_FOUND",
			"Practitioner not found.")

	case errors.Is(err, services.ErrSlotUnavailable):
		response.JSONError(w, http.StatusConflict, "SLOT_UNAVAILABLE",
			"This slot has just been booked. Choose another time.")

	case errors.Is(err, services.ErrPastAppointment):
		response.JSONError(w, http.StatusBadRequest, "PAST_APPOINTMENT",
			"Appointment date can't be in the past.")

	case errors.Is(err, services.ErrScheduleUnavailable):
		response.JSONError(w, http.StatusBadRequest, "SCHEDULE_UNAVAILABLE",
			"The practitioner isn't available at that time.")

	case errors.Is(err, validation.ErrInvalidTime):
		response.JSONError(w, http.StatusBadRequest, "INVALID_TIME",
			"start_time must use HH:MM format, for example 09:30.")

	case errors.Is(err, validation.ErrInvalidSlot):
		response.JSONError(w, http.StatusBadRequest, "INVALID_SLOT",
			"Appointments start on the hour or half hour, for example 09:00 or 09:30.")

	default:
		log.Printf("appointment error: %v", err)
		response.JSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Something went wrong on our side. Try again.")
	}
}

// requireUserID writes a 401 and returns false when the request has no
// authenticated user. The JWT middleware normally guarantees one.
func requireUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := getUserID(r)
	if !ok {
		response.JSONError(w, http.StatusUnauthorized, "UNAUTHORIZED",
			"authentication is required")
	}

	return userID, ok
}

// parseSlot validates the date and start time shared by create and
// reschedule requests.
func parseSlot(w http.ResponseWriter, date, startTime string) (time.Time, bool) {
	if !validation.Required(date) {
		response.JSONError(w, http.StatusBadRequest, "DATE_REQUIRED",
			"appointment_date is required")
		return time.Time{}, false
	}

	if !validation.Required(startTime) {
		response.JSONError(w, http.StatusBadRequest, "TIME_REQUIRED",
			"start_time is required")
		return time.Time{}, false
	}

	parsed, err := validation.Date(date)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "INVALID_DATE",
			"appointment_date must use YYYY-MM-DD format")
		return time.Time{}, false
	}

	return parsed, true
}

func appointmentIDFromPath(r *http.Request) string {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/appointments/")
	return strings.TrimSuffix(id, "/reschedule")
}

// Create books an appointment.
func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	var req CreateAppointmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "INVALID_REQUEST",
			"invalid request body")
		return
	}

	if !validation.Required(req.PractitionerID) {
		response.JSONError(w, http.StatusBadRequest, "PRACTITIONER_ID_REQUIRED",
			"practitioner_id is required")
		return
	}

	date, ok := parseSlot(w, req.AppointmentDate, req.StartTime)
	if !ok {
		return
	}

	err := h.service.Create(
		r.Context(),
		userID,
		req.PractitionerID,
		date,
		req.StartTime,
	)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "appointment booked successfully",
	})
}

// List returns the patient's appointments.
func (h *AppointmentHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	appointments, err := h.service.GetPatientAppointments(r.Context(), userID)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}

	// Return [] instead of null when the patient has no appointments.
	if appointments == nil {
		appointments = []services.Appointment{}
	}

	response.JSON(w, http.StatusOK, appointments)
}

// Cancel cancels one of the patient's booked appointments.
func (h *AppointmentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	appointmentID := appointmentIDFromPath(r)
	if appointmentID == "" {
		response.JSONError(w, http.StatusBadRequest, "APPOINTMENT_ID_REQUIRED",
			"appointment id is required")
		return
	}

	if err := h.service.Cancel(r.Context(), userID, appointmentID); err != nil {
		writeAppointmentError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "appointment cancelled successfully",
	})
}

// Reschedule moves one of the patient's booked appointments to a new
// date and time with the same practitioner.
func (h *AppointmentHandler) Reschedule(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	appointmentID := appointmentIDFromPath(r)
	if appointmentID == "" {
		response.JSONError(w, http.StatusBadRequest, "APPOINTMENT_ID_REQUIRED",
			"appointment id is required")
		return
	}

	var req RescheduleAppointmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "INVALID_REQUEST",
			"invalid request body")
		return
	}

	date, ok := parseSlot(w, req.AppointmentDate, req.StartTime)
	if !ok {
		return
	}

	err := h.service.Reschedule(
		r.Context(),
		userID,
		appointmentID,
		date,
		req.StartTime,
	)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "appointment rescheduled successfully",
	})
}
