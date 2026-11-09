-- Development data only. Every seeded user's password is: Password123!
-- (password_hash below is the bcrypt hash of that password.)

-- ============================================================
-- HMIS TEST DATA
-- 5 practitioners
-- 10 patients
-- Future schedules for testing availability and appointments
-- ============================================================


-- ============================================================
-- 1. PRACTITIONER USERS
-- ============================================================

INSERT INTO users (email, password_hash, role)
VALUES
    ('doctor1@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PRACTITIONER'),
    ('doctor2@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PRACTITIONER'),
    ('doctor3@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PRACTITIONER'),
    ('doctor4@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PRACTITIONER'),
    ('doctor5@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PRACTITIONER');


-- ============================================================
-- 2. PRACTITIONER PROFILES
-- ============================================================

INSERT INTO practitioners (
    user_id,
    first_name,
    last_name,
    specialty
)
SELECT id, 'Jane', 'Smith', 'General Medicine'
FROM users
WHERE email = 'doctor1@example.com';

INSERT INTO practitioners (
    user_id,
    first_name,
    last_name,
    specialty
)
SELECT id, 'John', 'Kamau', 'Cardiology'
FROM users
WHERE email = 'doctor2@example.com';

INSERT INTO practitioners (
    user_id,
    first_name,
    last_name,
    specialty
)
SELECT id, 'Mary', 'Achieng', 'Pediatrics'
FROM users
WHERE email = 'doctor3@example.com';

INSERT INTO practitioners (
    user_id,
    first_name,
    last_name,
    specialty
)
SELECT id, 'David', 'Mwangi', 'Dermatology'
FROM users
WHERE email = 'doctor4@example.com';

INSERT INTO practitioners (
    user_id,
    first_name,
    last_name,
    specialty
)
SELECT id, 'Grace', 'Wanjiku', 'Family Medicine'
FROM users
WHERE email = 'doctor5@example.com';


-- ============================================================
-- 3. PATIENT USERS
-- ============================================================

INSERT INTO users (email, password_hash, role)
VALUES
    ('patient1@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT'),
    ('patient2@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT'),
    ('patient3@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT'),
    ('patient4@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT'),
    ('patient5@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT'),
    ('patient6@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT'),
    ('patient7@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT'),
    ('patient8@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT'),
    ('patient9@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT'),
    ('patient10@example.com', '$2a$10$ac47aIkRzOKx19TxZe7.TelCHGk4KhtZx.PYP69WRbunEW6zGhxwa', 'PATIENT');


-- ============================================================
-- 4. PATIENT PROFILES
-- ============================================================

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'One'
FROM users
WHERE email = 'patient1@example.com';

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'Two'
FROM users
WHERE email = 'patient2@example.com';

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'Three'
FROM users
WHERE email = 'patient3@example.com';

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'Four'
FROM users
WHERE email = 'patient4@example.com';

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'Five'
FROM users
WHERE email = 'patient5@example.com';

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'Six'
FROM users
WHERE email = 'patient6@example.com';

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'Seven'
FROM users
WHERE email = 'patient7@example.com';

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'Eight'
FROM users
WHERE email = 'patient8@example.com';

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'Nine'
FROM users
WHERE email = 'patient9@example.com';

INSERT INTO patients (user_id, first_name, last_name)
SELECT id, 'Patient', 'Ten'
FROM users
WHERE email = 'patient10@example.com';


-- ============================================================
-- 5. PRACTITIONER SCHEDULES
--
-- V1 behavior:
-- A practitioner is considered available on a date only when
-- that date exists in schedules.
--
-- We intentionally seed several future dates so the frontend
-- availability and booking flow can be tested.
-- ============================================================


-- Jane Smith
INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-05', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Jane' AND last_name = 'Smith';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-06', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Jane' AND last_name = 'Smith';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-07', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Jane' AND last_name = 'Smith';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-08', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Jane' AND last_name = 'Smith';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-09', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Jane' AND last_name = 'Smith';


-- John Kamau
INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-05', '10:00', '13:00'
FROM practitioners
WHERE first_name = 'John' AND last_name = 'Kamau';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-06', '10:00', '13:00'
FROM practitioners
WHERE first_name = 'John' AND last_name = 'Kamau';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-07', '10:00', '13:00'
FROM practitioners
WHERE first_name = 'John' AND last_name = 'Kamau';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-08', '10:00', '13:00'
FROM practitioners
WHERE first_name = 'John' AND last_name = 'Kamau';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-09', '10:00', '13:00'
FROM practitioners
WHERE first_name = 'John' AND last_name = 'Kamau';


-- Mary Achieng
INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-05', '08:00', '11:00'
FROM practitioners
WHERE first_name = 'Mary' AND last_name = 'Achieng';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-06', '08:00', '11:00'
FROM practitioners
WHERE first_name = 'Mary' AND last_name = 'Achieng';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-07', '08:00', '11:00'
FROM practitioners
WHERE first_name = 'Mary' AND last_name = 'Achieng';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-08', '08:00', '11:00'
FROM practitioners
WHERE first_name = 'Mary' AND last_name = 'Achieng';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-09', '08:00', '11:00'
FROM practitioners
WHERE first_name = 'Mary' AND last_name = 'Achieng';


-- David Mwangi
INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-05', '13:00', '16:00'
FROM practitioners
WHERE first_name = 'David' AND last_name = 'Mwangi';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-06', '13:00', '16:00'
FROM practitioners
WHERE first_name = 'David' AND last_name = 'Mwangi';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-07', '13:00', '16:00'
FROM practitioners
WHERE first_name = 'David' AND last_name = 'Mwangi';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-08', '13:00', '16:00'
FROM practitioners
WHERE first_name = 'David' AND last_name = 'Mwangi';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-09', '13:00', '16:00'
FROM practitioners
WHERE first_name = 'David' AND last_name = 'Mwangi';


-- Grace Wanjiku
INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-05', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Grace' AND last_name = 'Wanjiku';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-06', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Grace' AND last_name = 'Wanjiku';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-07', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Grace' AND last_name = 'Wanjiku';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-08', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Grace' AND last_name = 'Wanjiku';

INSERT INTO schedules (
    practitioner_id,
    schedule_date,
    start_time,
    end_time
)
SELECT id, '2026-10-09', '09:00', '12:00'
FROM practitioners
WHERE first_name = 'Grace' AND last_name = 'Wanjiku';