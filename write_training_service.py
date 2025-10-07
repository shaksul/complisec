from pathlib import Path

go_code = """package domain

import (
	\"context\"
	\"database/sql\"
	\"errors\"
	\"fmt\"
	\"math\"
	\"strconv\"
	\"strings\"
	\"time\"

	\"github.com/google/uuid\"

	\"risknexus/backend/internal/dto\"
	\"risknexus/backend/internal/repo\"
)

const (
	assignmentStatusAssigned   = \"assigned\"
	assignmentStatusInProgress = \"in_progress\"
	assignmentStatusCompleted  = \"completed\"
	assignmentStatusOverdue    = \"overdue\"

	defaultPriorityValue = \"normal\"
	reminderWindowDays   = 3
)

type TrainingService struct {
	trainingRepo repo.TrainingRepoInterface
}

func NewTrainingService(trainingRepo repo.TrainingRepoInterface) *TrainingService {
	return &TrainingService{
		trainingRepo: trainingRepo,
	}
}

// Materials ----------------------------------------------------------------

func (s *TrainingService) CreateMaterial(ctx context.Context, tenantID string, req dto.CreateMaterialRequest, createdBy string) (*repo.Material, error) {
	material := &repo.Material{
		ID:              generateID(),
		TenantID:        tenantID,
		Title:           req.Title,
		Description:     req.Description,
		URI:             req.URI,
		Type:            req.Type,
		MaterialType:    req.MaterialType,
		DurationMinutes: req.DurationMinutes,
		Tags:            req.Tags,
		IsRequired:      req.IsRequired,
		PassingScore:    req.PassingScore,
		AttemptsLimit:   req.AttemptsLimit,
		Metadata:        copyMetadata(req.Metadata),
		CreatedBy:       stringPtr(createdBy),
	}

	if err := s.trainingRepo.CreateMaterial(ctx, *material); err != nil {
		return nil, fmt.Errorf("failed to create material: %w", err)
	}

	return material, nil
}

func (s *TrainingService) GetMaterial(ctx context.Context, id string) (*repo.Material, error) {
	material, err := s.trainingRepo.GetMaterialByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get material: %w", err)
	}

	return material, nil
}

func (s *TrainingService) ListMaterials(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.Material, error) {
	materials, err := s.trainingRepo.ListMaterials(ctx, tenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list materials: %w", err)
	}

	return materials, nil
}

func (s *TrainingService) UpdateMaterial(ctx context.Context, id string, req dto.UpdateMaterialRequest, updatedBy string) error {
	material, err := s.trainingRepo.GetMaterialByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get material: %w", err)
	}

	if req.Title != nil {
		material.Title = *req.Title
	}
	if req.Description != nil {
		material.Description = req.Description
	}
	if req.URI != nil {
		material.URI = *req.URI
	}
	if req.Type != nil {
		material.Type = *req.Type
	}
	if req.MaterialType != nil {
		material.MaterialType = *req.MaterialType
	}
	if req.DurationMinutes != nil {
		material.DurationMinutes = req.DurationMinutes
	}
	if req.Tags != nil {
		material.Tags = req.Tags
	}
	if req.IsRequired != nil {
		material.IsRequired = *req.IsRequired
	}
	if req.PassingScore != nil {
		material.PassingScore = *req.PassingScore
	}
	if req.AttemptsLimit != nil {
		material.AttemptsLimit = req.AttemptsLimit
	}
	if req.Metadata != nil {
		material.Metadata = copyMetadata(req.Metadata)
	}

	if err := s.trainingRepo.UpdateMaterial(ctx, *material); err != nil {
		return fmt.Errorf("failed to update material: %w", err)
	}

	return nil
}

func (s *TrainingService) DeleteMaterial(ctx context.Context, id string, deletedBy string) error {
	if err := s.trainingRepo.DeleteMaterial(ctx, id); err != nil {
		return fmt.Errorf("failed to delete material: %w", err)
	}

	return nil
}

// Courses -------------------------------------------------------------------

func (s *TrainingService) CreateCourse(ctx context.Context, tenantID string, req dto.CreateCourseRequest, createdBy string) (*repo.TrainingCourse, error) {
	course := &repo.TrainingCourse{
		ID:          generateID(),
		TenantID:    tenantID,
		Title:       req.Title,
		Description: req.Description,
		IsActive:    req.IsActive,
		CreatedBy:   stringPtr(createdBy),
	}

	if err := s.trainingRepo.CreateCourse(ctx, *course); err != nil {
		return nil, fmt.Errorf("failed to create course: %w", err)
	}

	return course, nil
}

func (s *TrainingService) GetCourse(ctx context.Context, id string) (*repo.TrainingCourse, error) {
	course, err := s.trainingRepo.GetCourseByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	return course, nil
}

func (s *TrainingService) ListCourses(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.TrainingCourse, error) {
	courses, err := s.trainingRepo.ListCourses(ctx, tenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list courses: %w", err)
	}

	return courses, nil
}

func (s *TrainingService) UpdateCourse(ctx context.Context, id string, req dto.UpdateCourseRequest, updatedBy string) error {
	course, err := s.trainingRepo.GetCourseByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get course: %w", err)
	}

	if req.Title != nil {
		course.Title = *req.Title
	}
	if req.Description != nil {
		course.Description = req.Description
	}
	if req.IsActive != nil {
		course.IsActive = *req.IsActive
	}

	if err := s.trainingRepo.UpdateCourse(ctx, *course); err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	return nil
}

func (s *TrainingService) DeleteCourse(ctx context.Context, id string, deletedBy string) error {
	if err := s.trainingRepo.DeleteCourse(ctx, id); err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	return nil
}

func (s *TrainingService) AddMaterialToCourse(ctx context.Context, courseID, materialID string, req dto.CourseMaterialRequest, addedBy string) error {
	courseMaterial := repo.CourseMaterial{
		ID:         generateID(),
		CourseID:   courseID,
		MaterialID: materialID,
		OrderIndex: req.OrderIndex,
		IsRequired: req.IsRequired,
	}

	if err := s.trainingRepo.AddMaterialToCourse(ctx, courseMaterial); err != nil {
		return fmt.Errorf("failed to add material to course: %w", err)
	}

	return nil
}

func (s *TrainingService) RemoveMaterialFromCourse(ctx context.Context, courseID, materialID string, removedBy string) error {
	if err := s.trainingRepo.RemoveMaterialFromCourse(ctx, courseID, materialID); err != nil {
		return fmt.Errorf("failed to remove material from course: %w", err)
	}

	return nil
}

func (s *TrainingService) GetCourseMaterials(ctx context.Context, courseID string) ([]repo.CourseMaterial, error) {
	materials, err := s.trainingRepo.GetCourseMaterials(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course materials: %w", err)
	}

	return materials, nil
}

// Assignments ---------------------------------------------------------------

func (s *TrainingService) AssignMaterial(ctx context.Context, tenantID string, req dto.AssignMaterialRequest, assignedBy string) (*repo.TrainingAssignment, error) {
	if req.MaterialID == \"\" {
		return nil, errors.New("material_id is required")
	}
	if len(req.UserIDs) == 0 {
		return nil, errors.New("at least one user_id is required")
	}

	if _, err := s.trainingRepo.GetMaterialByID(ctx, req.MaterialID); err != nil {
		return nil, fmt.Errorf("failed to fetch material: %w", err)
	}

	var created *repo.TrainingAssignment
	for _, userID := range req.UserIDs {
		assignment := repo.TrainingAssignment{
			ID:                 generateID(),
			TenantID:           tenantID,
			UserID:             userID,
			Status:             assignmentStatusAssigned,
			DueAt:              req.DueAt,
			AssignedBy:         stringPtr(assignedBy),
			Priority:           ensurePriority(req.Priority),
			ProgressPercentage: 0,
			TimeSpentMinutes:   0,
			Metadata:           copyMetadata(req.Metadata),
		}

		materialID := req.MaterialID
		assignment.MaterialID = &materialID

		if req.DueAt != nil && req.DueAt.Before(time.Now()) {
			assignment.Status = assignmentStatusOverdue
		}

		if err := s.trainingRepo.CreateAssignment(ctx, assignment); err != nil {
			return nil, fmt.Errorf("failed to create assignment: %w", err)
		}

		if created == nil {
			copy := assignment
			created = &copy
		}
	}

	return created, nil
}

func (s *TrainingService) AssignCourse(ctx context.Context, tenantID string, req dto.AssignCourseRequest, assignedBy string) (*repo.TrainingAssignment, error) {
	if req.CourseID == \"\" {
		return nil, errors.New("course_id is required")
	}
	if len(req.UserIDs) == 0 {
		return nil, errors.New("at least one user_id is required")
	}

	if _, err := s.trainingRepo.GetCourseByID(ctx, req.CourseID); err != nil {
		return nil, fmt.Errorf("failed to fetch course: %w", err)
	}

	var created *repo.TrainingAssignment
	for _, userID := range req.UserIDs {
		assignment := repo.TrainingAssignment{
			ID:                 generateID(),
			TenantID:           tenantID,
			UserID:             userID,
			Status:             assignmentStatusAssigned,
			DueAt:              req.DueAt,
			AssignedBy:         stringPtr(assignedBy),
			Priority:           ensurePriority(req.Priority),
			ProgressPercentage: 0,
			TimeSpentMinutes:   0,
			Metadata:           copyMetadata(req.Metadata),
		}

		courseID := req.CourseID
		assignment.CourseID = &courseID

		if req.DueAt != nil && req.DueAt.Before(time.Now()) {
			assignment.Status = assignmentStatusOverdue
		}

		if err := s.trainingRepo.CreateAssignment(ctx, assignment); err != nil {
			return nil, fmt.Errorf("failed to create course assignment: %w", err)
		}

		if created == nil {
			copy := assignment
			created = &copy
		}
	}

	return created, nil
}

func (s *TrainingService) AssignToRole(ctx context.Context, tenantID string, req dto.AssignToRoleRequest, assignedBy string) error {
	if req.MaterialID == nil && req.CourseID == nil {
		return errors.New("either material_id or course_id must be provided")
	}

	assignment := repo.RoleTrainingAssignment{
		ID:        generateID(),
		TenantID:  tenantID,
		RoleID:    req.RoleID,
		IsRequired: req.IsRequired,
		DueDays:   req.DueDays,
		AssignedBy: stringPtr(assignedBy),
	}

	if req.MaterialID != nil {
		assignment.MaterialID = req.MaterialID
	}
	if req.CourseID != nil {
		assignment.CourseID = req.CourseID
	}

	if err := s.trainingRepo.CreateRoleAssignment(ctx, assignment); err != nil {
		return fmt.Errorf("failed to create role assignment: %w", err)
	}

	return nil
}

func (s *TrainingService) GetUserAssignments(ctx context.Context, userID string, filters map[string]interface{}) ([]repo.TrainingAssignment, error) {
	assignments, err := s.trainingRepo.GetUserAssignments(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get user assignments: %w", err)
	}

	return assignments, nil
}

func (s *TrainingService) GetAssignment(ctx context.Context, id string) (*repo.TrainingAssignment, error) {
	assignment, err := s.trainingRepo.GetAssignmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment: %w", err)
	}

	return assignment, nil
}

func (s *TrainingService) UpdateAssignment(ctx context.Context, id string, req dto.UpdateAssignmentRequest, updatedBy string) error {
	assignment, err := s.trainingRepo.GetAssignmentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	if req.Status != nil {
		assignment.Status = *req.Status
		if *req.Status == assignmentStatusCompleted {
			if assignment.CompletedAt == nil {
				now := time.Now()
				assignment.CompletedAt = &now
			}
			assignment.ProgressPercentage = 100
		} else if *req.Status == assignmentStatusAssigned {
			assignment.CompletedAt = nil
			assignment.ProgressPercentage = 0
		}
	}
	if req.DueAt != nil {
		assignment.DueAt = req.DueAt
	}
	if req.Priority != nil {
		assignment.Priority = ensurePriority(*req.Priority)
	}
	if req.Metadata != nil {
		assignment.Metadata = copyMetadata(req.Metadata)
	}

	assignment.Priority = ensurePriority(assignment.Priority)
	if assignment.Status != assignmentStatusCompleted && assignment.DueAt != nil && assignment.DueAt.Before(time.Now()) {
		assignment.Status = assignmentStatusOverdue
	}

	if err := s.trainingRepo.UpdateAssignment(ctx, *assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	return nil
}

func (s *TrainingService) DeleteAssignment(ctx context.Context, id string, deletedBy string) error {
	if err := s.trainingRepo.DeleteAssignment(ctx, id); err != nil {
		return fmt.Errorf("failed to delete assignment: %w", err)
	}

	return nil
}

func (s *TrainingService) UpdateProgress(ctx context.Context, assignmentID, materialID string, req dto.UpdateProgressRequest, updatedBy string) error {
	assignment, err := s.trainingRepo.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	existing, err := s.trainingRepo.GetProgressByAssignmentAndMaterial(ctx, assignmentID, materialID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to get progress: %w", err)
	}

	completedAt := req.CompletedAt
	if completedAt == nil && req.ProgressPercentage >= 100 {
		now := time.Now()
		completedAt = &now
	}

	if err == nil && existing != nil {
		existing.ProgressPercentage = req.ProgressPercentage
		existing.TimeSpentMinutes = req.TimeSpentMinutes
		existing.LastPosition = req.LastPosition
		existing.CompletedAt = completedAt
		if err := s.trainingRepo.UpdateProgress(ctx, *existing); err != nil {
			return fmt.Errorf("failed to update progress: %w", err)
		}
	} else {
		progress := repo.TrainingProgress{
			ID:                 generateID(),
			AssignmentID:       assignmentID,
			MaterialID:         materialID,
			ProgressPercentage: req.ProgressPercentage,
			TimeSpentMinutes:   req.TimeSpentMinutes,
			LastPosition:       req.LastPosition,
			CompletedAt:        completedAt,
		}

		if err := s.trainingRepo.CreateProgress(ctx, progress); err != nil {
			return fmt.Errorf("failed to create progress: %w", err)
		}
	}

	return s.updateAssignmentState(ctx, assignment, req.ProgressPercentage, req.TimeSpentMinutes, completedAt)
}

func (s *TrainingService) GetProgress(ctx context.Context, assignmentID string) ([]repo.TrainingProgress, error) {
	progress, err := s.trainingRepo.GetProgressByAssignment(ctx, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	return progress, nil
}

func (s *TrainingService) MarkAsCompleted(ctx context.Context, assignmentID, materialID string, completedBy string) error {
	assignment, err := s.trainingRepo.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	now := time.Now()

	if err := s.updateAssignmentState(ctx, assignment, 100, assignment.TimeSpentMinutes, &now); err != nil {
		return err
	}

	return nil
}

func (s *TrainingService) SubmitQuizAttempt(ctx context.Context, assignmentID, materialID string, req dto.SubmitQuizAttemptRequest, submittedBy string) (*repo.QuizAttempt, error) {
	questions, err := s.trainingRepo.ListQuizQuestions(ctx, materialID)
	if err != nil {
		return nil, fmt.Errorf("failed to list quiz questions: %w", err)
	}

	score, maxScore := calculateQuizScore(questions, req.AnswersJSON)
	percentage := percentFromScore(score, maxScore)

	passed := false
	material, err := s.trainingRepo.GetMaterialByID(ctx, materialID)
	if err == nil && maxScore > 0 {
		passed = percentage >= material.PassingScore
	} else if maxScore > 0 {
		passed = score == maxScore
	}

	attempt := &repo.QuizAttempt{
		ID:          generateID(),
		UserID:      submittedBy,
		MaterialID:  materialID,
		Score:       score,
		Passed:      passed,
		AnswersJSON: req.AnswersJSON,
		AttemptedAt: time.Now(),
	}

	if maxScore > 0 {
		attempt.MaxScore = intPtr(maxScore)
	}
	if assignmentID != "" {
		attempt.AssignmentID = &assignmentID
	}
	if req.TimeSpentMinutes != nil {
		attempt.TimeSpentMinutes = req.TimeSpentMinutes
	}

	if err := s.trainingRepo.CreateQuizAttempt(ctx, *attempt); err != nil {
		return nil, fmt.Errorf("failed to create quiz attempt: %w", err)
	}

	if assignmentID != "" {
		if assignment, err := s.trainingRepo.GetAssignmentByID(ctx, assignmentID); err == nil {
			timeSpent := assignment.TimeSpentMinutes
			if req.TimeSpentMinutes != nil {
				timeSpent += *req.TimeSpentMinutes
			}
			var completedAt *time.Time
			if passed {
				now := time.Now()
				completedAt = &now
			}
			_ = s.updateAssignmentState(ctx, assignment, percentage, timeSpent, completedAt)
		}
	}

	return attempt, nil
}

func (s *TrainingService) GetQuizAttempts(ctx context.Context, assignmentID, materialID string) ([]repo.QuizAttempt, error) {
	attempts, err := s.trainingRepo.GetQuizAttempts(ctx, assignmentID, materialID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz attempts: %w", err)
	}

	return attempts, nil
}

func (s *TrainingService) GetQuizAttempt(ctx context.Context, id string) (*repo.QuizAttempt, error) {
	attempt, err := s.trainingRepo.GetQuizAttemptByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz attempt: %w", err)
	}

	return attempt, nil
}

// Certificates -------------------------------------------------------------

func (s *TrainingService) GenerateCertificate(ctx context.Context, assignmentID string, generatedBy string) (*repo.Certificate, error) {
	assignment, err := s.trainingRepo.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment: %w", err)
	}
	if assignment.Status != assignmentStatusCompleted {
		return nil, errors.New("assignment is not completed")
	}

	certificateNumber := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), \"-\", \"\"))
	certificate := &repo.Certificate{
		ID:                generateID(),
		TenantID:          assignment.TenantID,
		AssignmentID:      assignment.ID,
		UserID:            assignment.UserID,
		MaterialID:        assignment.MaterialID,
		CourseID:          assignment.CourseID,
		CertificateNumber: fmt.Sprintf("CERT-%s", certificateNumber[:8]),
		IssuedAt:          time.Now(),
		IsValid:           true,
		Metadata:          map[string]any{"generated_by": generatedBy},
	}

	if err := s.trainingRepo.CreateCertificate(ctx, *certificate); err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	return certificate, nil
}

func (s *TrainingService) GetUserCertificates(ctx context.Context, userID string, filters map[string]interface{}) ([]repo.Certificate, error) {
	certificates, err := s.trainingRepo.GetUserCertificates(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificates: %w", err)
	}

	return certificates, nil
}

func (s *TrainingService) GetCertificate(ctx context.Context, id string) (*repo.Certificate, error) {
	certificate, err := s.trainingRepo.GetCertificateByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate: %w", err)
	}

	return certificate, nil
}

func (s *TrainingService) ValidateCertificate(ctx context.Context, certificateNumber string) (*repo.Certificate, error) {
	certificate, err := s.trainingRepo.GetCertificateByNumber(ctx, certificateNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate by number: %w", err)
	}
	if !certificate.IsValid {
		return nil, errors.New("certificate is no longer valid")
	}

	return certificate, nil
}

// Notifications ------------------------------------------------------------

func (s *TrainingService) CreateNotification(ctx context.Context, tenantID string, req dto.CreateNotificationRequest, createdBy string) (*repo.TrainingNotification, error) {
	assignment, err := s.trainingRepo.GetAssignmentByID(ctx, req.AssignmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment: %w", err)
	}
	if assignment.TenantID != tenantID {
		return nil, errors.New("assignment does not belong to tenant")
	}

	notification := &repo.TrainingNotification{
		ID:           generateID(),
		TenantID:     tenantID,
		AssignmentID: req.AssignmentID,
		UserID:       req.UserID,
		Type:         req.Type,
		Title:        req.Title,
		Message:      req.Message,
		SentAt:       time.Now(),
		IsRead:       false,
	}

	if notification.UserID == \"\" {
		notification.UserID = assignment.UserID
	}

	if err := s.trainingRepo.CreateNotification(ctx, *notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

func (s *TrainingService) GetUserNotifications(ctx context.Context, userID string, unreadOnly bool) ([]repo.TrainingNotification, error) {
	notifications, err := s.trainingRepo.GetUserNotifications(ctx, userID, unreadOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}

func (s *TrainingService) MarkNotificationAsRead(ctx context.Context, notificationID string, userID string) error {
	if err := s.trainingRepo.MarkNotificationAsRead(ctx, notificationID, userID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}

// Analytics ----------------------------------------------------------------

func (s *TrainingService) GetUserProgress(ctx context.Context, userID string) (*repo.TrainingAnalytics, error) {
	analytics, err := s.trainingRepo.GetUserAnalytics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user analytics: %w", err)
	}

	return analytics, nil
}

func (s *TrainingService) GetCourseProgress(ctx context.Context, courseID string) (*repo.CourseAnalytics, error) {
	analytics, err := s.trainingRepo.GetCourseAnalytics(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course analytics: %w", err)
	}

	return analytics, nil
}

func (s *TrainingService) GetOrganizationAnalytics(ctx context.Context, tenantID string) (*repo.OrganizationAnalytics, error) {
	analytics, err := s.trainingRepo.GetOrganizationAnalytics(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization analytics: %w", err)
	}

	return analytics, nil
}

func (s *TrainingService) RecordAnalytics(ctx context.Context, tenantID string, req dto.RecordAnalyticsRequest) error {
	analytics := repo.TrainingAnalytics{
		ID:         generateID(),
		TenantID:   tenantID,
		UserID:     req.UserID,
		MaterialID: req.MaterialID,
		CourseID:   req.CourseID,
		MetricType: req.MetricType,
		MetricValue: req.MetricValue,
		RecordedAt: time.Now(),
	}

	if err := s.trainingRepo.CreateAnalytics(ctx, analytics); err != nil {
		return fmt.Errorf("failed to record analytics: %w", err)
	}

	return nil
}

// Bulk operations ----------------------------------------------------------

func (s *TrainingService) BulkAssignMaterial(ctx context.Context, tenantID string, req dto.BulkAssignMaterialRequest, assignedBy string) error {
	assignReq := dto.AssignMaterialRequest{
		MaterialID: req.MaterialID,
		UserIDs:    req.UserIDs,
		DueAt:      req.DueAt,
		Priority:   req.Priority,
	}

	_, err := s.AssignMaterial(ctx, tenantID, assignReq, assignedBy)
	return err
}

func (s *TrainingService) BulkAssignCourse(ctx context.Context, tenantID string, req dto.BulkAssignCourseRequest, assignedBy string) error {
	assignReq := dto.AssignCourseRequest{
		CourseID: req.CourseID,
		UserIDs:  req.UserIDs,
		DueAt:    req.DueAt,
		Priority: req.Priority,
	}

	_, err := s.AssignCourse(ctx, tenantID, assignReq, assignedBy)
	return err
}

func (s *TrainingService) BulkUpdateProgress(ctx context.Context, req dto.BulkUpdateProgressRequest, updatedBy string) error {
	updateReq := dto.UpdateProgressRequest{
		ProgressPercentage: req.ProgressPercentage,
		TimeSpentMinutes:   req.TimeSpentMinutes,
	}
	return s.UpdateProgress(ctx, req.AssignmentID, req.MaterialID, updateReq, updatedBy)
}

// Reporting helpers --------------------------------------------------------

func (s *TrainingService) GetOverdueAssignments(ctx context.Context, tenantID string) ([]repo.TrainingAssignment, error) {
	assignments, err := s.trainingRepo.GetOverdueAssignments(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue assignments: %w", err)
	}

	return assignments, nil
}

func (s *TrainingService) SendReminderNotifications(ctx context.Context, tenantID string) error {
	assignments, err := s.trainingRepo.GetUpcomingDeadlines(ctx, tenantID, reminderWindowDays)
	if err != nil {
		return fmt.Errorf("failed to get upcoming deadlines: %w", err)
	}

	overdue, err := s.trainingRepo.GetOverdueAssignments(ctx, tenantID)
	if err == nil {
		assignments = append(assignments, overdue...)
	}

	now := time.Now()
	seen := make(map[string]struct{})
	for _, assignment := range assignments {
		if _, ok := seen[assignment.ID]; ok {
			continue
		}
		seen[assignment.ID] = struct{}{}

		if assignment.ReminderSentAt != nil && now.Sub(*assignment.ReminderSentAt) < 12*time.Hour {
			continue
		}

		dueText := \"soon\"
		notificationType := \"reminder\"
		if assignment.DueAt != nil {
			dueText = assignment.DueAt.Format(time.RFC3339)
			if assignment.DueAt.Before(now) {
				notificationType = \"deadline\"
			}
		}

		message := fmt.Sprintf("Training assignment %s is due %s.", assignment.ID, dueText)
		notification := repo.TrainingNotification{
			ID:           generateID(),
			TenantID:     tenantID,
			AssignmentID: assignment.ID,
			UserID:       assignment.UserID,
			Type:         notificationType,
			Title:        \"Training deadline reminder\",
			Message:      message,
			SentAt:       now,
		}

		if err := s.trainingRepo.CreateNotification(ctx, notification); err != nil {
			return fmt.Errorf("failed to create reminder: %w", err)
		}

		assignment.ReminderSentAt = timePtr(now)
		assignment.Priority = ensurePriority(assignment.Priority)
		if err := s.trainingRepo.UpdateAssignment(ctx, assignment); err != nil {
			return fmt.Errorf("failed to update assignment after reminder: %w", err)
		}
	}

	return nil
}

func (s *TrainingService) GetUpcomingDeadlines(ctx context.Context, tenantID string, days int) ([]repo.TrainingAssignment, error) {
	assignments, err := s.trainingRepo.GetUpcomingDeadlines(ctx, tenantID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming deadlines: %w", err)
	}

	return assignments, nil
}

// Helpers ------------------------------------------------------------------

func (s *TrainingService) updateAssignmentState(ctx context.Context, assignment *repo.TrainingAssignment, progress int, timeSpent int, completedAt *time.Time) error {
	progress = clamp(progress, 0, 100)
	assignment.ProgressPercentage = progress
	if timeSpent >= 0 {
		assignment.TimeSpentMinutes = timeSpent
	}

	now := time.Now()
	if completedAt != nil {
		assignment.CompletedAt = completedAt
		assignment.Status = assignmentStatusCompleted
	} else if progress >= 100 {
		t := now
		assignment.CompletedAt = &t
		assignment.Status = assignmentStatusCompleted
	} else if progress > 0 && assignment.Status == assignmentStatusAssigned {
		assignment.Status = assignmentStatusInProgress
		assignment.CompletedAt = nil
	}

	assignment.Priority = ensurePriority(assignment.Priority)

	if assignment.Status != assignmentStatusCompleted && assignment.DueAt != nil && assignment.DueAt.Before(now) {
		assignment.Status = assignmentStatusOverdue
	}

	return s.trainingRepo.UpdateAssignment(ctx, *assignment)
}

func copyMetadata(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func ensurePriority(value string) string {
	if strings.TrimSpace(value) == "" {
		return defaultPriorityValue
	}
	return value
}

func percentFromScore(score, max int) int {
	if max <= 0 {
		return 0
	}
	return int(math.Round((float64(score) / float64(max)) * 100))
}

func calculateQuizScore(questions []repo.QuizQuestion, answers map[string]any) (int, int) {
	score := 0
	maxScore := 0
	for _, question := range questions {
		points := question.Points
		if points <= 0 {
			points = 1
		}
		maxScore += points

		selected, ok := extractSelectedIndex(answers[question.ID])
		if ok && selected == question.CorrectIndex {
			score += points
		}
	}

	return score, maxScore
}

func extractSelectedIndex(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func stringPtr(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	v := value
	return &v
}

func timePtr(t time.Time) *time.Time {
	v := t
	return &v
}

func intPtr(v int) *int {
	value := v
	return &value
}

func generateID() string {
	return uuid.NewString()
}
"""

Path(r"d:/CompliSec/apps/backend/internal/domain/training_service_impl.go").write_text(go_code, encoding="utf-8")
