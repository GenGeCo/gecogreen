package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gecogreen/backend/internal/models"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, phone, roles, status, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7::user_role[], $8::user_status, $9, $10, $11)
	`

	user.ID = uuid.New()
	user.Status = models.UserStatusActive // Auto-activate for now
	user.EmailVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if len(user.Roles) == 0 {
		user.Roles = []models.UserRole{models.RoleBuyer}
	}

	// Convert roles to string slice for PostgreSQL
	rolesStr := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		rolesStr[i] = string(r)
	}

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, rolesStr, string(user.Status), user.EmailVerified, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil && err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
		return ErrEmailAlreadyExists
	}
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name,
		       COALESCE(phone, ''), COALESCE(city, ''), roles::text[], status::text,
		       email_verified, COALESCE(avatar_url, ''), stripe_customer_id, stripe_account_id,
		       last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`

	user := &models.User{}
	var roles []string
	var status string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.City, &roles, &status, &user.EmailVerified,
		&user.AvatarURL, &user.StripeCustomerID, &user.StripeAccountID, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Convert string slices to typed values
	for _, r := range roles {
		user.Roles = append(user.Roles, models.UserRole(r))
	}
	user.Status = models.UserStatus(status)

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name,
		       COALESCE(phone, ''), COALESCE(city, ''), roles::text[], status::text,
		       email_verified, COALESCE(avatar_url, ''), stripe_customer_id, stripe_account_id,
		       last_login_at, created_at, updated_at
		FROM users WHERE email = $1
	`

	user := &models.User{}
	var roles []string
	var status string
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.City, &roles, &status, &user.EmailVerified,
		&user.AvatarURL, &user.StripeCustomerID, &user.StripeAccountID, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Convert string slices to typed values
	for _, r := range roles {
		user.Roles = append(user.Roles, models.UserRole(r))
	}
	user.Status = models.UserStatus(status)

	return user, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $1, updated_at = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}
