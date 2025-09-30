package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Asset struct {
	ID        string
	TenantID  string
	Name      string
	InvCode   *string
	Type      string
	Status    string
	OwnerID   *string
	Location  *string
	Software  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AssetRepo struct {
	db *DB
}

func NewAssetRepo(db *DB) *AssetRepo {
	return &AssetRepo{db: db}
}

func (r *AssetRepo) Create(ctx context.Context, asset Asset) error {
	_, err := r.db.Exec(`
		INSERT INTO assets (id, tenant_id, name, inv_code, type, status, owner_id, location, software)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, asset.ID, asset.TenantID, asset.Name, asset.InvCode, asset.Type, asset.Status, asset.OwnerID, asset.Location, asset.Software)
	return err
}

func (r *AssetRepo) GetByID(ctx context.Context, id string) (*Asset, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, name, inv_code, type, status, owner_id, location, software, created_at, updated_at
		FROM assets WHERE id = $1
	`, id)

	var asset Asset
	err := row.Scan(&asset.ID, &asset.TenantID, &asset.Name, &asset.InvCode, &asset.Type, &asset.Status, &asset.OwnerID, &asset.Location, &asset.Software, &asset.CreatedAt, &asset.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepo) List(ctx context.Context, tenantID string) ([]Asset, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, name, inv_code, type, status, owner_id, location, software, created_at, updated_at
		FROM assets WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		err := rows.Scan(&asset.ID, &asset.TenantID, &asset.Name, &asset.InvCode, &asset.Type, &asset.Status, &asset.OwnerID, &asset.Location, &asset.Software, &asset.CreatedAt, &asset.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

func (r *AssetRepo) Update(ctx context.Context, asset Asset) error {
	_, err := r.db.Exec(`
		UPDATE assets SET name = $1, inv_code = $2, type = $3, status = $4, owner_id = $5, location = $6, software = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`, asset.Name, asset.InvCode, asset.Type, asset.Status, asset.OwnerID, asset.Location, asset.Software, asset.ID)
	return err
}

func (r *AssetRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM assets WHERE id = $1", id)
	return err
}

