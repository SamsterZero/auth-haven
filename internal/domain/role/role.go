package role

import "time"

type Role struct {
	ID          int64     `db:"role_id" json:"id"`
	TenantID    string    `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Permissions []byte    `db:"permissions" json:"permissions"` // store JSON as []byte or map[string]any
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type UpdateRole struct {
	Name        *string
	Permissions *[]byte
}
