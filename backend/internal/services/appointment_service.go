package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Tasks2/hmis/internal/validation"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSlotUnavailable      = errors.New("appointment slot is unavailable")
	ErrPatientNotFound      = errors.New("patient not found")
	ErrAppointmentNotFound  = errors.New("appointment not found")
	ErrPractitionerNotFound = errors.New("practitioner not found")
	ErrScheduleUnavailable  = errors.New("schedule unavailable")
	ErrPastAppointment      = errors.New("appointment date is in the past")
)

type Appointment struct {
	ID              string `json:"id"`
	PractitionerID  string `json:"practitioner_id"`
	AppointmentDate string `json:"appointment_date"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	Status          string `json:"status"`
}

type AppointmentService struct {
	db *pgxpool.Pool
}

func NewAppointmentService(db *pgxpool.Pool) *AppointmentService {
	return &AppointmentService{db: db}
}

//Get Patient Id

func (s *AppointmentService) getPatientID(
	ctx context.Context,
	userID string,
) (string, error) {

	var patientID string

	err := s.db.QueryRow(
		ctx,
		`SELECT id
		 FROM patients
		 WHERE user_id = $1`,
		userID,
	).Scan(&patientID)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrPatientNotFound
	}

	if err != nil {
		return "", err
	}

	return patientID, nil
}

//Create appointment

func (s *AppointmentService) Create(
	ctx context.Context,
	userID string,
	practitionerID string,
	appointmentDate time.Time,
	startTime string,
) error {

	patientID, err := s.getPatientID(ctx, userID)
	if err != nil {
		return err
	}

	var endTime string

	err = s.db.QueryRow(
		ctx,
		`SELECT ($1::time + INTERVAL '30 minutes')::time`,
		startTime,
	).Scan(&endTime)

	if err != nil {
		return fmt.Errorf("calculate appointment end time: %w", err)
	}

	if err := s.validateSlot(
		ctx,
		practitionerID,
		appointmentDate,
		startTime,
	); err != nil {
		return err
	}

	_, err = s.db.Exec(
		ctx,
		`INSERT INTO appointments (
			patient_id,
			practitioner_id,
			appointment_date,
			start_time,
			end_time,
			status
		)
		VALUES ($1, $2, $3, $4, $5, 'BOOKED')`,
		patientID,
		practitionerID,
		appointmentDate,
		startTime,
		endTime,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrSlotUnavailable
		}

		return err
	}

	return nil
}

// Cancel appointment
func (s *AppointmentService) Cancel(
	ctx context.Context,
	userID string,
	appointmentID string,
) error {

	patientID, err := s.getPatientID(ctx, userID)
	if err != nil {
		return err
	}

	result, err := s.db.Exec(
		ctx,
		`UPDATE appointments
		 SET status = 'CANCELLED',
		     updated_at = NOW()
		 WHERE id = $1
		   AND patient_id = $2
		   AND status = 'BOOKED'`,
		appointmentID,
		patientID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrAppointmentNotFound
	}

	return nil
}

// Reschedule Appointment
func (s *AppointmentService) Reschedule(
	ctx context.Context,
	userID string,
	appointmentID string,
	appointmentDate time.Time,
	startTime string,
) error {

	patientID, err := s.getPatientID(ctx, userID)
	if err != nil {
		return err
	}

	var endTime string

	err = s.db.QueryRow(
		ctx,
		`SELECT ($1::time + INTERVAL '30 minutes')::time`,
		startTime,
	).Scan(&endTime)

	if err != nil {
		return fmt.Errorf("calculate appointment end time: %w", err)
	}

	result, err := s.db.Exec(
		ctx,
		`UPDATE appointments
		 SET appointment_date = $1,
		     start_time = $2,
		     end_time = $3,
		     updated_at = NOW()
		 WHERE id = $4
		   AND patient_id = $5
		   AND status = 'BOOKED'`,
		appointmentDate,
		startTime,
		endTime,
		appointmentID,
		patientID,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrSlotUnavailable
		}

		return err
	}

	if result.RowsAffected() == 0 {
		return ErrAppointmentNotFound
	}

	return nil
}

// Get Appointments
func (s *AppointmentService) GetPatientAppointments(
	ctx context.Context,
	userID string,
) ([]Appointment, error) {

	patientID, err := s.getPatientID(ctx, userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(
		ctx,
		`SELECT
			id,
			practitioner_id,
			appointment_date,
			start_time,
			end_time,
			status
		 FROM appointments
		 WHERE patient_id = $1
		 ORDER BY appointment_date, start_time`,
		patientID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment

	for rows.Next() {
		var appointment Appointment

		var (
			date      time.Time
			startTime time.Time
			endTime   time.Time
		)

		err := rows.Scan(
			&appointment.ID,
			&appointment.PractitionerID,
			&date,
			&startTime,
			&endTime,
			&appointment.Status,
		)

		if err != nil {
			return nil, err
		}

		appointment.AppointmentDate = date.Format("2006-01-02")
		appointment.StartTime = startTime.Format("15:04")
		appointment.EndTime = endTime.Format("15:04")

		appointments = append(appointments, appointment)
	}

	return appointments, rows.Err()
}

// Schedule Validation Helper
func (s *AppointmentService) validateSlot(
	ctx context.Context,
	practitionerID string,
	appointmentDate time.Time,
	startTime string,
) error {
	var practitionerExists bool

	err := s.db.QueryRow(
		ctx,
		`SELECT EXISTS(
			SELECT 1
			FROM practitioners
			WHERE id = $1
		)`,
		practitionerID,
	).Scan(&practitionerExists)

	if err != nil {
		return fmt.Errorf("check practitioner: %w", err)
	}

	if !practitionerExists {
		return ErrPractitionerNotFound
	}

	today := time.Now()

	requestedDate := appointmentDate.Format("2006-01-02")
	currentDate := today.Format("2006-01-02")

	if requestedDate < currentDate {
		return ErrPastAppointment
	}

	requestedStart, err := time.Parse("15:04", startTime)
	if err != nil {
		return validation.ErrInvalidTime
	}

	if !validation.Is30MinuteSlot(startTime) {
		return validation.ErrInvalidSlot
	}

	requestedEnd := requestedStart.Add(30 * time.Minute)

	rows, err := s.db.Query(
		ctx,
		`SELECT start_time, end_time, slot_duration_minutes
		 FROM schedules
		 WHERE practitioner_id = $1
		   AND schedule_date = $2`,
		practitionerID,
		appointmentDate,
	)
	if err != nil {
		return fmt.Errorf("query practitioner schedule: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var scheduleStart time.Time
		var scheduleEnd time.Time
		var duration int

		if err := rows.Scan(
			&scheduleStart,
			&scheduleEnd,
			&duration,
		); err != nil {
			return fmt.Errorf("scan practitioner schedule: %w", err)
		}

		// V1 uses fixed 30-minute appointment slots.
		if duration != 30 {
			continue
		}

		scheduleStartTime := time.Date(
			0, 1, 1,
			scheduleStart.Hour(),
			scheduleStart.Minute(),
			0,
			0,
			time.UTC,
		)

		scheduleEndTime := time.Date(
			0, 1, 1,
			scheduleEnd.Hour(),
			scheduleEnd.Minute(),
			0,
			0,
			time.UTC,
		)

		requestedStartTime := time.Date(
			0, 1, 1,
			requestedStart.Hour(),
			requestedStart.Minute(),
			0,
			0,
			time.UTC,
		)

		requestedEndTime := time.Date(
			0, 1, 1,
			requestedEnd.Hour(),
			requestedEnd.Minute(),
			0,
			0,
			time.UTC,
		)

		if !requestedStartTime.Before(scheduleStartTime) &&
			!requestedEndTime.After(scheduleEndTime) {
			return nil
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate practitioner schedule: %w", err)
	}

	return ErrScheduleUnavailable
}
