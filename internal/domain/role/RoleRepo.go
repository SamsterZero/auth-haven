package role

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
	ErrRoleAlreadyExists = errors.New("role with this name already exists for this tenant")
	ErrRoleNotFound      = errors.New("role not found")
)

type RoleRepository interface {
	Create(ctx context.Context, r *Role) (*Role, error)
	FindById(ctx context.Context, roleID int64) (*Role, error)
	FindByTenantAndName(ctx context.Context, tenantID, name string) (*Role, error)
	Update(ctx context.Context, roleID int64, r *UpdateRole) error
	Delete(ctx context.Context, roleID int64) error
}

type roleRepository struct {
	db db.DBTX
}

func RoleRepoImpl(db db.DBTX) RoleRepository {
	return &roleRepository{db: db}
}

// Create inserts a new role
func (r *roleRepository) Create(ctx context.Context, role *Role) (*Role, error) {
	query := `INSERT INTO roles (tenant_id, name, permissions)
              VALUES ($1, $2, $3)
              RETURNING role_id, created_at`
	err := r.db.QueryRowContext(ctx, query, role.TenantID, role.Name, role.Permissions).
		Scan(&role.ID, &role.CreatedAt)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Constraint, "roles_tenant_id_name_key") {
				return nil, ErrRoleAlreadyExists
			}
		}
		return nil, fmt.Errorf("RoleRepo.Create: %w", err)
	}
	return role, nil
}

// FindById returns a role by ID
func (r *roleRepository) FindById(ctx context.Context, roleID int64) (*Role, error) {
	query := `SELECT role_id, tenant_id, name, permissions, created_at
              FROM roles WHERE role_id=$1`
	row := r.db.QueryRowContext(ctx, query, roleID)

	role := &Role{}
	err := row.Scan(&role.ID, &role.TenantID, &role.Name, &role.Permissions, &role.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("RoleRepo.FindById: %w", err)
	}
	return role, nil
}

// FindByTenantAndName returns a role by tenant ID and role name
func (r *roleRepository) FindByTenantAndName(ctx context.Context, tenantID, name string) (*Role, error) {
	query := `SELECT role_id, tenant_id, name, permissions, created_at
              FROM roles WHERE tenant_id=$1 AND name=$2`
	row := r.db.QueryRowContext(ctx, query, tenantID, name)

	role := &Role{}
	err := row.Scan(&role.ID, &role.TenantID, &role.Name, &role.Permissions, &role.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("RoleRepo.FindByTenantAndName: %w", err)
	}
	return role, nil
}

// Update modifies an existing role
func (r *roleRepository) Update(ctx context.Context, roleID int64, role *UpdateRole) error {
	fields := []string{}
	args := []interface{}{}
	argPos := 1

	if role.Name != nil {
		fields = append(fields, fmt.Sprintf("name=$%d", argPos))
		args = append(args, *role.Name)
		argPos++
	}
	if role.Permissions != nil {
		fields = append(fields, fmt.Sprintf("permissions=$%d", argPos))
		args = append(args, *role.Permissions)
		argPos++
	}

	if len(fields) == 0 {
		return errors.New("nothing to update")
	}

	query := fmt.Sprintf(`UPDATE roles SET %s WHERE role_id=$%d`,
		strings.Join(fields, ", "), argPos)
	args = append(args, roleID)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("RoleRepo.Update: %w", err)
	}
	return nil
}

// Delete removes a role
func (r *roleRepository) Delete(ctx context.Context, roleID int64) error {
	query := `DELETE FROM roles WHERE role_id=$1`
	res, err := r.db.ExecContext(ctx, query, roleID)
	if err != nil {
		return fmt.Errorf("RoleRepo.Delete: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrRoleNotFound
	}
	return nil
}
