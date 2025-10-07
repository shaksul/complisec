from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go")
content = path.read_text(encoding="utf-8")

roles = """// Role assignments ---------------------------------------------------------

func (r *TrainingRepo) CreateRoleAssignment(ctx context.Context, assignment RoleTrainingAssignment) error {
	query := `
		INSERT INTO role_training_assignments (
			id, tenant_id, role_id, material_id, course_id, is_required, due_days, assigned_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query,
		assignment.ID,
		assignment.TenantID,
		assignment.RoleID,
		assignment.MaterialID,
		assignment.CourseID,
		assignment.IsRequired,
		assignment.DueDays,
		assignment.AssignedBy,
	)
	return err
}

func (r *TrainingRepo) GetRoleAssignments(ctx context.Context, roleID string) ([]RoleTrainingAssignment, error) {
	query := `
		SELECT id, tenant_id, role_id, material_id, course_id, is_required, due_days, assigned_by, created_at
		FROM role_training_assignments
		WHERE role_id = $1
		ORDER BY created_at DESC`

	rs, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var assignments []RoleTrainingAssignment
	for rs.Next() {
		var assignment RoleTrainingAssignment
		var materialID sql.NullString
		var courseID sql.NullString
		var dueDays sql.NullInt64
		var assignedBy sql.NullString

		if err := rs.Scan(
			&assignment.ID,
			&assignment.TenantID,
			&assignment.RoleID,
			&materialID,
			&courseID,
			&assignment.IsRequired,
			&dueDays,
			&assignedBy,
			&assignment.CreatedAt,
		); err != nil {
			return nil, err
		}

		assignment.MaterialID = stringPointer(materialID)
		assignment.CourseID = stringPointer(courseID)
		assignment.DueDays = intPointer(dueDays)
		assignment.AssignedBy = stringPointer(assignedBy)

		assignments = append(assignments, assignment)
	}

	return assignments, rs.Err()
}

func (r *TrainingRepo) DeleteRoleAssignment(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM role_training_assignments WHERE id = $1", id)
	return err
}
"""

content += "\n\n" + roles

path.write_text(content, encoding="utf-8")
