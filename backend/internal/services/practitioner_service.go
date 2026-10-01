package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Practitioner struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Specialty string `json:"specialty"`
}

type AvailableSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type PractitionerService struct {
	db *pgxpool.Pool
}

func NewPractitionerService(db *pgxpool.Pool) *PractitionerService {
	return &PractitionerService{db: db}
}

func (s *PractitionerService) GetPractitioners(
	ctx context.Context,
) ([]Practitioner, error) {

	rows, err := s.db.Query(ctx, `
		SELECT id, first_name, last_name, specialty
		FROM practitioners
		ORDER BY last_name, first_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var practitioners []Practitioner

	for rows.Next() {
		var p Practitioner

		if err := rows.Scan(
			&p.ID,
			&p.FirstName,
			&p.LastName,
			&p.Specialty,
		); err != nil {
			return nil, err
		}

		practitioners = append(practitioners, p)
	}

	return practitioners, rows.Err()
}

func (s *PractitionerService) GetAvailability(
	ctx context.Context,
	practitionerID string,
	date time.Time,
) ([]AvailableSlot, error) {

	rows, err := s.db.Query(ctx, `
		SELECT start_time, end_time, slot_duration_minutes
		FROM schedules
		WHERE practitioner_id = $1
		  AND schedule_date = $2
		ORDER BY start_time
	`, practitionerID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []AvailableSlot

	for rows.Next() {
		var (
			startTime time.Time
			endTime   time.Time
			duration  int
		)

		if err := rows.Scan(
			&startTime,
			&endTime,
			&duration,
		); err != nil {
			return nil, err
		}

		for current := startTime; current.Before(endTime); current = current.Add(time.Duration(duration) * time.Minute) {
			slotEnd := current.Add(time.Duration(duration) * time.Minute)

			if slotEnd.After(endTime) {
				break
			}

			var booked bool

			err := s.db.QueryRow(
				ctx,
				`SELECT EXISTS (
					SELECT 1
					FROM appointments
					WHERE practitioner_id = $1
					  AND appointment_date = $2
					  AND start_time = $3
					  AND status = 'BOOKED'
				)`,
				practitionerID,
				date,
				current.Format("15:04:05"),
			).Scan(&booked)

			if err != nil {
				return nil, err
			}

			if !booked {
				slots = append(slots, AvailableSlot{
					StartTime: current.Format("15:04"),
					EndTime:   slotEnd.Format("15:04"),
				})
			}
		}
	}

	return slots, rows.Err()
}
