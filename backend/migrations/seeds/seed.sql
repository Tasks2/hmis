INSERT INTO users (email, password_hash, role)
VALUES
    ('patient@example.com', 'development-only', 'PATIENT'),
    ('doctor@example.com', 'development-only', 'PRACTITIONER');

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'John', 'Doe'
FROM users
WHERE email = 'patient@example.com';

INSERT INTO practitioners (user_id, first_name, last_name, specialty)
SELECT id, 'Jane', 'Smith', 'General Medicine'
FROM users
WHERE email = 'doctor@example.com';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT
    id,
    '2026-10-01',
    '09:00',
    '12:00'
FROM practitioners
WHERE first_name = 'Jane'
  AND last_name = 'Smith';