package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSlotUnavailable     = errors.New("appointment slot is unavailable")
	ErrPatientNotFound     = errors.New("patient not found")
	ErrAppointmentNotFound = errors.New("appointment not found")
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
