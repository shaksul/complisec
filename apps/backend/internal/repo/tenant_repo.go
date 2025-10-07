package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Tenant struct {
	ID        string
	Name      string
	Domain    *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TenantRepo struct {
	db *DB
}

func NewTenantRepo(db *DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) Create(ctx context.Context, tenant Tenant) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tenants (id, name, domain)
		VALUES ($1, $2, $3)
	`, tenant.ID, tenant.Name, tenant.Domain)
	return err
}

func (r *TenantRepo) GetByID(ctx context.Context, id string) (*Tenant, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, domain, created_at, updated_at
		FROM tenants WHERE id = $1
	`, id)

	var tenant Tenant
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepo) GetByDomain(ctx context.Context, domain string) (*Tenant, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, domain, created_at, updated_at
		FROM tenants WHERE domain = $1
	`, domain)

	var tenant Tenant
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepo) List(ctx context.Context, limit, offset int) ([]Tenant, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, domain, created_at, updated_at
		FROM tenants 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []Tenant
	for rows.Next() {
		var tenant Tenant
		err := rows.Scan(&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.CreatedAt, &tenant.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

func (r *TenantRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenants`).Scan(&count)
	return count, err
}

func (r *TenantRepo) Update(ctx context.Context, id, name string, domain *string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tenants 
		SET name = $1, domain = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, name, domain, id)
	return err
}

func (r *TenantRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tenants WHERE id = $1`, id)
	return err
}

func (r *TenantRepo) ExistsByDomain(ctx context.Context, domain string, excludeID *string) (bool, error) {
	var query string
	var args []interface{}

	if excludeID != nil {
		query = `SELECT EXISTS(SELECT 1 FROM tenants WHERE domain = $1 AND id != $2)`
		args = []interface{}{domain, *excludeID}
	} else {
		query = `SELECT EXISTS(SELECT 1 FROM tenants WHERE domain = $1)`
		args = []interface{}{domain}
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	return exists, err
}
