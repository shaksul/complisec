from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go")
content = path.read_text(encoding="utf-8")

quiz_attempts = """// Quiz attempts ------------------------------------------------------------

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
"""

content += "\n\n" + quiz_attempts

path.write_text(content, encoding="utf-8")
