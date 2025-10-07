from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/domain/training_service_impl.go")
content = path.read_text(encoding="utf-8")

anchor = "func (s *TrainingService) SubmitQuizAttempt"

quiz_section = """// Quiz questions -----------------------------------------------------------

func (s *TrainingService) CreateQuizQuestion(ctx context.Context, materialID string, req dto.CreateQuizQuestionRequest, createdBy string) (*repo.QuizQuestion, error) {
	if materialID == "" {
		return nil, errors.New("material_id is required")
	}

	if _, err := s.trainingRepo.GetMaterialByID(ctx, materialID); err != nil {
		return nil, fmt.Errorf("failed to fetch material: %w", err)
	}

	question := &repo.QuizQuestion{
		ID:           generateID(),
		MaterialID:   materialID,
		Text:         req.Text,
		OptionsJSON:  copyMetadata(req.OptionsJSON),
		CorrectIndex: req.CorrectIndex,
		QuestionType: req.QuestionType,
		Points:       req.Points,
		Explanation:  req.Explanation,
		OrderIndex:   req.OrderIndex,
	}

	if err := s.trainingRepo.CreateQuizQuestion(ctx, *question); err != nil {
		return nil, fmt.Errorf("failed to create quiz question: %w", err)
	}

	return question, nil
}

func (s *TrainingService) GetQuizQuestion(ctx context.Context, id string) (*repo.QuizQuestion, error) {
	question, err := s.trainingRepo.GetQuizQuestionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz question: %w", err)
	}

	return question, nil
}

func (s *TrainingService) ListQuizQuestions(ctx context.Context, materialID string) ([]repo.QuizQuestion, error) {
	questions, err := s.trainingRepo.ListQuizQuestions(ctx, materialID)
	if err != nil {
		return nil, fmt.Errorf("failed to list quiz questions: %w", err)
	}

	return questions, nil
}

func (s *TrainingService) UpdateQuizQuestion(ctx context.Context, id string, req dto.UpdateQuizQuestionRequest, updatedBy string) error {
	question, err := s.trainingRepo.GetQuizQuestionByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get quiz question: %w", err)
	}

	if req.Text != nil {
		question.Text = *req.Text
	}
	if req.OptionsJSON != nil {
		question.OptionsJSON = copyMetadata(req.OptionsJSON)
	}
	if req.CorrectIndex != nil {
		question.CorrectIndex = *req.CorrectIndex
	}
	if req.QuestionType != nil {
		question.QuestionType = *req.QuestionType
	}
	if req.Points != nil {
		question.Points = *req.Points
	}
	if req.Explanation != nil {
		question.Explanation = req.Explanation
	}
	if req.OrderIndex != nil {
		question.OrderIndex = *req.OrderIndex
	}

	if err := s.trainingRepo.UpdateQuizQuestion(ctx, *question); err != nil {
		return fmt.Errorf("failed to update quiz question: %w", err)
	}

	return nil
}

func (s *TrainingService) DeleteQuizQuestion(ctx context.Context, id string, deletedBy string) error {
	if err := s.trainingRepo.DeleteQuizQuestion(ctx, id); err != nil {
		return fmt.Errorf("failed to delete quiz question: %w", err)
	}

	return nil
}

"""

if anchor not in content:
    raise SystemExit("anchor not found")

content = content.replace(anchor, quiz_section + "\n" + anchor, 1)

path.write_text(content, encoding="utf-8")
