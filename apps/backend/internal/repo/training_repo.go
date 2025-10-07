package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/lib/pq"
)

type TrainingRepo struct {
	db DBInterface
}

func NewTrainingRepo(db DBInterface) *TrainingRepo {
	return &TrainingRepo{db: db}
}

// Materials ----------------------------------------------------------------

func (r *TrainingRepo) CreateMaterial(ctx context.Context, material Material) error {
	query := `
		INSERT INTO materials (
			id, tenant_id, title, description, uri, type, material_type,
			duration_minutes, tags, is_required, passing_score, attempts_limit,
			metadata, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13, $14
		)`

	metadataJSON := marshalJSON(material.Metadata)

	_, err := r.db.ExecContext(
		ctx,
		query,
		material.ID,
		material.TenantID,
		material.Title,
		material.Description,
		material.URI,
		material.Type,
		material.MaterialType,
		material.DurationMinutes,
		pq.Array(material.Tags),
		material.IsRequired,
		material.PassingScore,
		material.AttemptsLimit,
		metadataJSON,
		material.CreatedBy,
	)

	return err
}

func (r *TrainingRepo) GetMaterialByID(ctx context.Context, id string) (*Material, error) {
	query := `
		SELECT id, tenant_id, title, description, uri, type, material_type,
			duration_minutes, tags, is_required, passing_score, attempts_limit,
			metadata, created_by, created_at, updated_at
		FROM materials
		WHERE id = $1`

	var material Material
	var tags pq.StringArray
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&material.ID,
		&material.TenantID,
		&material.Title,
		&material.Description,
		&material.URI,
		&material.Type,
		&material.MaterialType,
		&material.DurationMinutes,
		(*pq.StringArray)(&tags),
		&material.IsRequired,
		&material.PassingScore,
		&material.AttemptsLimit,
		&metadataJSON,
		&material.CreatedBy,
		&material.CreatedAt,
		&material.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	material.Tags = cloneStringSlice([]string(tags))
	material.Metadata = unmarshalJSONMap(metadataJSON)

	return &material, nil
}

func (r *TrainingRepo) ListMaterials(ctx context.Context, tenantID string, filters map[string]interface{}) ([]Material, error) {
	query := `
		SELECT id, tenant_id, title, description, uri, type, material_type,
			duration_minutes, tags, is_required, passing_score, attempts_limit,
			metadata, created_by, created_at, updated_at
		FROM materials
		WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	argIdx := 2

	if v, ok := filters["material_type"]; ok {
		query += fmt.Sprintf(" AND material_type = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["is_required"]; ok {
		query += fmt.Sprintf(" AND is_required = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["search"]; ok {
		query += fmt.Sprintf(" AND (LOWER(title) LIKE $%d OR LOWER(COALESCE(description, '')) LIKE $%d)", argIdx, argIdx)
		args = append(args, patternForSearch(v))
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var materials []Material
	for rs.Next() {
		var material Material
		var tags pq.StringArray
		var metadataJSON []byte

		if err := rs.Scan(
			&material.ID,
			&material.TenantID,
			&material.Title,
			&material.Description,
			&material.URI,
			&material.Type,
			&material.MaterialType,
			&material.DurationMinutes,
			(*pq.StringArray)(&tags),
			&material.IsRequired,
			&material.PassingScore,
			&material.AttemptsLimit,
			&metadataJSON,
			&material.CreatedBy,
			&material.CreatedAt,
			&material.UpdatedAt,
		); err != nil {
			return nil, err
		}

		material.Tags = cloneStringSlice([]string(tags))
		material.Metadata = unmarshalJSONMap(metadataJSON)
		materials = append(materials, material)
	}

	return materials, rs.Err()
}

func (r *TrainingRepo) UpdateMaterial(ctx context.Context, material Material) error {
	query := `
		UPDATE materials SET
			tenant_id = $2,
			title = $3,
			description = $4,
			uri = $5,
			type = $6,
			material_type = $7,
			duration_minutes = $8,
			tags = $9,
			is_required = $10,
			passing_score = $11,
			attempts_limit = $12,
			metadata = $13,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	metadataJSON := marshalJSON(material.Metadata)

	_, err := r.db.ExecContext(
		ctx,
		query,
		material.ID,
		material.TenantID,
		material.Title,
		material.Description,
		material.URI,
		material.Type,
		material.MaterialType,
		material.DurationMinutes,
		pq.Array(material.Tags),
		material.IsRequired,
		material.PassingScore,
		material.AttemptsLimit,
		metadataJSON,
	)
	return err
}

func (r *TrainingRepo) DeleteMaterial(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM materials WHERE id = $1", id)
	return err
}

// Courses -------------------------------------------------------------------

func (r *TrainingRepo) CreateCourse(ctx context.Context, course TrainingCourse) error {
	query := `
		INSERT INTO training_courses (
			id, tenant_id, title, description, is_active, created_by
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		course.ID,
		course.TenantID,
		course.Title,
		course.Description,
		course.IsActive,
		course.CreatedBy,
	)
	return err
}

func (r *TrainingRepo) GetCourseByID(ctx context.Context, id string) (*TrainingCourse, error) {
	query := `
		SELECT id, tenant_id, title, description, is_active, created_by, created_at, updated_at
		FROM training_courses
		WHERE id = $1`

	var course TrainingCourse
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.TenantID,
		&course.Title,
		&course.Description,
		&course.IsActive,
		&course.CreatedBy,
		&course.CreatedAt,
		&course.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *TrainingRepo) ListCourses(ctx context.Context, tenantID string, filters map[string]interface{}) ([]TrainingCourse, error) {
	query := `
		SELECT id, tenant_id, title, description, is_active, created_by, created_at, updated_at
		FROM training_courses
		WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	argIdx := 2

	if v, ok := filters["is_active"]; ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["search"]; ok {
		query += fmt.Sprintf(" AND LOWER(title) LIKE $%d", argIdx)
		args = append(args, patternForSearch(v))
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var courses []TrainingCourse
	for rs.Next() {
		var course TrainingCourse
		if err := rs.Scan(
			&course.ID,
			&course.TenantID,
			&course.Title,
			&course.Description,
			&course.IsActive,
			&course.CreatedBy,
			&course.CreatedAt,
			&course.UpdatedAt,
		); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, rs.Err()
}

func (r *TrainingRepo) UpdateCourse(ctx context.Context, course TrainingCourse) error {
	query := `
		UPDATE training_courses SET
			tenant_id = $2,
			title = $3,
			description = $4,
			is_active = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		course.ID,
		course.TenantID,
		course.Title,
		course.Description,
		course.IsActive,
	)
	return err
}

func (r *TrainingRepo) DeleteCourse(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM training_courses WHERE id = $1", id)
	return err
}

func (r *TrainingRepo) AddMaterialToCourse(ctx context.Context, cm CourseMaterial) error {
	query := `
		INSERT INTO course_materials (id, course_id, material_id, order_index, is_required)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (course_id, material_id) DO UPDATE
		SET order_index = EXCLUDED.order_index,
			is_required = EXCLUDED.is_required`

	_, err := r.db.ExecContext(ctx, query,
		cm.ID,
		cm.CourseID,
		cm.MaterialID,
		cm.OrderIndex,
		cm.IsRequired,
	)
	return err
}

func (r *TrainingRepo) RemoveMaterialFromCourse(ctx context.Context, courseID, materialID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM course_materials WHERE course_id = $1 AND material_id = $2", courseID, materialID)
	return err
}

func (r *TrainingRepo) GetCourseMaterials(ctx context.Context, courseID string) ([]CourseMaterial, error) {
	query := `
		SELECT id, course_id, material_id, order_index, is_required, created_at
		FROM course_materials
		WHERE course_id = $1
		ORDER BY order_index ASC, created_at ASC`

	rs, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var materials []CourseMaterial
	for rs.Next() {
		var cm CourseMaterial
		if err := rs.Scan(
			&cm.ID,
			&cm.CourseID,
			&cm.MaterialID,
			&cm.OrderIndex,
			&cm.IsRequired,
			&cm.CreatedAt,
		); err != nil {
			return nil, err
		}
		materials = append(materials, cm)
	}

	return materials, rs.Err()
}

// Assignments ---------------------------------------------------------------

func (r *TrainingRepo) CreateAssignment(ctx context.Context, assignment TrainingAssignment) error {
	query := `
		INSERT INTO train_assignments (
			id, tenant_id, material_id, course_id, user_id, status, due_at, completed_at,
			assigned_by, priority, progress_percentage, time_spent_minutes,
			last_accessed_at, reminder_sent_at, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15
		)`

	metadataJSON := marshalJSON(assignment.Metadata)

	_, err := r.db.ExecContext(ctx, query,
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

func (r *TrainingRepo) GetAssignmentByID(ctx context.Context, id string) (*TrainingAssignment, error) {
	query := `
		SELECT id, tenant_id, material_id, course_id, user_id, status, due_at, completed_at,
			assigned_by, priority, progress_percentage, time_spent_minutes,
			last_accessed_at, reminder_sent_at, metadata, created_at
		FROM train_assignments
		WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	assignment, err := scanAssignment(row)
	if err != nil {
		return nil, err
	}
	return assignment, nil
}

func (r *TrainingRepo) GetUserAssignments(ctx context.Context, userID string, filters map[string]interface{}) ([]TrainingAssignment, error) {
	query := `
		SELECT id, tenant_id, material_id, course_id, user_id, status, due_at, completed_at,
			assigned_by, priority, progress_percentage, time_spent_minutes,
			last_accessed_at, reminder_sent_at, metadata, created_at
		FROM train_assignments
		WHERE user_id = $1`

	args := []interface{}{userID}
	argIdx := 2

	if v, ok := filters["tenant_id"]; ok {
		query += fmt.Sprintf(" AND tenant_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["status"]; ok {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["course_id"]; ok {
		query += fmt.Sprintf(" AND course_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["material_id"]; ok {
		query += fmt.Sprintf(" AND material_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["due_before"]; ok {
		query += fmt.Sprintf(" AND due_at <= $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["due_after"]; ok {
		query += fmt.Sprintf(" AND due_at >= $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["overdue_only"]; ok {
		if flag, ok := v.(bool); ok && flag {
			query += " AND due_at IS NOT NULL AND due_at < CURRENT_TIMESTAMP AND status <> 'completed'"
		}
	}

	query += " ORDER BY created_at DESC"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var assignments []TrainingAssignment
	for rs.Next() {
		assignment, err := scanAssignment(rs)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, *assignment)
	}

	return assignments, rs.Err()
}

func (r *TrainingRepo) UpdateAssignment(ctx context.Context, assignment TrainingAssignment) error {
	query := `
		UPDATE train_assignments SET
			tenant_id = $2,
			material_id = $3,
			course_id = $4,
			user_id = $5,
			status = $6,
			due_at = $7,
			completed_at = $8,
			assigned_by = $9,
			priority = $10,
			progress_percentage = $11,
			time_spent_minutes = $12,
			last_accessed_at = $13,
			reminder_sent_at = $14,
			metadata = $15
		WHERE id = $1`

	metadataJSON := marshalJSON(assignment.Metadata)

	_, err := r.db.ExecContext(ctx, query,
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

func (r *TrainingRepo) DeleteAssignment(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM train_assignments WHERE id = $1", id)
	return err
}

func (r *TrainingRepo) GetOverdueAssignments(ctx context.Context, tenantID string) ([]TrainingAssignment, error) {
	query := `
		SELECT id, tenant_id, material_id, course_id, user_id, status, due_at, completed_at,
			assigned_by, priority, progress_percentage, time_spent_minutes,
			last_accessed_at, reminder_sent_at, metadata, created_at
		FROM train_assignments
		WHERE tenant_id = $1
			AND status <> 'completed'
			AND due_at IS NOT NULL
			AND due_at < CURRENT_TIMESTAMP
		ORDER BY due_at ASC`

	rs, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var assignments []TrainingAssignment
	for rs.Next() {
		assignment, err := scanAssignment(rs)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, *assignment)
	}

	return assignments, rs.Err()
}

func (r *TrainingRepo) GetUpcomingDeadlines(ctx context.Context, tenantID string, days int) ([]TrainingAssignment, error) {
	query := `
		SELECT id, tenant_id, material_id, course_id, user_id, status, due_at, completed_at,
			assigned_by, priority, progress_percentage, time_spent_minutes,
			last_accessed_at, reminder_sent_at, metadata, created_at
		FROM train_assignments
		WHERE tenant_id = $1
			AND status <> 'completed'
			AND due_at IS NOT NULL
			AND due_at BETWEEN CURRENT_TIMESTAMP AND CURRENT_TIMESTAMP + make_interval(days => $2)
		ORDER BY due_at ASC`

	rs, err := r.db.QueryContext(ctx, query, tenantID, days)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var assignments []TrainingAssignment
	for rs.Next() {
		assignment, err := scanAssignment(rs)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, *assignment)
	}

	return assignments, rs.Err()
}

// Progress -----------------------------------------------------------------

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

// Quiz questions -----------------------------------------------------------

func (r *TrainingRepo) CreateQuizQuestion(ctx context.Context, question QuizQuestion) error {
	query := `
		INSERT INTO quiz_questions (
			id, material_id, text, options_json, correct_index, question_type,
			points, explanation, order_index
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	optionsJSON := marshalJSON(question.OptionsJSON)

	_, err := r.db.ExecContext(ctx, query,
		question.ID,
		question.MaterialID,
		question.Text,
		optionsJSON,
		question.CorrectIndex,
		question.QuestionType,
		question.Points,
		question.Explanation,
		question.OrderIndex,
	)
	return err
}

func (r *TrainingRepo) GetQuizQuestionByID(ctx context.Context, id string) (*QuizQuestion, error) {
	query := `
		SELECT id, material_id, text, options_json, correct_index, question_type,
			points, explanation, order_index, created_at
		FROM quiz_questions
		WHERE id = $1`

	var question QuizQuestion
	var optionsJSON []byte
	var explanation sql.NullString

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&question.ID,
		&question.MaterialID,
		&question.Text,
		&optionsJSON,
		&question.CorrectIndex,
		&question.QuestionType,
		&question.Points,
		&explanation,
		&question.OrderIndex,
		&question.CreatedAt,
	); err != nil {
		return nil, err
	}

	question.OptionsJSON = unmarshalJSONMap(optionsJSON)
	question.Explanation = stringPointer(explanation)

	return &question, nil
}

func (r *TrainingRepo) ListQuizQuestions(ctx context.Context, materialID string) ([]QuizQuestion, error) {
	query := `
		SELECT id, material_id, text, options_json, correct_index, question_type,
			points, explanation, order_index, created_at
		FROM quiz_questions
		WHERE material_id = $1
		ORDER BY order_index ASC, created_at ASC`

	rs, err := r.db.QueryContext(ctx, query, materialID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var questions []QuizQuestion
	for rs.Next() {
		var question QuizQuestion
		var optionsJSON []byte
		var explanation sql.NullString

		if err := rs.Scan(
			&question.ID,
			&question.MaterialID,
			&question.Text,
			&optionsJSON,
			&question.CorrectIndex,
			&question.QuestionType,
			&question.Points,
			&explanation,
			&question.OrderIndex,
			&question.CreatedAt,
		); err != nil {
			return nil, err
		}

		question.OptionsJSON = unmarshalJSONMap(optionsJSON)
		question.Explanation = stringPointer(explanation)
		questions = append(questions, question)
	}

	return questions, rs.Err()
}

func (r *TrainingRepo) UpdateQuizQuestion(ctx context.Context, question QuizQuestion) error {
	query := `
		UPDATE quiz_questions SET
			text = $2,
			options_json = $3,
			correct_index = $4,
			question_type = $5,
			points = $6,
			explanation = $7,
			order_index = $8
		WHERE id = $1`

	optionsJSON := marshalJSON(question.OptionsJSON)

	_, err := r.db.ExecContext(ctx, query,
		question.ID,
		question.Text,
		optionsJSON,
		question.CorrectIndex,
		question.QuestionType,
		question.Points,
		question.Explanation,
		question.OrderIndex,
	)
	return err
}

func (r *TrainingRepo) DeleteQuizQuestion(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM quiz_questions WHERE id = $1", id)
	return err
}

// Quiz attempts ------------------------------------------------------------

func (r *TrainingRepo) CreateQuizAttempt(ctx context.Context, attempt QuizAttempt) error {
	query := `
		INSERT INTO quiz_attempts (
			id, user_id, material_id, assignment_id, score, max_score, passed,
			answers_json, time_spent_minutes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	answersJSON := marshalJSON(attempt.AnswersJSON)

	_, err := r.db.ExecContext(ctx, query,
		attempt.ID,
		attempt.UserID,
		attempt.MaterialID,
		attempt.AssignmentID,
		attempt.Score,
		attempt.MaxScore,
		attempt.Passed,
		answersJSON,
		attempt.TimeSpentMinutes,
	)
	return err
}

func (r *TrainingRepo) GetQuizAttemptByID(ctx context.Context, id string) (*QuizAttempt, error) {
	query := `
		SELECT id, user_id, material_id, assignment_id, score, max_score, passed,
			answers_json, time_spent_minutes, attempted_at
		FROM quiz_attempts
		WHERE id = $1`

	var attempt QuizAttempt
	var answersJSON []byte
	var assignmentID sql.NullString
	var maxScore sql.NullInt64
	var timeSpent sql.NullInt64

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.MaterialID,
		&assignmentID,
		&attempt.Score,
		&maxScore,
		&attempt.Passed,
		&answersJSON,
		&timeSpent,
		&attempt.AttemptedAt,
	); err != nil {
		return nil, err
	}

	attempt.AssignmentID = stringPointer(assignmentID)
	attempt.MaxScore = intPointer(maxScore)
	attempt.TimeSpentMinutes = intPointer(timeSpent)
	attempt.AnswersJSON = unmarshalJSONMap(answersJSON)

	return &attempt, nil
}

func (r *TrainingRepo) GetQuizAttempts(ctx context.Context, assignmentID, materialID string) ([]QuizAttempt, error) {
	query := `
		SELECT id, user_id, material_id, assignment_id, score, max_score, passed,
			answers_json, time_spent_minutes, attempted_at
		FROM quiz_attempts
		WHERE material_id = $1`

	args := []interface{}{materialID}
	argIdx := 2

	if assignmentID != "" {
		query += fmt.Sprintf(" AND assignment_id = $%d", argIdx)
		args = append(args, assignmentID)
		argIdx++
	}

	query += " ORDER BY attempted_at DESC"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var attempts []QuizAttempt
	for rs.Next() {
		var attempt QuizAttempt
		var answersJSON []byte
		var assignment sql.NullString
		var maxScore sql.NullInt64
		var timeSpent sql.NullInt64

		if err := rs.Scan(
			&attempt.ID,
			&attempt.UserID,
			&attempt.MaterialID,
			&assignment,
			&attempt.Score,
			&maxScore,
			&attempt.Passed,
			&answersJSON,
			&timeSpent,
			&attempt.AttemptedAt,
		); err != nil {
			return nil, err
		}

		attempt.AssignmentID = stringPointer(assignment)
		attempt.MaxScore = intPointer(maxScore)
		attempt.TimeSpentMinutes = intPointer(timeSpent)
		attempt.AnswersJSON = unmarshalJSONMap(answersJSON)
		attempts = append(attempts, attempt)
	}

	return attempts, rs.Err()
}

// Certificates -------------------------------------------------------------

func (r *TrainingRepo) CreateCertificate(ctx context.Context, certificate Certificate) error {
	query := `
		INSERT INTO certificates (
			id, tenant_id, assignment_id, user_id, material_id, course_id,
			certificate_number, issued_at, expires_at, is_valid, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	metadataJSON := marshalJSON(certificate.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		certificate.ID,
		certificate.TenantID,
		certificate.AssignmentID,
		certificate.UserID,
		certificate.MaterialID,
		certificate.CourseID,
		certificate.CertificateNumber,
		certificate.IssuedAt,
		certificate.ExpiresAt,
		certificate.IsValid,
		metadataJSON,
	)
	return err
}

func (r *TrainingRepo) GetCertificateByID(ctx context.Context, id string) (*Certificate, error) {
	return r.getCertificate(ctx, "id = $1", id)
}

func (r *TrainingRepo) GetCertificateByNumber(ctx context.Context, certificateNumber string) (*Certificate, error) {
	return r.getCertificate(ctx, "certificate_number = $1", certificateNumber)
}

func (r *TrainingRepo) GetUserCertificates(ctx context.Context, userID string, filters map[string]interface{}) ([]Certificate, error) {
	query := `
		SELECT id, tenant_id, assignment_id, user_id, material_id, course_id,
			certificate_number, issued_at, expires_at, is_valid, metadata, created_at
		FROM certificates
		WHERE user_id = $1`

	args := []interface{}{userID}
	argIdx := 2

	if v, ok := filters["tenant_id"]; ok {
		query += fmt.Sprintf(" AND tenant_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["course_id"]; ok {
		query += fmt.Sprintf(" AND course_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["material_id"]; ok {
		query += fmt.Sprintf(" AND material_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["is_valid"]; ok {
		query += fmt.Sprintf(" AND is_valid = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	query += " ORDER BY issued_at DESC"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var certificates []Certificate
	for rs.Next() {
		certificate, err := scanCertificate(rs)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, *certificate)
	}

	return certificates, rs.Err()
}

// Notifications ------------------------------------------------------------

func (r *TrainingRepo) CreateNotification(ctx context.Context, notification TrainingNotification) error {
	query := `
		INSERT INTO training_notifications (
			id, tenant_id, assignment_id, user_id, type, title, message,
			is_read, sent_at, read_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query,
		notification.ID,
		notification.TenantID,
		notification.AssignmentID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.IsRead,
		notification.SentAt,
		notification.ReadAt,
	)
	return err
}

func (r *TrainingRepo) GetUserNotifications(ctx context.Context, userID string, unreadOnly bool) ([]TrainingNotification, error) {
	query := `
		SELECT id, tenant_id, assignment_id, user_id, type, title, message,
			is_read, sent_at, read_at
		FROM training_notifications
		WHERE user_id = $1`

	args := []interface{}{userID}

	if unreadOnly {
		query += " AND is_read = false"
	}

	query += " ORDER BY sent_at DESC"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var notifications []TrainingNotification
	for rs.Next() {
		notification, err := scanNotification(rs)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, *notification)
	}

	return notifications, rs.Err()
}

func (r *TrainingRepo) MarkNotificationAsRead(ctx context.Context, notificationID, userID string) error {
	query := `
		UPDATE training_notifications
		SET is_read = true, read_at = COALESCE(read_at, CURRENT_TIMESTAMP)
		WHERE id = $1 AND user_id = $2`

	_, err := r.db.ExecContext(ctx, query, notificationID, userID)
	return err
}

// Analytics ----------------------------------------------------------------

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

// Role assignments ---------------------------------------------------------

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

// Helpers ------------------------------------------------------------------

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanAssignment(scanner rowScanner) (*TrainingAssignment, error) {
	var assignment TrainingAssignment
	var materialID sql.NullString
	var courseID sql.NullString
	var dueAt sql.NullTime
	var completedAt sql.NullTime
	var assignedBy sql.NullString
	var lastAccessed sql.NullTime
	var reminderSent sql.NullTime
	var metadataJSON []byte

	if err := scanner.Scan(
		&assignment.ID,
		&assignment.TenantID,
		&materialID,
		&courseID,
		&assignment.UserID,
		&assignment.Status,
		&dueAt,
		&completedAt,
		&assignedBy,
		&assignment.Priority,
		&assignment.ProgressPercentage,
		&assignment.TimeSpentMinutes,
		&lastAccessed,
		&reminderSent,
		&metadataJSON,
		&assignment.CreatedAt,
	); err != nil {
		return nil, err
	}

	assignment.MaterialID = stringPointer(materialID)
	assignment.CourseID = stringPointer(courseID)
	assignment.DueAt = timePointer(dueAt)
	assignment.CompletedAt = timePointer(completedAt)
	assignment.AssignedBy = stringPointer(assignedBy)
	assignment.LastAccessedAt = timePointer(lastAccessed)
	assignment.ReminderSentAt = timePointer(reminderSent)
	assignment.Metadata = unmarshalJSONMap(metadataJSON)

	return &assignment, nil
}

func scanCertificate(scanner rowScanner) (*Certificate, error) {
	var certificate Certificate
	var materialID sql.NullString
	var courseID sql.NullString
	var expiresAt sql.NullTime
	var metadataJSON []byte

	if err := scanner.Scan(
		&certificate.ID,
		&certificate.TenantID,
		&certificate.AssignmentID,
		&certificate.UserID,
		&materialID,
		&courseID,
		&certificate.CertificateNumber,
		&certificate.IssuedAt,
		&expiresAt,
		&certificate.IsValid,
		&metadataJSON,
		&certificate.CreatedAt,
	); err != nil {
		return nil, err
	}

	certificate.MaterialID = stringPointer(materialID)
	certificate.CourseID = stringPointer(courseID)
	certificate.ExpiresAt = timePointer(expiresAt)
	certificate.Metadata = unmarshalJSONMap(metadataJSON)

	return &certificate, nil
}

func scanNotification(scanner rowScanner) (*TrainingNotification, error) {
	var notification TrainingNotification
	var readAt sql.NullTime

	if err := scanner.Scan(
		&notification.ID,
		&notification.TenantID,
		&notification.AssignmentID,
		&notification.UserID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&notification.IsRead,
		&notification.SentAt,
		&readAt,
	); err != nil {
		return nil, err
	}

	notification.ReadAt = timePointer(readAt)
	return &notification, nil
}

func (r *TrainingRepo) getCertificate(ctx context.Context, predicate string, arg interface{}) (*Certificate, error) {
	query := fmt.Sprintf(
		"SELECT id, tenant_id, assignment_id, user_id, material_id, course_id, certificate_number, issued_at, expires_at, is_valid, metadata, created_at FROM certificates WHERE %s",
		predicate,
	)

	row := r.db.QueryRowContext(ctx, query, arg)
	certificate, err := scanCertificate(row)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}

func marshalJSON(data map[string]any) []byte {
	if len(data) == 0 {
		return nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return b
}

func unmarshalJSONMap(data []byte) map[string]any {
	if len(data) == 0 {
		return nil
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}

func patternForSearch(value interface{}) string {
	text := strings.ToLower(fmt.Sprint(value))
	return "%" + text + "%"
}

func cloneStringSlice(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func stringPointer(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

func timePointer(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	t := nt.Time
	return &t
}

func intPointer(ni sql.NullInt64) *int {
	if !ni.Valid {
		return nil
	}
	v := int(ni.Int64)
	return &v
}
