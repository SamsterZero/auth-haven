package tenant

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
	ErrTenantAlreadyExists = errors.New("tenant with this domain already exists")
	ErrTenantNotFound      = errors.New("tenant not found")
)

type TenantRepository interface {
	Create(ctx context.Context, t *Tenant) (*Tenant, error)
	FindById(ctx context.Context, tenantID string) (*Tenant, error)
	FindByDomain(ctx context.Context, domain string) (*Tenant, error)
	Update(ctx context.Context, tenantID string, t *UpdateTenant) error
	Delete(ctx context.Context, tenantID string) error
}

type tenantRepository struct {
	db db.DBTX
}

func TenantRepoImpl(db db.DBTX) TenantRepository {
	return &tenantRepository{db: db}
}

// Create implements TenantRepository.
func (r *tenantRepository) Create(ctx context.Context, t *Tenant) (*Tenant, error) {
	query := `INSERT INTO tenants (name, domain, status)
              VALUES ($1, $2, $3)
              RETURNING tenant_id, created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query, t.Name, t.Domain, t.Status).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" && pgErr.Constraint == "tenants_domain_key" {
				return nil, ErrTenantAlreadyExists
			}
		}
		return nil, fmt.Errorf("TenantRepo.Create: %w", err)
	}
	return t, nil
}

// FindByDomain implements TenantRepository.
func (r *tenantRepository) FindByDomain(ctx context.Context, domain string) (*Tenant, error) {
	query := `SELECT tenant_id, name, domain, status, created_at, updated_at
              FROM tenants WHERE domain=$1`
	row := r.db.QueryRowContext(ctx, query, domain)

	t := &Tenant{}
	err := row.Scan(&t.ID, &t.Name, &t.Domain, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("TenantRepo.FindByDomain: %w", err)
	}
	return t, nil
}

// FindById implements TenantRepository.
func (r *tenantRepository) FindById(ctx context.Context, tenantID string) (*Tenant, error) {
	query := `SELECT tenant_id, name, domain, status, created_at, updated_at
              FROM tenants WHERE tenant_id=$1`
	row := r.db.QueryRowContext(ctx, query, tenantID)

	t := &Tenant{}
	err := row.Scan(&t.ID, &t.Name, &t.Domain, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("TenantRepo.FindById: %w", err)
	}
	return t, nil
}

// Update implements TenantRepository.
func (r *tenantRepository) Update(ctx context.Context, tenantID string, t *UpdateTenant) error {
	fields := []string{}
	args := []interface{}{}
	argPos := 1

	if t.Name != nil {
		fields = append(fields, fmt.Sprintf("name=$%d", argPos))
		args = append(args, t.Name)
		argPos++
	}
	if t.Domain != nil {
		fields = append(fields, fmt.Sprintf("domain=$%d", argPos))
		args = append(args, t.Domain)
		argPos++
	}
	if t.Status != nil {
		fields = append(fields, fmt.Sprintf("status=$%d", argPos))
		args = append(args, t.Status)
		argPos++
	}

	if len(fields) == 0 {
		return errors.New("nothing to update")
	}

	query := fmt.Sprintf(`UPDATE tenants SET %s, updated_at=NOW() WHERE tenant_id=$%d`,
		strings.Join(fields, ", "), argPos)
	args = append(args, tenantID)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("TenantRepo.Update: %w", err)
	}
	return nil
}

// Delete implements TenantRepository.
func (r *tenantRepository) Delete(ctx context.Context, tenantID string) error {
	query := `DELETE FROM tenants WHERE tenant_id=$1`
	res, err := r.db.ExecContext(ctx, query, tenantID)
	if err != nil {
		return fmt.Errorf("TenantRepo.Delete: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrTenantNotFound
	}
	return nil
}
