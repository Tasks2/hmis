package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Tasks2/hmis/internal/middleware"
	"github.com/Tasks2/hmis/internal/response"
	"github.com/Tasks2/hmis/internal/services"
	"github.com/Tasks2/hmis/internal/validation"
)

// assertJSONError checks the status code and that the body is the
// standard {"error": {"code", "message"}} envelope with the expected code.
func assertJSONError(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()

	if rec.Code != wantStatus {
		t.Fatalf("expected status %d, got %d (body: %s)", wantStatus, rec.Code, rec.Body.String())
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var body response.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not a JSON error: %v (body: %s)", err, rec.Body.String())
	}

	if body.Error.Code != wantCode {
		t.Errorf("expected code %s, got %s", wantCode, body.Error.Code)
	}

	if body.Error.Message == "" {
		t.Error("expected a non-empty error message")
	}
}

func withUser(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, "user-123")
	return r.WithContext(ctx)
}

// These requests all fail validation before reaching a service, so the
// handlers can be built with nil services.
func TestHandlers_ValidationErrorsAreJSON(t *testing.T) {
	auth := NewAuthHandler(nil, nil)
	appt := NewAppointmentHandler(nil)
	prac := NewPractitionerHandler(nil)

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		method     string
		path       string
		body       string
		withUser   bool
		wantStatus int
		wantCode   string
	}{
		{"register: body not JSON", auth.Register, http.MethodPost, "/api/v1/auth/register", "{", false, 400, "INVALID_REQUEST"},
		{"register: email missing", auth.Register, http.MethodPost, "/api/v1/auth/register", `{"password":"x"}`, false, 400, "EMAIL_REQUIRED"},
		{"login: body not JSON", auth.Login, http.MethodPost, "/api/v1/auth/login", "nope", false, 400, "INVALID_REQUEST"},
		{"login: password missing", auth.Login, http.MethodPost, "/api/v1/auth/login", `{"email":"a@b.co"}`, false, 400, "PASSWORD_REQUIRED"},

		{"create: no user", appt.Create, http.MethodPost, "/api/v1/appointments", `{}`, false, 401, "UNAUTHORIZED"},
		{"create: body not JSON", appt.Create, http.MethodPost, "/api/v1/appointments", "{", true, 400, "INVALID_REQUEST"},
		{"create: practitioner missing", appt.Create, http.MethodPost, "/api/v1/appointments", `{"appointment_date":"2030-01-01","start_time":"09:00"}`, true, 400, "PRACTITIONER_ID_REQUIRED"},
		{"create: date missing", appt.Create, http.MethodPost, "/api/v1/appointments", `{"practitioner_id":"p1","start_time":"09:00"}`, true, 400, "DATE_REQUIRED"},
		{"create: time missing", appt.Create, http.MethodPost, "/api/v1/appointments", `{"practitioner_id":"p1","appointment_date":"2030-01-01"}`, true, 400, "TIME_REQUIRED"},
		{"create: bad date", appt.Create, http.MethodPost, "/api/v1/appointments", `{"practitioner_id":"p1","appointment_date":"01/01/2030","start_time":"09:00"}`, true, 400, "INVALID_DATE"},

		{"list: no user", appt.List, http.MethodGet, "/api/v1/appointments", "", false, 401, "UNAUTHORIZED"},
		{"cancel: no user", appt.Cancel, http.MethodDelete, "/api/v1/appointments/a1", "", false, 401, "UNAUTHORIZED"},
		{"cancel: id missing", appt.Cancel, http.MethodDelete, "/api/v1/appointments/", "", true, 400, "APPOINTMENT_ID_REQUIRED"},

		{"reschedule: no user", appt.Reschedule, http.MethodPatch, "/api/v1/appointments/a1", `{}`, false, 401, "UNAUTHORIZED"},
		{"reschedule: body not JSON", appt.Reschedule, http.MethodPatch, "/api/v1/appointments/a1", "{", true, 400, "INVALID_REQUEST"},
		{"reschedule: time missing", appt.Reschedule, http.MethodPatch, "/api/v1/appointments/a1", `{"appointment_date":"2030-01-01"}`, true, 400, "TIME_REQUIRED"},
		{"reschedule: bad date", appt.Reschedule, http.MethodPatch, "/api/v1/appointments/a1", `{"appointment_date":"tomorrow","start_time":"09:00"}`, true, 400, "INVALID_DATE"},

		{"availability: unknown path", prac.Availability, http.MethodGet, "/api/v1/practitioners/p1/other", "", false, 404, "NOT_FOUND"},
		{"availability: wrong method", prac.Availability, http.MethodPost, "/api/v1/practitioners/p1/availability?date=2030-01-01", "", false, 405, "METHOD_NOT_ALLOWED"},
		{"availability: date missing", prac.Availability, http.MethodGet, "/api/v1/practitioners/p1/availability", "", false, 400, "DATE_REQUIRED"},
		{"availability: bad date", prac.Availability, http.MethodGet, "/api/v1/practitioners/p1/availability?date=2030-13-01", "", false, 400, "INVALID_DATE"},
		{"practitioners: wrong method", prac.List, http.MethodPost, "/api/v1/practitioners", "", false, 405, "METHOD_NOT_ALLOWED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.withUser {
				req = withUser(req)
			}
			rec := httptest.NewRecorder()

			tt.handler(rec, req)

			assertJSONError(t, rec, tt.wantStatus, tt.wantCode)
		})
	}
}

func TestWriteAppointmentError(t *testing.T) {
	tests := []struct {
		err        error
		wantStatus int
		wantCode   string
	}{
		{services.ErrPatientNotFound, 404, "PATIENT_NOT_FOUND"},
		{services.ErrAppointmentNotFound, 404, "APPOINTMENT_NOT_FOUND"},
		{services.ErrPractitionerNotFound, 404, "PRACTITIONER_NOT_FOUND"},
		{services.ErrSlotUnavailable, 409, "SLOT_UNAVAILABLE"},
		{services.ErrPastAppointment, 400, "PAST_APPOINTMENT"},
		{services.ErrScheduleUnavailable, 400, "SCHEDULE_UNAVAILABLE"},
		// Regression: these two used to fall through to a 500.
		{validation.ErrInvalidTime, 400, "INVALID_TIME"},
		{validation.ErrInvalidSlot, 400, "INVALID_SLOT"},
		// Wrapped errors are still recognised.
		{fmt.Errorf("booking: %w", services.ErrSlotUnavailable), 409, "SLOT_UNAVAILABLE"},
		{errors.New("database exploded"), 500, "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.wantCode, func(t *testing.T) {
			rec := httptest.NewRecorder()

			writeAppointmentError(rec, tt.err)

			assertJSONError(t, rec, tt.wantStatus, tt.wantCode)
		})
	}

	t.Run("internal errors don't leak details", func(t *testing.T) {
		rec := httptest.NewRecorder()

		writeAppointmentError(rec, errors.New("pq: password authentication failed for user hmis_app"))

		if strings.Contains(rec.Body.String(), "hmis_app") {
			t.Errorf("internal error details leaked into the response: %s", rec.Body.String())
		}
	})
}
