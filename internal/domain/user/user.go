package user

import "time"

type User struct {
	ID           string     `db:"user_id" json:"id"`
	TenantID     string     `db:"tenant_id" json:"tenant_id"`
	RoleId       int64      `db:"role_id" json:"role_id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	FullName     string     `db:"full_name" json:"full_name"`
	Status       string     `db:"status" json:"status"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
}

type UpdateUser struct {
	FullName     *string
	Email        *string
	PasswordHash *string
	Status       *string
	LastLoginAt  *time.Time
}
