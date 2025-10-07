from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go")
text = path.read_text(encoding="utf-8")

def normalize(value: str) -> str:
    value = value.lstrip("\n")
    return value.replace("\n", "\r\n")

def apply(old: str, new: str):
    global text
    old_norm = normalize(old)
    new_norm = normalize(new)
    if old_norm not in text:
        lines = [line for line in old.splitlines() if line.strip()]
        hint = lines[0] if lines else "<unknown>"
        raise SystemExit(f"pattern not found: {hint}")
    text = text.replace(old_norm, new_norm)

replacements = [
(
"""
func (r *TrainingRepo) CreateAssignment(ctx context.Context, assignment TrainingAssignment) error {
	// TODO: Implement
	return nil
}
""",
"""
func (r *TrainingRepo) CreateAssignment(ctx context.Context, assignment TrainingAssignment) error {
	query := `
		INSERT INTO train_assignments (
			id, tenant_id, material_id, course_id, user_id, status, due_at, completed_at,
			assigned_by, priority, progress_percentage, time_spent_minutes, last_accessed_at,
			reminder_sent_at, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	metadataJSON, _ := json.Marshal(assignment.Metadata)

	_, err := r.db.ExecContext(
		ctx,
		query,
		assignment.ID,
		assignment.TenantID,
		assignment.MaterialID,
		assignment.CourseID,
		assignment.UserID,
		assignment.Status,
		assignment.DueAt,
		assignment.CompletedAt,
		assignment.AssignedBy,
		assignment.Priority,
		assignment.ProgressPercentage,
		assignment.TimeSpentMinutes,
		assignment.LastAccessedAt,
		assignment.ReminderSentAt,
		metadataJSON,
	)

	return err
}
"""),
(
"""
func (r *TrainingRepo) GetAssignmentByID(ctx context.Context, id string) (*TrainingAssignment, error) {
	// TODO: Implement
	return nil, nil
}
""",
"""
func (r *TrainingRepo) GetAssignmentByID(ctx context.Context, id string) (*TrainingAssignment, error) {
	query := `
		SELECT id, tenant_id, material_id, course_id, user_id, status, due_at, completed_at,
			assigned_by, priority, progress_percentage, time_spent_minutes, last_accessed_at,
			reminder_sent_at, metadata, created_at
		FROM train_assignments
		WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	assignment, err := scanAssignment(row)
	if err != nil {
		return nil, err
	}

	return assignment, nil
}
"""),
(
"""
func (r *TrainingRepo) GetUserAssignments(ctx context.Context, userID string, filters map[string]interface{}) ([]TrainingAssignment, error) {
	// TODO: Implement
	return nil, nil
}
""",
"""
func (r *TrainingRepo) GetUserAssignments(ctx context.Context, userID string, filters map[string]interface{}) ([]TrainingAssignment, error) {
	query := `
		SELECT id, tenant_id, material_id, course_id, user_id, status, due_at, completed_at,
			assigned_by, priority, progress_percentage, time_spent_minutes, last_accessed_at,
			reminder_sent_at, metadata, created_at
		FROM train_assignments
		WHERE user_id = $1`

	args := []interface{}{userID}
	argIdx := 2

	if tenantID, ok := filters["tenant_id"]; ok {
		query += fmt.Sprintf(" AND tenant_id = $%d", argIdx)
		args = append(args, tenantID)
		argIdx++
	}

	if status, ok := filters["status"]; ok {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	if courseID, ok := filters["course_id"]; ok {
		query += fmt.Sprintf(" AND course_id = $%d", argIdx)
		args = append(args, courseID)
		argIdx++
	}

	if materialID, ok := filters["material_id"]; ok {
		query += fmt.Sprintf(" AND material_id = $%d", argIdx)
		args = append(args, materialID)
		argIdx++
	}

	if dueBefore, ok := filters["due_before"]; ok {
		query += fmt.Sprintf(" AND due_at <= $%d", argIdx)
		args = append(args, dueBefore)
		argIdx++
	}

	if dueAfter, ok := filters["due_after"]; ok {
		query += fmt.Sprintf(" AND due_at >= $%d", argIdx)
		args = append(args, dueAfter)
		argIdx++
	}

	if overdueOnly, ok := filters["overdue_only"]; ok {
		if flag, cast := overdueOnly.(bool); cast && flag {
			query += " AND due_at IS NOT NULL AND due_at < CURRENT_TIMESTAMP AND status <> 'completed'"
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []TrainingAssignment
	for rows.Next() {
		assignment, err := scanAssignment(rows)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, *assignment)
	}

	return assignments, rows.Err()
}
"""),
