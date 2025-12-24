package refreshtoken

import "time"

type RefreshToken struct {
	ID        string    `db:"token_id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	TokenHash string    `db:"token_hash" json:"-"`
	Revoked   bool      `db:"revoked" json:"revoked"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
