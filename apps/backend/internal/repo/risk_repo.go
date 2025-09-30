package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Risk struct {
	ID          string
	TenantID    string
	Title       string
	Description *string
	Category    *string
	Likelihood  *int
	Impact      *int
	Level       *int
	Status      string
	OwnerID     *string
	AssetID     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RiskRepo struct {
	db *DB
}

func NewRiskRepo(db *DB) *RiskRepo {
	return &RiskRepo{db: db}
}

func (r *RiskRepo) Create(ctx context.Context, risk Risk) error {
	_, err := r.db.Exec(`
		INSERT INTO risks (id, tenant_id, title, description, category, likelihood, impact, status, owner_id, asset_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, risk.ID, risk.TenantID, risk.Title, risk.Description, risk.Category, risk.Likelihood, risk.Impact, risk.Status, risk.OwnerID, risk.AssetID)
	return err
}

func (r *RiskRepo) GetByID(ctx context.Context, id string) (*Risk, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, title, description, category, likelihood, impact, level, status, owner_id, asset_id, created_at, updated_at
		FROM risks WHERE id = $1
	`, id)

	var risk Risk
	err := row.Scan(&risk.ID, &risk.TenantID, &risk.Title, &risk.Description, &risk.Category, &risk.Likelihood, &risk.Impact, &risk.Level, &risk.Status, &risk.OwnerID, &risk.AssetID, &risk.CreatedAt, &risk.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &risk, nil
}

func (r *RiskRepo) List(ctx context.Context, tenantID string) ([]Risk, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, title, description, category, likelihood, impact, level, status, owner_id, asset_id, created_at, updated_at
		FROM risks WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var risks []Risk
	for rows.Next() {
		var risk Risk
		err := rows.Scan(&risk.ID, &risk.TenantID, &risk.Title, &risk.Description, &risk.Category, &risk.Likelihood, &risk.Impact, &risk.Level, &risk.Status, &risk.OwnerID, &risk.AssetID, &risk.CreatedAt, &risk.UpdatedAt)
		if err != nil {
			return nil, err
		}
		risks = append(risks, risk)
	}
	return risks, nil
}

func (r *RiskRepo) Update(ctx context.Context, risk Risk) error {
	_, err := r.db.Exec(`
		UPDATE risks SET title = $1, description = $2, category = $3, likelihood = $4, impact = $5, status = $6, owner_id = $7, asset_id = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
	`, risk.Title, risk.Description, risk.Category, risk.Likelihood, risk.Impact, risk.Status, risk.OwnerID, risk.AssetID, risk.ID)
	return err
}

func (r *RiskRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(`
		UPDATE risks SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, status, id)
	return err
}

func (r *RiskRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM risks WHERE id = $1", id)
	return err
}

