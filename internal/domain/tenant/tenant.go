package tenant

import "time"

type Tenant struct {
	ID        string    `db:"tenant_id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Domain    string    `db:"domain" json:"domain"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UpdateTenant struct {
	Name   *string
	Domain *string
	Status *string
}
