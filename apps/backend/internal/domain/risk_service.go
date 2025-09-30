package domain

import (
	"context"
	"time"

	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

type RiskService struct {
	riskRepo  *repo.RiskRepo
	auditRepo *repo.AuditRepo
}

func NewRiskService(riskRepo *repo.RiskRepo, auditRepo *repo.AuditRepo) *RiskService {
	return &RiskService{
		riskRepo:  riskRepo,
		auditRepo: auditRepo,
	}
}

func (s *RiskService) CreateRisk(ctx context.Context, tenantID, title string, description, category *string, likelihood, impact int, ownerID, assetID *string) (*repo.Risk, error) {
	risk := repo.Risk{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Title:       title,
		Description: description,
		Category:    category,
		Likelihood:  &likelihood,
		Impact:      &impact,
		Status:      "draft",
		OwnerID:     ownerID,
		AssetID:     assetID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.riskRepo.Create(ctx, risk)
	if err != nil {
		return nil, err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, tenantID, "system", "create", "risk", &risk.ID, risk)

	return &risk, nil
}

func (s *RiskService) GetRisk(ctx context.Context, id string) (*repo.Risk, error) {
	return s.riskRepo.GetByID(ctx, id)
}

func (s *RiskService) ListRisks(ctx context.Context, tenantID string) ([]repo.Risk, error) {
	return s.riskRepo.List(ctx, tenantID)
}

func (s *RiskService) UpdateRisk(ctx context.Context, id, title string, description, category *string, likelihood, impact int, ownerID, assetID *string) error {
	risk, err := s.riskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if risk == nil {
		return nil
	}

	risk.Title = title
	risk.Description = description
	risk.Category = category
	risk.Likelihood = &likelihood
	risk.Impact = &impact
	risk.OwnerID = ownerID
	risk.AssetID = assetID
	risk.UpdatedAt = time.Now()

	err = s.riskRepo.Update(ctx, *risk)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, risk.TenantID, "system", "update", "risk", &id, risk)

	return nil
}

func (s *RiskService) UpdateRiskStatus(ctx context.Context, id, status string) error {
	risk, err := s.riskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if risk == nil {
		return nil
	}

	err = s.riskRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, risk.TenantID, "system", "update_status", "risk", &id, map[string]string{"status": status})

	return nil
}

func (s *RiskService) DeleteRisk(ctx context.Context, id string) error {
	risk, err := s.riskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if risk == nil {
		return nil
	}

	err = s.riskRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, risk.TenantID, "system", "delete", "risk", &id, nil)

	return nil
}
