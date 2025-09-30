package domain

import (
	"context"
	"time"

	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

type TrainingService struct {
	trainingRepo *repo.TrainingRepo
	auditRepo    *repo.AuditRepo
}

func NewTrainingService(trainingRepo *repo.TrainingRepo, auditRepo *repo.AuditRepo) *TrainingService {
	return &TrainingService{
		trainingRepo: trainingRepo,
		auditRepo:    auditRepo,
	}
}

func (s *TrainingService) CreateMaterial(ctx context.Context, tenantID, title, materialType, uri string, description, createdBy *string) (*repo.Material, error) {
	material := repo.Material{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Title:       title,
		Description: description,
		URI:         uri,
		Type:        materialType,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.trainingRepo.CreateMaterial(ctx, material)
	if err != nil {
		return nil, err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, tenantID, "system", "create", "material", &material.ID, material)

	return &material, nil
}

func (s *TrainingService) ListMaterials(ctx context.Context, tenantID string) ([]repo.Material, error) {
	return s.trainingRepo.ListMaterials(ctx, tenantID)
}

func (s *TrainingService) CreateAssignment(ctx context.Context, tenantID, materialID, userID string, dueAt *time.Time) (*repo.TrainingAssignment, error) {
	assignment := repo.TrainingAssignment{
		ID:         uuid.New().String(),
		TenantID:   tenantID,
		MaterialID: materialID,
		UserID:     userID,
		Status:     "assigned",
		DueAt:      dueAt,
		CreatedAt:  time.Now(),
	}

	err := s.trainingRepo.CreateAssignment(ctx, assignment)
	if err != nil {
		return nil, err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, tenantID, "system", "create", "assignment", &assignment.ID, assignment)

	return &assignment, nil
}

func (s *TrainingService) GetUserAssignments(ctx context.Context, userID string) ([]repo.TrainingAssignment, error) {
	return s.trainingRepo.GetUserAssignments(ctx, userID)
}

func (s *TrainingService) CompleteAssignment(ctx context.Context, assignmentID string) error {
	err := s.trainingRepo.UpdateAssignmentStatus(ctx, assignmentID, "completed")
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "system", "system", "complete", "assignment", &assignmentID, map[string]string{"assignment_id": assignmentID})

	return nil
}
