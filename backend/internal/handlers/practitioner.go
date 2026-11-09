package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/Tasks2/hmis/internal/response"
	"github.com/Tasks2/hmis/internal/services"
	"github.com/Tasks2/hmis/internal/validation"
)

type PractitionerHandler struct {
	service *services.PractitionerService
}

func NewPractitionerHandler(
	service *services.PractitionerService,
) *PractitionerHandler {
	return &PractitionerHandler{service: service}
}

// List returns every practitioner.
func (h *PractitionerHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		MethodNotAllowed(w, http.MethodGet)
		return
	}

	practitioners, err := h.service.GetPractitioners(r.Context())
	if err != nil {
		log.Printf("list practitioners error: %v", err)
		response.JSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Something went wrong on our side. Try again.")
		return
	}

	if practitioners == nil {
		practitioners = []services.Practitioner{}
	}

	response.JSON(w, http.StatusOK, practitioners)
}

// Availability returns a practitioner's free slots on one date:
// GET /api/v1/practitioners/{id}/availability?date=YYYY-MM-DD
func (h *PractitionerHandler) Availability(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/practitioners/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] != "availability" {
		response.JSONError(w, http.StatusNotFound, "NOT_FOUND",
			"endpoint not found")
		return
	}

	if r.Method != http.MethodGet {
		MethodNotAllowed(w, http.MethodGet)
		return
	}

	practitionerID := parts[0]
	dateString := r.URL.Query().Get("date")

	if !validation.Required(dateString) {
		response.JSONError(w, http.StatusBadRequest, "DATE_REQUIRED",
			"date is required")
		return
	}

	date, err := validation.Date(dateString)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "INVALID_DATE",
			"date must use YYYY-MM-DD format")
		return
	}

	slots, err := h.service.GetAvailability(r.Context(), practitionerID, date)
	if err != nil {
		log.Printf("availability error: %v", err)
		response.JSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR",
			"Something went wrong on our side. Try again.")
		return
	}

	if slots == nil {
		slots = []services.AvailableSlot{}
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"practitioner_id": practitionerID,
		"date":            dateString,
		"slots":           slots,
	})
}

// MethodNotAllowed writes a JSON 405 with the Allow header set.
func MethodNotAllowed(w http.ResponseWriter, allowed ...string) {
	w.Header().Set("Allow", strings.Join(allowed, ", "))
	response.JSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED",
		"method not allowed")
}
