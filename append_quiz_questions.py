from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go")
content = path.read_text(encoding="utf-8")

quiz_questions = """// Quiz questions -----------------------------------------------------------

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
"""

content += "\n\n" + quiz_questions

path.write_text(content, encoding="utf-8")
