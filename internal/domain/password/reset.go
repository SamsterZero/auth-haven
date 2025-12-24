package password

import "time"

type PasswordReset struct {
	ID        string    `db:"reset_id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"-"`
	Status    string    `db:"status" json:"status"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
