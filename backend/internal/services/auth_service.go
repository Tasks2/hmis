package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{
		db: db,
	}
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

func (s *AuthService) CreatePatient(
	ctx context.Context,
	email string,
	password string,
	firstName string,
	lastName string,
) error {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var userID string

	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash, role)
		 VALUES ($1, $2, 'PATIENT')
		 RETURNING id`,
		email,
		passwordHash,
	).Scan(&userID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO patients (user_id, first_name, last_name)
		 VALUES ($1, $2, $3)`,
		userID,
		firstName,
		lastName,
	)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, string, error) {

	var (
		userID       string
		passwordHash string
		role         string
	)

	err := s.db.QueryRow(
		ctx,
		`SELECT id, password_hash, role
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(
		&userID,
		&passwordHash,
		&role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", ErrInvalidCredentials
	}

	if err != nil {
		return "", "", err
	}

	if err := CheckPassword(password, passwordHash); err != nil {
		return "", "", ErrInvalidCredentials
	}

	return userID, role, nil
}
