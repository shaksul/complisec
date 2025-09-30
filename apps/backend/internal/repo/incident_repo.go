package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Incident struct {
	ID          string
	TenantID    string
	Title       string
	Description *string
	Severity    string
	Status      string
	AssetID     *string
	RiskID      *string
	AssignedTo  *string
	CreatedBy   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type IncidentRepo struct {
	db *DB
}

func NewIncidentRepo(db *DB) *IncidentRepo {
	return &IncidentRepo{db: db}
}

func (r *IncidentRepo) Create(ctx context.Context, incident Incident) error {
	_, err := r.db.Exec(`
		INSERT INTO incidents (id, tenant_id, title, description, severity, status, asset_id, risk_id, assigned_to, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, incident.ID, incident.TenantID, incident.Title, incident.Description, incident.Severity, incident.Status, incident.AssetID, incident.RiskID, incident.AssignedTo, incident.CreatedBy)
	return err
}

func (r *IncidentRepo) GetByID(ctx context.Context, id string) (*Incident, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, title, description, severity, status, asset_id, risk_id, assigned_to, created_by, created_at, updated_at
		FROM incidents WHERE id = $1
	`, id)

	var incident Incident
	err := row.Scan(&incident.ID, &incident.TenantID, &incident.Title, &incident.Description, &incident.Severity, &incident.Status, &incident.AssetID, &incident.RiskID, &incident.AssignedTo, &incident.CreatedBy, &incident.CreatedAt, &incident.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &incident, nil
}

func (r *IncidentRepo) List(ctx context.Context, tenantID string) ([]Incident, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, title, description, severity, status, asset_id, risk_id, assigned_to, created_by, created_at, updated_at
		FROM incidents WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []Incident
	for rows.Next() {
		var incident Incident
		err := rows.Scan(&incident.ID, &incident.TenantID, &incident.Title, &incident.Description, &incident.Severity, &incident.Status, &incident.AssetID, &incident.RiskID, &incident.AssignedTo, &incident.CreatedBy, &incident.CreatedAt, &incident.UpdatedAt)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}
	return incidents, nil
}

func (r *IncidentRepo) Update(ctx context.Context, incident Incident) error {
	_, err := r.db.Exec(`
		UPDATE incidents SET title = $1, description = $2, severity = $3, status = $4, asset_id = $5, risk_id = $6, assigned_to = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`, incident.Title, incident.Description, incident.Severity, incident.Status, incident.AssetID, incident.RiskID, incident.AssignedTo, incident.ID)
	return err
}

func (r *IncidentRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(`
		UPDATE incidents SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, status, id)
	return err
}

func (r *IncidentRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM incidents WHERE id = $1", id)
	return err
}

