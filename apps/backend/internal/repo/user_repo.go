package repo

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	TenantID     string
	Email        string
	PasswordHash string
	FirstName    *string
	LastName     *string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	log.Printf("DEBUG: user_repo.Create inserting user tenant=%s email=%s", u.TenantID, u.Email)
	_, err = r.db.ExecContext(ctx, `
        INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, is_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, u.ID, u.TenantID, u.Email, string(hashedPassword), u.FirstName, u.LastName, u.IsActive)
	if err != nil {
		log.Printf("ERROR: user_repo.Create insert failed: %v", err)
	}
	return err
}

func (r *UserRepo) GetByEmail(ctx context.Context, tenantID, email string) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE tenant_id = $1 AND email = $2
    `, tenantID, email)

	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("DEBUG: user_repo.GetByEmail no rows tenant=%s email=%s", tenantID, email)
			return nil, nil
		}
		log.Printf("ERROR: user_repo.GetByEmail query failed tenant=%s email=%s: %v", tenantID, email, err)
		return nil, err
	}
	log.Printf("DEBUG: user_repo.GetByEmail found user id=%s", u.ID)
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`, id)

	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context, tenantID string) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
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

func (r *UserRepo) ListPaginated(ctx context.Context, tenantID string, page, pageSize int) ([]User, int64, error) {
	offset := (page - 1) * pageSize

	// Получаем общее количество записей
	var total int64
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM users WHERE tenant_id = $1
	`, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Получаем данные с пагинацией
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, tenantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

// GetUsersByTenant retrieves all users for a tenant
func (r *UserRepo) GetUsersByTenant(ctx context.Context, tenantID string) ([]User, error) {
	return r.List(ctx, tenantID)
}

func (r *UserRepo) Update(ctx context.Context, u User) error {
	_, err := r.db.Exec(`
		UPDATE users SET first_name = $1, last_name = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, u.FirstName, u.LastName, u.IsActive, u.ID)
	return err
}

func (r *UserRepo) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT r.name FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *UserRepo) SetUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	// Remove existing roles
	_, err := r.db.Exec("DELETE FROM user_roles WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Add new roles
	for _, roleID := range roleIDs {
		_, err = r.db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepo) GetUserWithRoles(ctx context.Context, userID string) (*UserWithRoles, error) {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	roles, err := r.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserWithRoles{
		User:  *user,
		Roles: roles,
	}, nil
}

func (r *UserRepo) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT p.code FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
	`, userID)
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

type UserWithRoles struct {
	User
	Roles []string
}
