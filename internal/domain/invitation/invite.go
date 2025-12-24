package invitation

import "time"

type Invitation struct {
	ID        string     `db:"invitation_id" json:"id"`
	TenantID  string     `db:"tenant_id" json:"tenant_id"`
	RoleId    int64      `db:"role_id" json:"role_id"`
	Email     string     `db:"email" json:"email"`
	Token     string     `db:"token" json:"-"`
	Status    string     `db:"status" json:"status"`
	ExpiresAt *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}
