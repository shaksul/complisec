from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go")
content = path.read_text(encoding="utf-8")

progress = """// Progress -----------------------------------------------------------------

func (r *TrainingRepo) CreateProgress(ctx context.Context, progress TrainingProgress) error {
	query := `
		INSERT INTO training_progress (
			id, assignment_id, material_id, progress_percentage, time_spent_minutes,
			last_position, completed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		progress.ID,
		progress.AssignmentID,
		progress.MaterialID,
		progress.ProgressPercentage,
		progress.TimeSpentMinutes,
		progress.LastPosition,
		progress.CompletedAt,
	)
	return err
}

func (r *TrainingRepo) UpdateProgress(ctx context.Context, progress TrainingProgress) error {
	query := `
		UPDATE training_progress SET
			progress_percentage = $4,
			time_spent_minutes = $5,
			last_position = $6,
			completed_at = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND assignment_id = $2 AND material_id = $3`

	_, err := r.db.ExecContext(ctx, query,
		progress.ID,
		progress.AssignmentID,
		progress.MaterialID,
		progress.ProgressPercentage,
		progress.TimeSpentMinutes,
		progress.LastPosition,
		progress.CompletedAt,
	)
	return err
}

func (r *TrainingRepo) GetProgressByAssignment(ctx context.Context, assignmentID string) ([]TrainingProgress, error) {
	query := `
		SELECT id, assignment_id, material_id, progress_percentage, time_spent_minutes,
			last_position, completed_at, created_at, updated_at
		FROM training_progress
		WHERE assignment_id = $1`

	rs, err := r.db.QueryContext(ctx, query, assignmentID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var result []TrainingProgress
	for rs.Next() {
		var progress TrainingProgress
		var lastPosition sql.NullInt64
		var completedAt sql.NullTime

		if err := rs.Scan(
			&progress.ID,
			&progress.AssignmentID,
			&progress.MaterialID,
			&progress.ProgressPercentage,
			&progress.TimeSpentMinutes,
			&lastPosition,
			&completedAt,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		); err != nil {
			return nil, err
		}

		progress.LastPosition = intPointer(lastPosition)
		progress.CompletedAt = timePointer(completedAt)
		result = append(result, progress)
	}

	return result, rs.Err()
}

func (r *TrainingRepo) GetProgressByAssignmentAndMaterial(ctx context.Context, assignmentID, materialID string) (*TrainingProgress, error) {
	query := `
		SELECT id, assignment_id, material_id, progress_percentage, time_spent_minutes,
			last_position, completed_at, created_at, updated_at
		FROM training_progress
		WHERE assignment_id = $1 AND material_id = $2`

	var progress TrainingProgress
	var lastPosition sql.NullInt64
	var completedAt sql.NullTime

	if err := r.db.QueryRowContext(ctx, query, assignmentID, materialID).Scan(
		&progress.ID,
		&progress.AssignmentID,
		&progress.MaterialID,
		&progress.ProgressPercentage,
		&progress.TimeSpentMinutes,
		&lastPosition,
		&completedAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	); err != nil {
		return nil, err
	}

	progress.LastPosition = intPointer(lastPosition)
	progress.CompletedAt = timePointer(completedAt)
	return &progress, nil
}
"""

content += "\n\n" + progress

path.write_text(content, encoding="utf-8")
