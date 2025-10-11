package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type AIProvider struct {
	ID             string   `json:"id"`
	TenantID       string   `json:"tenant_id"`
	Name           string   `json:"name"`
	BaseURL        string   `json:"base_url"`
	APIKey         *string  `json:"api_key,omitempty"`
	Roles          []string `json:"roles"`
	PromptTemplate *string  `json:"prompt_template,omitempty"`
	Models         []string `json:"models"`
	DefaultModel   string   `json:"default_model"`
	IsActive       bool     `json:"is_active"`
}

type AIRepo struct {
	db *DB
}

func NewAIRepo(db *DB) *AIRepo {
	return &AIRepo{db: db}
}

func (r *AIRepo) List(ctx context.Context, tenantID string) ([]AIProvider, error) {
	rows, err := r.db.Query(`SELECT id, tenant_id, name, base_url, api_key, roles, prompt_template, models, default_model, is_active FROM ai_providers WHERE tenant_id=$1`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AIProvider
	for rows.Next() {
		var p AIProvider
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.BaseURL, &p.APIKey, pq.Array(&p.Roles), &p.PromptTemplate, pq.Array(&p.Models), &p.DefaultModel, &p.IsActive); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, nil
}

func (r *AIRepo) Create(ctx context.Context, p AIProvider) error {
	_, err := r.db.Exec(`INSERT INTO ai_providers(id,tenant_id,name,base_url,api_key,roles,prompt_template,models,default_model,is_active) VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,$9)`, p.TenantID, p.Name, p.BaseURL, p.APIKey, pq.Array(p.Roles), p.PromptTemplate, pq.Array(p.Models), p.DefaultModel, p.IsActive)
	return err
}

func (r *AIRepo) Get(ctx context.Context, id string) (*AIProvider, error) {
	row := r.db.QueryRow(`SELECT id, tenant_id, name, base_url, api_key, roles, prompt_template, models, default_model, is_active FROM ai_providers WHERE id=$1`, id)
	var p AIProvider
	if err := row.Scan(&p.ID, &p.TenantID, &p.Name, &p.BaseURL, &p.APIKey, pq.Array(&p.Roles), &p.PromptTemplate, pq.Array(&p.Models), &p.DefaultModel, &p.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *AIRepo) Update(ctx context.Context, p AIProvider) error {
	_, err := r.db.Exec(`UPDATE ai_providers SET name=$1, base_url=$2, api_key=$3, roles=$4, prompt_template=$5, models=$6, default_model=$7, is_active=$8, updated_at=CURRENT_TIMESTAMP WHERE id=$9`, p.Name, p.BaseURL, p.APIKey, pq.Array(p.Roles), p.PromptTemplate, pq.Array(p.Models), p.DefaultModel, p.IsActive, p.ID)
	return err
}

func (r *AIRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM ai_providers WHERE id=$1`, id)
	return err
}
