package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Tasks2/hmis/internal/services"
)

type PractitionerHandler struct {
	service *services.PractitionerService
}

func NewPractitionerHandler(
	service *services.PractitionerService,
) *PractitionerHandler {
	return &PractitionerHandler{service: service}
}

func (h *PractitionerHandler) List(w http.ResponseWriter, r *http.Request) {
	practitioners, err := h.service.GetPractitioners(r.Context())

	if err != nil {
		http.Error(w, "could not retrieve practitioners", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(practitioners)
}

func (h *PractitionerHandler) Availability(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/practitioners/",
	)

	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "availability" {
		http.NotFound(w, r)
		return
	}

	practitionerID := parts[0]

	dateString := r.URL.Query().Get("date")

	if dateString == "" {
		http.Error(w, "date is required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	slots, err := h.service.GetAvailability(
		r.Context(),
		practitionerID,
		date,
	)

	if err != nil {
		http.Error(w, "could not retrieve availability", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"practitioner_id": practitionerID,
		"date":            dateString,
		"slots":           slots,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
