package domain

import (
	"context"
	"time"

	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

type IncidentService struct {
	incidentRepo *repo.IncidentRepo
	auditRepo    *repo.AuditRepo
}

func NewIncidentService(incidentRepo *repo.IncidentRepo, auditRepo *repo.AuditRepo) *IncidentService {
	return &IncidentService{
		incidentRepo: incidentRepo,
		auditRepo:    auditRepo,
	}
}

func (s *IncidentService) CreateIncident(ctx context.Context, tenantID, title, severity string, description, assetID, riskID, assignedTo, createdBy *string) (*repo.Incident, error) {
	incident := repo.Incident{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Title:       title,
		Description: description,
		Severity:    severity,
		Status:      "new",
		AssetID:     assetID,
		RiskID:      riskID,
		AssignedTo:  assignedTo,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.incidentRepo.Create(ctx, incident)
	if err != nil {
		return nil, err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, tenantID, "system", "create", "incident", &incident.ID, incident)

	return &incident, nil
}

func (s *IncidentService) GetIncident(ctx context.Context, id string) (*repo.Incident, error) {
	return s.incidentRepo.GetByID(ctx, id)
}

func (s *IncidentService) ListIncidents(ctx context.Context, tenantID string) ([]repo.Incident, error) {
	return s.incidentRepo.List(ctx, tenantID)
}

func (s *IncidentService) UpdateIncident(ctx context.Context, id, title, severity string, description, assetID, riskID, assignedTo *string) error {
	incident, err := s.incidentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if incident == nil {
		return nil
	}

	incident.Title = title
	incident.Description = description
	incident.Severity = severity
	incident.AssetID = assetID
	incident.RiskID = riskID
	incident.AssignedTo = assignedTo
	incident.UpdatedAt = time.Now()

	err = s.incidentRepo.Update(ctx, *incident)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, incident.TenantID, "system", "update", "incident", &id, incident)

	return nil
}

func (s *IncidentService) UpdateIncidentStatus(ctx context.Context, id, status string) error {
	incident, err := s.incidentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if incident == nil {
		return nil
	}

	err = s.incidentRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, incident.TenantID, "system", "update_status", "incident", &id, map[string]string{"status": status})

	return nil
}

func (s *IncidentService) DeleteIncident(ctx context.Context, id string) error {
	incident, err := s.incidentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if incident == nil {
		return nil
	}

	err = s.incidentRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, incident.TenantID, "system", "delete", "incident", &id, nil)

	return nil
}
