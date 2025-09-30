package domain

import (
	"context"
	"risknexus/backend/internal/repo"
)

type ComplianceService struct {
	repo *repo.ComplianceRepo
}

func NewComplianceService(r *repo.ComplianceRepo) *ComplianceService {
	return &ComplianceService{repo: r}
}

func (s *ComplianceService) ListStandards(ctx context.Context, tenantID string) ([]repo.ComplianceStandard, error) {
	return s.repo.ListStandards(ctx, tenantID)
}

func (s *ComplianceService) CreateStandard(ctx context.Context, standard repo.ComplianceStandard) error {
	return s.repo.CreateStandard(ctx, standard)
}

func (s *ComplianceService) ListRequirements(ctx context.Context, standardID string) ([]repo.ComplianceRequirement, error) {
	return s.repo.ListRequirements(ctx, standardID)
}

func (s *ComplianceService) CreateRequirement(ctx context.Context, requirement repo.ComplianceRequirement) error {
	return s.repo.CreateRequirement(ctx, requirement)
}

func (s *ComplianceService) ListAssessments(ctx context.Context, tenantID string) ([]repo.ComplianceAssessment, error) {
	return s.repo.ListAssessments(ctx, tenantID)
}

func (s *ComplianceService) CreateAssessment(ctx context.Context, assessment repo.ComplianceAssessment) error {
	return s.repo.CreateAssessment(ctx, assessment)
}

func (s *ComplianceService) UpdateAssessment(ctx context.Context, id string, assessment repo.ComplianceAssessment) error {
	return s.repo.UpdateAssessment(ctx, id, assessment)
}

func (s *ComplianceService) ListGaps(ctx context.Context, assessmentID string) ([]repo.ComplianceGap, error) {
	return s.repo.ListGaps(ctx, assessmentID)
}

func (s *ComplianceService) CreateGap(ctx context.Context, gap repo.ComplianceGap) error {
	return s.repo.CreateGap(ctx, gap)
}

func (s *ComplianceService) UpdateGap(ctx context.Context, id string, gap repo.ComplianceGap) error {
	return s.repo.UpdateGap(ctx, id, gap)
}
