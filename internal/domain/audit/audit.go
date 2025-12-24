package audit

import "time"

type AuditLog struct {
	ID        string    `db:"log_id" json:"id"`
	UserID    *string   `db:"user_id" json:"user_id,omitempty"`
	TenantID  *string   `db:"tenant_id" json:"tenant_id,omitempty"`
	Action    string    `db:"action" json:"action"`
	IPAddress *string   `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent *string   `db:"user_agent" json:"user_agent,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
