from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go")
content = path.read_text(encoding="utf-8")

analytics = """// Analytics ----------------------------------------------------------------

func (r *TrainingRepo) CreateAnalytics(ctx context.Context, analytics TrainingAnalytics) error {
	query := `
		INSERT INTO training_analytics (
			id, tenant_id, user_id, material_id, course_id, metric_type, metric_value, recorded_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query,
		analytics.ID,
		analytics.TenantID,
		analytics.UserID,
		analytics.MaterialID,
		analytics.CourseID,
		analytics.MetricType,
		analytics.MetricValue,
		analytics.RecordedAt,
	)
	return err
}

func (r *TrainingRepo) GetUserAnalytics(ctx context.Context, userID string) (*TrainingAnalytics, error) {
	query := `
		SELECT id, tenant_id, user_id, material_id, course_id, metric_type, metric_value, recorded_at
		FROM training_analytics
		WHERE user_id = $1
		ORDER BY recorded_at DESC
		LIMIT 1`

	var analytics TrainingAnalytics
	var user sql.NullString
	var materialID sql.NullString
	var courseID sql.NullString

	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&analytics.ID,
		&analytics.TenantID,
		&user,
		&materialID,
		&courseID,
		&analytics.MetricType,
		&analytics.MetricValue,
		&analytics.RecordedAt,
	); err != nil {
		return nil, err
	}

	analytics.UserID = stringPointer(user)
	analytics.MaterialID = stringPointer(materialID)
	analytics.CourseID = stringPointer(courseID)

	return &analytics, nil
}

func (r *TrainingRepo) GetCourseAnalytics(ctx context.Context, courseID string) (*CourseAnalytics, error) {
	query := `
		WITH latest_scores AS (
			SELECT assignment_id, MAX(score) AS score
			FROM quiz_attempts
			WHERE assignment_id IS NOT NULL
			GROUP BY assignment_id
		)
		SELECT
			COUNT(a.*) AS total_assignments,
			COUNT(*) FILTER (WHERE a.status = 'completed') AS completed_assignments,
			COALESCE(AVG(a.time_spent_minutes), 0) AS avg_time_spent,
			COALESCE(AVG(ls.score::numeric), 0) AS avg_score
		FROM train_assignments a
		LEFT JOIN latest_scores ls ON ls.assignment_id = a.id
		WHERE a.course_id = $1`

	var total sql.NullInt64
	var completed sql.NullInt64
	var avgTime sql.NullFloat64
	var avgScore sql.NullFloat64

	if err := r.db.QueryRowContext(ctx, query, courseID).Scan(&total, &completed, &avgTime, &avgScore); err != nil {
		return nil, err
	}

	analytics := &CourseAnalytics{CourseID: courseID}
	if total.Valid {
		analytics.TotalAssignments = int(total.Int64)
	}
	if completed.Valid {
		analytics.CompletedAssignments = int(completed.Int64)
	}
	if analytics.TotalAssignments > 0 {
		analytics.CompletionRate = float64(analytics.CompletedAssignments) / float64(analytics.TotalAssignments)
	}
	if avgTime.Valid {
		analytics.AverageTimeSpent = int(math.Round(avgTime.Float64))
	}
	if avgScore.Valid {
		analytics.AverageScore = avgScore.Float64
	}

	return analytics, nil
}

func (r *TrainingRepo) GetOrganizationAnalytics(ctx context.Context, tenantID string) (*OrganizationAnalytics, error) {
	query := `
		SELECT
			COALESCE((SELECT COUNT(*) FROM materials WHERE tenant_id = $1), 0) AS total_materials,
			COALESCE((SELECT COUNT(*) FROM training_courses WHERE tenant_id = $1), 0) AS total_courses,
			COUNT(a.*) AS total_assignments,
			COUNT(*) FILTER (WHERE a.status = 'completed') AS completed_assignments,
			COUNT(*) FILTER (WHERE a.status <> 'completed' AND a.due_at IS NOT NULL AND a.due_at < CURRENT_TIMESTAMP) AS overdue_assignments,
			COALESCE(AVG(a.time_spent_minutes), 0) AS avg_time_spent
		FROM train_assignments a
		WHERE a.tenant_id = $1`

	var totalMaterials sql.NullInt64
	var totalCourses sql.NullInt64
	var totalAssignments sql.NullInt64
	var completedAssignments sql.NullInt64
	var overdueAssignments sql.NullInt64
	var avgTime sql.NullFloat64

	if err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&totalMaterials,
		&totalCourses,
		&totalAssignments,
		&completedAssignments,
		&overdueAssignments,
		&avgTime,
	); err != nil {
		return nil, err
	}

	analytics := &OrganizationAnalytics{TenantID: tenantID}
	if totalMaterials.Valid {
		analytics.TotalMaterials = int(totalMaterials.Int64)
	}
	if totalCourses.Valid {
		analytics.TotalCourses = int(totalCourses.Int64)
	}
	if totalAssignments.Valid {
		analytics.TotalAssignments = int(totalAssignments.Int64)
	}
	if completedAssignments.Valid {
		analytics.CompletedAssignments = int(completedAssignments.Int64)
	}
	if overdueAssignments.Valid {
		analytics.OverdueAssignments = int(overdueAssignments.Int64)
	}
	if analytics.TotalAssignments > 0 {
		analytics.CompletionRate = float64(analytics.CompletedAssignments) / float64(analytics.TotalAssignments)
	}
	if avgTime.Valid {
		analytics.AverageTimeSpent = int(math.Round(avgTime.Float64))
	}

	return analytics, nil
}
"""

content += "\n\n" + analytics

path.write_text(content, encoding="utf-8")
