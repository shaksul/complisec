package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type AIProvider struct {
	ID             string
	TenantID       string
	Name           string
	BaseURL        string
	APIKey         *string
	Roles          []string
	PromptTemplate *string
	IsActive       bool
}

type AIRepo struct {
	db *DB
}

func NewAIRepo(db *DB) *AIRepo {
	return &AIRepo{db: db}
}

func (r *AIRepo) List(ctx context.Context, tenantID string) ([]AIProvider, error) {
	rows, err := r.db.Query(`SELECT id, tenant_id, name, base_url, api_key, roles, prompt_template, is_active FROM ai_providers WHERE tenant_id=$1`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AIProvider
	for rows.Next() {
		var p AIProvider
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.BaseURL, &p.APIKey, &p.Roles, &p.PromptTemplate, &p.IsActive); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, nil
}

func (r *AIRepo) Create(ctx context.Context, p AIProvider) error {
	_, err := r.db.Exec(`INSERT INTO ai_providers(id,tenant_id,name,base_url,api_key,roles,prompt_template,is_active) VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7)`, p.TenantID, p.Name, p.BaseURL, p.APIKey, p.Roles, p.PromptTemplate, p.IsActive)
	return err
}

func (r *AIRepo) Get(ctx context.Context, id string) (*AIProvider, error) {
	row := r.db.QueryRow(`SELECT id, tenant_id, name, base_url, api_key, roles, prompt_template, is_active FROM ai_providers WHERE id=$1`, id)
	var p AIProvider
	if err := row.Scan(&p.ID, &p.TenantID, &p.Name, &p.BaseURL, &p.APIKey, &p.Roles, &p.PromptTemplate, &p.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

