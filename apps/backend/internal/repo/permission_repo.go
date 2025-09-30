package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PermissionRepo struct {
	db *DB
}

func NewPermissionRepo(db *DB) *PermissionRepo {
	return &PermissionRepo{db: db}
}

func (r *PermissionRepo) List(ctx context.Context) ([]Permission, error) {
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

func (r *PermissionRepo) GetByCode(ctx context.Context, code string) (*Permission, error) {
	row := r.db.QueryRow(`
		SELECT id, code, module, description, created_at
		FROM permissions WHERE code = $1
	`, code)

	var perm Permission
	err := row.Scan(&perm.ID, &perm.Code, &perm.Module, &perm.Description, &perm.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

