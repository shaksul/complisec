package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Role struct {
	ID          string
	TenantID    string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permission struct {
	ID          string
	Code        string
	Module      string
	Description *string
	CreatedAt   time.Time
}

type RoleRepo struct {
	db *DB
}

func NewRoleRepo(db *DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) Create(ctx context.Context, tenantID, name, description string) (*Role, error) {
	role := &Role{
		ID:          generateUUID(),
		TenantID:    tenantID,
		Name:        name,
		Description: &description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := r.db.Exec(`
		INSERT INTO roles (id, tenant_id, name, description)
		VALUES ($1, $2, $3, $4)
	`, role.ID, role.TenantID, role.Name, role.Description)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *RoleRepo) GetByID(ctx context.Context, id string) (*Role, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, name, description, created_at, updated_at
		FROM roles WHERE id = $1
	`, id)

	var role Role
	err := row.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) GetByName(ctx context.Context, tenantID, name string) (*Role, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, name, description, created_at, updated_at
		FROM roles WHERE tenant_id = $1 AND name = $2
	`, tenantID, name)

	var role Role
	err := row.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepo) List(ctx context.Context, tenantID string) ([]Role, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, name, description, created_at, updated_at
		FROM roles WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepo) Update(ctx context.Context, id, name, description string) error {
	_, err := r.db.Exec(`
		UPDATE roles SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, name, description, id)
	return err
}

func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM roles WHERE id = $1", id)
	return err
}

func (r *RoleRepo) GetPermissions(ctx context.Context, tenantID string) ([]Permission, error) {
	rows, err := r.db.Query(`
		SELECT id, code, module, description, created_at
		FROM permissions ORDER BY module, code
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		err := rows.Scan(&perm.ID, &perm.Code, &perm.Module, &perm.Description, &perm.CreatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

func (r *RoleRepo) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT p.code FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

func (r *RoleRepo) SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing permissions
	_, err = tx.ExecContext(ctx, "DELETE FROM role_permissions WHERE role_id = $1", roleID)
	if err != nil {
		return err
	}

	// Add new permissions
	for _, permID := range permissionIDs {
		_, err = tx.ExecContext(ctx, "INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)", roleID, permID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *RoleRepo) GetRoleWithPermissions(ctx context.Context, roleID string) (*RoleWithPermissions, error) {
	// Оптимизированный запрос с JOIN для получения роли и прав за один запрос
	row := r.db.QueryRow(`
		SELECT r.id, r.tenant_id, r.name, r.description, r.created_at, r.updated_at,
		       COALESCE(array_agg(p.code ORDER BY p.code), '{}') as permissions
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE r.id = $1
		GROUP BY r.id, r.tenant_id, r.name, r.description, r.created_at, r.updated_at
	`, roleID)

	var role Role
	var permissions pq.StringArray
	err := row.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description,
		&role.CreatedAt, &role.UpdatedAt, &permissions)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &RoleWithPermissions{
		Role:        role,
		Permissions: []string(permissions),
	}, nil
}

func (r *RoleRepo) GetUsersByRole(ctx context.Context, roleID string) ([]User, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.tenant_id, u.email, u.password_hash, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		WHERE ur.role_id = $1
		ORDER BY u.created_at DESC
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

type RoleWithPermissions struct {
	Role
	Permissions []string
}

// ListWithPermissions получает список ролей с правами за один запрос
func (r *RoleRepo) ListWithPermissions(ctx context.Context, tenantID string) ([]RoleWithPermissions, error) {
	rows, err := r.db.Query(`
		SELECT r.id, r.tenant_id, r.name, r.description, r.created_at, r.updated_at,
		       COALESCE(array_agg(p.code ORDER BY p.code), '{}') as permissions
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE r.tenant_id = $1
		GROUP BY r.id, r.tenant_id, r.name, r.description, r.created_at, r.updated_at
		ORDER BY r.created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []RoleWithPermissions
	for rows.Next() {
		var role Role
		var permissions pq.StringArray
		err := rows.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description,
			&role.CreatedAt, &role.UpdatedAt, &permissions)
		if err != nil {
			return nil, err
		}
		roles = append(roles, RoleWithPermissions{
			Role:        role,
			Permissions: []string(permissions),
		})
	}
	return roles, nil
}
