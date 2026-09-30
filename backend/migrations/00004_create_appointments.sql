-- +goose Up

CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL,
    practitioner_id UUID NOT NULL,
    appointment_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'BOOKED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT appointments_patient_fk
        FOREIGN KEY (patient_id)
        REFERENCES patients(id)
        ON DELETE RESTRICT,

    CONSTRAINT appointments_practitioner_fk
        FOREIGN KEY (practitioner_id)
        REFERENCES practitioners(id)
        ON DELETE RESTRICT,

    CONSTRAINT appointments_time_check
        CHECK (end_time > start_time),

    CONSTRAINT appointments_status_check
        CHECK (status IN ('BOOKED', 'CANCELLED'))
);

CREATE INDEX idx_appointments_patient
    ON appointments (patient_id);

CREATE INDEX idx_appointments_practitioner_date
    ON appointments (practitioner_id, appointment_date);

CREATE UNIQUE INDEX idx_active_appointment_slot
    ON appointments (
        practitioner_id,
        appointment_date,
        start_time
    )
    WHERE status = 'BOOKED';

-- +goose Down

DROP TABLE appointments;