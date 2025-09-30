package repo

import (
	"context"
	"time"
)

type Material struct {
	ID          string
	TenantID    string
	Title       string
	Description *string
	URI         string
	Type        string
	CreatedBy   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TrainingAssignment struct {
	ID          string
	TenantID    string
	MaterialID  string
	UserID      string
	Status      string
	DueAt       *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
}

type QuizQuestion struct {
	ID           string
	MaterialID   string
	Text         string
	OptionsJSON  string
	CorrectIndex int
	CreatedAt    time.Time
}

type QuizAttempt struct {
	ID          string
	UserID      string
	MaterialID  string
	Score       int
	Passed      bool
	AnswersJSON *string
	AttemptedAt time.Time
}

type TrainingRepo struct {
	db *DB
}

func NewTrainingRepo(db *DB) *TrainingRepo {
	return &TrainingRepo{db: db}
}

func (r *TrainingRepo) CreateMaterial(ctx context.Context, material Material) error {
	_, err := r.db.Exec(`
		INSERT INTO materials (id, tenant_id, title, description, uri, type, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, material.ID, material.TenantID, material.Title, material.Description, material.URI, material.Type, material.CreatedBy)
	return err
}

func (r *TrainingRepo) ListMaterials(ctx context.Context, tenantID string) ([]Material, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, title, description, uri, type, created_by, created_at, updated_at
		FROM materials WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []Material
	for rows.Next() {
		var material Material
		err := rows.Scan(&material.ID, &material.TenantID, &material.Title, &material.Description, &material.URI, &material.Type, &material.CreatedBy, &material.CreatedAt, &material.UpdatedAt)
		if err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}
	return materials, nil
}

func (r *TrainingRepo) CreateAssignment(ctx context.Context, assignment TrainingAssignment) error {
	_, err := r.db.Exec(`
		INSERT INTO train_assignments (id, tenant_id, material_id, user_id, status, due_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, assignment.ID, assignment.TenantID, assignment.MaterialID, assignment.UserID, assignment.Status, assignment.DueAt)
	return err
}

func (r *TrainingRepo) GetUserAssignments(ctx context.Context, userID string) ([]TrainingAssignment, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, material_id, user_id, status, due_at, completed_at, created_at
		FROM train_assignments WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []TrainingAssignment
	for rows.Next() {
		var assignment TrainingAssignment
		err := rows.Scan(&assignment.ID, &assignment.TenantID, &assignment.MaterialID, &assignment.UserID, &assignment.Status, &assignment.DueAt, &assignment.CompletedAt, &assignment.CreatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}
	return assignments, nil
}

func (r *TrainingRepo) UpdateAssignmentStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(`
		UPDATE train_assignments SET status = $1, completed_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, status, id)
	return err
}

