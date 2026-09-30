-- +goose Up

CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    practitioner_id UUID NOT NULL,
    schedule_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    slot_duration_minutes INTEGER NOT NULL DEFAULT 30,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT schedules_practitioner_fk
        FOREIGN KEY (practitioner_id)
        REFERENCES practitioners(id)
        ON DELETE CASCADE,

    CONSTRAINT schedules_time_check
        CHECK (end_time > start_time),

    CONSTRAINT schedules_slot_duration_check
        CHECK (slot_duration_minutes = 30)
);

CREATE INDEX idx_schedules_practitioner_date
    ON schedules (practitioner_id, schedule_date);

-- +goose Down

DROP TABLE schedules;