package user

import (
	"auth-haven/internal/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists for this tenant")
	ErrNothingToUpdate    = errors.New("nothing to update")
	ErrUserNotFound       = errors.New("user not found")
)

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	FindById(ctx context.Context, userID string) (*User, error)
	FindByEmail(ctx context.Context, tenantID string, email string) (*User, error)
	Update(ctx context.Context, userID string, u *UpdateUser) error
	Delete(ctx context.Context, userID string) error
}

type userRepository struct {
	db db.DBTX
}

func UserRepoImpl(db db.DBTX) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *User) (*User, error) {
	query := `INSERT INTO users (tenant_id, role_id, email, password_hash, full_name, status)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING user_id, created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query, u.TenantID, u.RoleId, u.Email, u.PasswordHash,
		u.FullName, u.Status).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		// Detect unique constraint violation (Postgres specific)
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Constraint, "users_email") {
				return nil, ErrEmailAlreadyExists
			}
		}
		return nil, fmt.Errorf("UserRepo.Create: %w", err)
	}
	return u, err
}

// FindById implements UserRepository.
func (r *userRepository) FindById(ctx context.Context, userID string) (*User, error) {
	query := `SELECT user_id, tenant_id, role_id, email, password_hash, full_name,
                     status, created_at, updated_at, last_login_at
              FROM users WHERE user_id=$1`
	row := r.db.QueryRowContext(ctx, query, userID)

	u := &User{}
	err := row.Scan(&u.ID, &u.TenantID, &u.RoleId, &u.Email, &u.PasswordHash,
		&u.FullName, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return u, err
}

// FindByEmail implements UserRepository.
func (r *userRepository) FindByEmail(ctx context.Context, tenantID string, email string) (*User, error) {
	query := `SELECT user_id, tenant_id, role_id, email, password_hash, full_name,
					status, created_at, updated_at, last_login_at
				FROM users 
				WHERE tenant_id=$1 AND email=$2`
	row := r.db.QueryRowContext(ctx, query, tenantID, email)
	u := &User{}
	err := row.Scan(&u.ID, &u.TenantID, &u.RoleId, &u.Email, &u.PasswordHash,
		&u.FullName, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return u, err
}

// Update implements UserRepository.
func (r *userRepository) Update(ctx context.Context, userID string, u *UpdateUser) error {
	fields := []string{}
	args := []interface{}{}
	argPos := 1

	if u.FullName != nil {
		fields = append(fields, fmt.Sprintf("full_name=$%d", argPos))
		args = append(args, *u.FullName)
		argPos++
	}
	if u.Email != nil {
		fields = append(fields, fmt.Sprintf("email=$%d", argPos))
		args = append(args, *u.Email)
		argPos++
	}
	if u.PasswordHash != nil {
		fields = append(fields, fmt.Sprintf("password_hash=$%d", argPos))
		args = append(args, *u.PasswordHash)
		argPos++
	}
	if u.Status != nil {
		fields = append(fields, fmt.Sprintf("status=$%d", argPos))
		args = append(args, *u.Status)
		argPos++
	}
	if u.LastLoginAt != nil {
		fields = append(fields, fmt.Sprintf("last_login_at=$%d", argPos))
		args = append(args, *u.LastLoginAt)
		argPos++
	}

	if len(fields) == 0 {
		return ErrNothingToUpdate
	}

	// build query
	query := fmt.Sprintf(`UPDATE users SET %s, updated_at=NOW() WHERE user_id=$%d`,
		strings.Join(fields, ", "), argPos)
	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("UserRepo.Update: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE user_id=$1`
	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("UserRepo.Delete: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
