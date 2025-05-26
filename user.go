package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserPostgresStore struct {
	db *PostgresDB
}

type UserStore interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

func (s *UserPostgresStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	row := s.db.pool.QueryRow(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if err := row.Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (s *UserPostgresStore) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := s.db.pool.QueryRow(ctx, query, id)

	var user User

	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserPostgresStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := s.db.pool.QueryRow(ctx, query, email)

	var user User

	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserPostgresStore) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := s.db.pool.Exec(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

func (s *UserPostgresStore) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := s.db.pool.Exec(ctx, query, id)

	return err
}

type UserService struct {
	userStore UserStore
}

func NewUserService(userStore UserStore) *UserService {
	return &UserService{userStore: userStore}
}

func validateUser(user *User) error {
	if user.Password == "" {
		return errors.New("password is required")
	}

	if user.FirstName == "" {
		return errors.New("first name is required")
	}

	if user.LastName == "" {
		return errors.New("last name is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

func (s *UserService) Create(ctx context.Context, user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := validateUser(user); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return errors.New("an internal error occurred")
	}

	user.Password = string(hashedPassword)

	existingUser, err := s.userStore.GetByEmail(ctx, user.Email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("failed to check for existing user", "error", err)
		return errors.New("an internal error occurred")
	}

	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	err = s.userStore.Create(ctx, user)

	if err != nil {
		slog.Error("failed to create user", "error", err)
		return errors.New("an internal error occurred")
	}

	return nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
	user, err := s.userStore.GetByID(ctx, id)

	if err != nil {
		slog.Error("failed to get user", "error", err)
		return nil, errors.New("an internal error occurred")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.userStore.GetByEmail(ctx, email)

	if err != nil {
		slog.Error("failed to get user", "error", err)
		return nil, errors.New("an internal error occurred")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (s *UserService) Update(ctx context.Context, id string, request *UpdateUserRequest) error {
	existingUser, err := s.userStore.GetByID(ctx, id)

	if err != nil {
		slog.Error("failed to get user", "error", err)
		return errors.New("an internal error occurred")
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	if request.FirstName == "" {
		existingUser.FirstName = request.FirstName
	}

	if request.LastName == "" {
		existingUser.LastName = request.LastName
	}

	existingUser.UpdatedAt = time.Now()

	err = s.userStore.Update(ctx, existingUser)

	if err != nil {
		slog.Error("failed to update user", "error", err)
		return errors.New("an internal error occurred")
	}

	return err
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

func (s *UserService) UpdatePassword(ctx context.Context, id string, request *UpdatePasswordRequest) error {
	existingUser, err := s.userStore.GetByID(ctx, id)

	if err != nil {
		slog.Error("failed to get user", "error", err)
		return errors.New("an internal error occurred")
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	if request.Password == "" {
		return errors.New("password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return errors.New("an internal error occurred")
	}

	existingUser.Password = string(hashedPassword)
	err = s.userStore.Update(ctx, existingUser)

	if err != nil {
		slog.Error("failed to update user", "error", err)
		return errors.New("an internal error occurred")
	}

	return nil
}

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

func (s *UserService) UpdateEmail(ctx context.Context, id string, request *UpdateEmailRequest) error {
	existingUser, err := s.userStore.GetByID(ctx, id)

	if err != nil {
		slog.Error("failed to get user", "error", err)
		return errors.New("an internal error occurred")
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	if request.Email == "" {
		return errors.New("email is required")
	}

	existingEmailUser, err := s.userStore.GetByEmail(ctx, request.Email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("failed to check for existing user", "error", err)
		return errors.New("an internal error occurred")
	}

	if existingEmailUser != nil {
		return errors.New("user with this email already exists")
	}

	existingUser.Email = request.Email
	existingUser.UpdatedAt = time.Now()

	err = s.userStore.Update(ctx, existingUser)

	if err != nil {
		slog.Error("failed to update user", "error", err)
		return errors.New("an internal error occurred")
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	user, err := s.userStore.GetByID(ctx, id)

	if err != nil {
		slog.Error("failed to get user", "error", err)
		return errors.New("an internal error occurred")
	}

	if user == nil {
		return errors.New("user not found")
	}

	err = s.userStore.Delete(ctx, id)

	if err != nil {
		slog.Error("failed to delete user", "error", err)
		return errors.New("an internal error occurred")
	}

	return nil
}
