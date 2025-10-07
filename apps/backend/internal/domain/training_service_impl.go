package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

// TrainingService - реализация интерфейса TrainingServiceInterface
type TrainingService struct {
	trainingRepo           repo.TrainingRepoInterface
	documentStorageService DocumentStorageServiceInterface
}

// NewTrainingService создает новый экземпляр TrainingService
func NewTrainingService(trainingRepo repo.TrainingRepoInterface, documentStorageService DocumentStorageServiceInterface) *TrainingService {
	return &TrainingService{
		trainingRepo:           trainingRepo,
		documentStorageService: documentStorageService,
	}
}

// Materials management
func (s *TrainingService) CreateMaterial(ctx context.Context, tenantID string, req dto.CreateMaterialRequest, createdBy string) (*repo.Material, error) {
	log.Printf("DEBUG: training_service.CreateMaterial tenant=%s title=%s", tenantID, req.Title)

	material := repo.Material{
		ID:              uuid.New().String(),
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
		Metadata:        req.Metadata,
		CreatedBy:       trainingStringPtr(createdBy),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := s.trainingRepo.CreateMaterial(ctx, material)
	if err != nil {
		log.Printf("ERROR: training_service.CreateMaterial CreateMaterial: %v", err)
		return nil, err
	}

	log.Printf("DEBUG: training_service.CreateMaterial success id=%s", material.ID)
	return &material, nil
}

func (s *TrainingService) GetMaterial(ctx context.Context, id string) (*repo.Material, error) {
	return s.trainingRepo.GetMaterialByID(ctx, id)
}

func (s *TrainingService) ListMaterials(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.Material, error) {
	return s.trainingRepo.ListMaterials(ctx, tenantID, filters)
}

func (s *TrainingService) UpdateMaterial(ctx context.Context, id string, req dto.UpdateMaterialRequest, updatedBy string) error {
	log.Printf("DEBUG: training_service.UpdateMaterial id=%s", id)

	material, err := s.trainingRepo.GetMaterialByID(ctx, id)
	if err != nil {
		return err
	}
	if material == nil {
		return errors.New("material not found")
	}

	// Обновляем поля
	if req.Title != nil {
		material.Title = *req.Title
	}
	if req.Description != nil {
		material.Description = req.Description
	}
	if req.Type != nil {
		material.Type = *req.Type
	}
	if req.URI != nil {
		material.URI = *req.URI
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
		material.Metadata = req.Metadata
	}
	material.UpdatedAt = time.Now()

	err = s.trainingRepo.UpdateMaterial(ctx, *material)
	if err != nil {
		log.Printf("ERROR: training_service.UpdateMaterial UpdateMaterial: %v", err)
		return err
	}

	log.Printf("DEBUG: training_service.UpdateMaterial success id=%s", id)
	return nil
}

func (s *TrainingService) DeleteMaterial(ctx context.Context, id string, deletedBy string) error {
	log.Printf("DEBUG: training_service.DeleteMaterial id=%s", id)

	err := s.trainingRepo.DeleteMaterial(ctx, id)
	if err != nil {
		log.Printf("ERROR: training_service.DeleteMaterial DeleteMaterial: %v", err)
		return err
	}

	log.Printf("DEBUG: training_service.DeleteMaterial success id=%s", id)
	return nil
}

// Courses management
func (s *TrainingService) CreateCourse(ctx context.Context, tenantID string, req dto.CreateCourseRequest, createdBy string) (*repo.TrainingCourse, error) {
	log.Printf("DEBUG: training_service.CreateCourse tenant=%s title=%s", tenantID, req.Title)

	course := repo.TrainingCourse{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Title:       req.Title,
		Description: req.Description,
		IsActive:    req.IsActive,
		CreatedBy:   trainingStringPtr(createdBy),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.trainingRepo.CreateCourse(ctx, course)
	if err != nil {
		log.Printf("ERROR: training_service.CreateCourse CreateCourse: %v", err)
		return nil, err
	}

	log.Printf("DEBUG: training_service.CreateCourse success id=%s", course.ID)
	return &course, nil
}

func (s *TrainingService) GetCourse(ctx context.Context, id string) (*repo.TrainingCourse, error) {
	return s.trainingRepo.GetCourseByID(ctx, id)
}

func (s *TrainingService) ListCourses(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.TrainingCourse, error) {
	return s.trainingRepo.ListCourses(ctx, tenantID, filters)
}

func (s *TrainingService) UpdateCourse(ctx context.Context, id string, req dto.UpdateCourseRequest, updatedBy string) error {
	log.Printf("DEBUG: training_service.UpdateCourse id=%s", id)

	course, err := s.trainingRepo.GetCourseByID(ctx, id)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("course not found")
	}

	// Обновляем поля
	if req.Title != nil {
		course.Title = *req.Title
	}
	if req.Description != nil {
		course.Description = req.Description
	}
	if req.IsActive != nil {
		course.IsActive = *req.IsActive
	}
	course.UpdatedAt = time.Now()

	err = s.trainingRepo.UpdateCourse(ctx, *course)
	if err != nil {
		log.Printf("ERROR: training_service.UpdateCourse UpdateCourse: %v", err)
		return err
	}

	log.Printf("DEBUG: training_service.UpdateCourse success id=%s", id)
	return nil
}

func (s *TrainingService) DeleteCourse(ctx context.Context, id string, deletedBy string) error {
	log.Printf("DEBUG: training_service.DeleteCourse id=%s", id)

	err := s.trainingRepo.DeleteCourse(ctx, id)
	if err != nil {
		log.Printf("ERROR: training_service.DeleteCourse DeleteCourse: %v", err)
		return err
	}

	log.Printf("DEBUG: training_service.DeleteCourse success id=%s", id)
	return nil
}

func (s *TrainingService) AddMaterialToCourse(ctx context.Context, courseID, materialID string, req dto.CourseMaterialRequest, addedBy string) error {
	log.Printf("DEBUG: training_service.AddMaterialToCourse courseID=%s materialID=%s", courseID, materialID)

	courseMaterial := repo.CourseMaterial{
		ID:         uuid.New().String(),
		CourseID:   courseID,
		MaterialID: materialID,
		OrderIndex: req.OrderIndex,
		IsRequired: req.IsRequired,
		CreatedAt:  time.Now(),
	}

	err := s.trainingRepo.AddMaterialToCourse(ctx, courseMaterial)
	if err != nil {
		log.Printf("ERROR: training_service.AddMaterialToCourse AddMaterialToCourse: %v", err)
		return err
	}

	log.Printf("DEBUG: training_service.AddMaterialToCourse success")
	return nil
}

func (s *TrainingService) RemoveMaterialFromCourse(ctx context.Context, courseID, materialID string, removedBy string) error {
	log.Printf("DEBUG: training_service.RemoveMaterialFromCourse courseID=%s materialID=%s", courseID, materialID)

	err := s.trainingRepo.RemoveMaterialFromCourse(ctx, courseID, materialID)
	if err != nil {
		log.Printf("ERROR: training_service.RemoveMaterialFromCourse RemoveMaterialFromCourse: %v", err)
		return err
	}

	log.Printf("DEBUG: training_service.RemoveMaterialFromCourse success")
	return nil
}

func (s *TrainingService) GetCourseMaterials(ctx context.Context, courseID string) ([]repo.CourseMaterial, error) {
	return s.trainingRepo.GetCourseMaterials(ctx, courseID)
}

// Training Document methods - интеграция с централизованным хранилищем документов
func (s *TrainingService) UploadTrainingDocument(ctx context.Context, trainingID, tenantID, documentType string, file multipart.File, header *multipart.FileHeader, req dto.UploadDocumentDTO, uploadedBy string) (*dto.DocumentDTO, error) {
	log.Printf("DEBUG: training_service.UploadTrainingDocument trainingID=%s type=%s", trainingID, documentType)

	// Обновляем запрос с информацией о материале обучения
	req.Name = fmt.Sprintf("%s - %s", documentType, req.Name)
	req.Tags = append(req.Tags, "#обучение", fmt.Sprintf("#%s", documentType))
	if req.LinkedTo == nil {
		req.LinkedTo = &dto.DocumentLinkDTO{
			Module:   "training",
			EntityID: trainingID,
		}
	}
	if req.Description == nil {
		req.Description = trainingStringPtr(fmt.Sprintf("Training document: %s", documentType))
	}

	// Загружаем документ в централизованное хранилище
	document, err := s.documentStorageService.UploadDocument(ctx, tenantID, file, header, req, uploadedBy)
	if err != nil {
		log.Printf("ERROR: training_service.UploadTrainingDocument UploadDocument: %v", err)
		return nil, err
	}

	log.Printf("DEBUG: training_service.UploadTrainingDocument success documentID=%s", document.ID)
	return document, nil
}

func (s *TrainingService) GetTrainingDocuments(ctx context.Context, trainingID, tenantID string) ([]dto.DocumentDTO, error) {
	log.Printf("DEBUG: training_service.GetTrainingDocuments trainingID=%s", trainingID)

	// Получаем документы из централизованного хранилища
	documents, err := s.documentStorageService.GetModuleDocuments(ctx, "training", trainingID, tenantID)
	if err != nil {
		log.Printf("ERROR: training_service.GetTrainingDocuments GetModuleDocuments: %v", err)
		return nil, err
	}

	log.Printf("DEBUG: training_service.GetTrainingDocuments found %d documents", len(documents))
	return documents, nil
}

func (s *TrainingService) LinkExistingDocumentToTraining(ctx context.Context, trainingID, documentID, tenantID, linkedBy string) error {
	log.Printf("DEBUG: training_service.LinkExistingDocumentToTraining trainingID=%s documentID=%s", trainingID, documentID)

	// Проверяем, что документ существует
	_, err := s.documentStorageService.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return err
	}

	// Связываем документ с материалом обучения
	err = s.documentStorageService.LinkDocumentToModule(ctx, documentID, "training", trainingID, "attachment", fmt.Sprintf("Linked to training material"), linkedBy)
	if err != nil {
		log.Printf("ERROR: training_service.LinkExistingDocumentToTraining LinkDocumentToModule: %v", err)
		return err
	}

	log.Printf("DEBUG: training_service.LinkExistingDocumentToTraining success")
	return nil
}

func (s *TrainingService) UnlinkDocumentFromTraining(ctx context.Context, trainingID, documentID, tenantID, unlinkedBy string) error {
	log.Printf("DEBUG: training_service.UnlinkDocumentFromTraining trainingID=%s documentID=%s", trainingID, documentID)

	// Отвязываем документ от материала обучения
	err := s.documentStorageService.UnlinkDocumentFromModule(ctx, documentID, "training", trainingID, unlinkedBy)
	if err != nil {
		log.Printf("ERROR: training_service.UnlinkDocumentFromTraining UnlinkDocumentFromModule: %v", err)
		return err
	}

	log.Printf("DEBUG: training_service.UnlinkDocumentFromTraining success")
	return nil
}

func (s *TrainingService) DeleteTrainingDocument(ctx context.Context, trainingID, documentID, tenantID, deletedBy string) error {
	log.Printf("DEBUG: training_service.DeleteTrainingDocument trainingID=%s documentID=%s", trainingID, documentID)

	// Получаем информацию о документе
	_, err := s.documentStorageService.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return err
	}

	// Удаляем документ из централизованного хранилища
	err = s.documentStorageService.DeleteDocument(ctx, documentID, tenantID, deletedBy)
	if err != nil {
		log.Printf("ERROR: training_service.DeleteTrainingDocument DeleteDocument: %v", err)
		return err
	}

	log.Printf("DEBUG: training_service.DeleteTrainingDocument success")
	return nil
}

// Остальные методы интерфейса (заглушки)
func (s *TrainingService) CreateQuizQuestion(ctx context.Context, materialID string, req dto.CreateQuizQuestionRequest, createdBy string) (*repo.QuizQuestion, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetQuizQuestion(ctx context.Context, id string) (*repo.QuizQuestion, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) ListQuizQuestions(ctx context.Context, materialID string) ([]repo.QuizQuestion, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) UpdateQuizQuestion(ctx context.Context, id string, req dto.UpdateQuizQuestionRequest, updatedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) DeleteQuizQuestion(ctx context.Context, id string, deletedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) AssignMaterial(ctx context.Context, tenantID string, req dto.AssignMaterialRequest, assignedBy string) (*repo.TrainingAssignment, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) AssignCourse(ctx context.Context, tenantID string, req dto.AssignCourseRequest, assignedBy string) (*repo.TrainingAssignment, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) AssignToRole(ctx context.Context, tenantID string, req dto.AssignToRoleRequest, assignedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetUserAssignments(ctx context.Context, userID string, filters map[string]interface{}) ([]repo.TrainingAssignment, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetAssignment(ctx context.Context, id string) (*repo.TrainingAssignment, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) UpdateAssignment(ctx context.Context, id string, req dto.UpdateAssignmentRequest, updatedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) DeleteAssignment(ctx context.Context, id string, deletedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) UpdateProgress(ctx context.Context, assignmentID, materialID string, req dto.UpdateProgressRequest, updatedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetProgress(ctx context.Context, assignmentID string) ([]repo.TrainingProgress, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) MarkAsCompleted(ctx context.Context, assignmentID, materialID string, completedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) SubmitQuizAttempt(ctx context.Context, assignmentID, materialID string, req dto.SubmitQuizAttemptRequest, submittedBy string) (*repo.QuizAttempt, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetQuizAttempts(ctx context.Context, assignmentID, materialID string) ([]repo.QuizAttempt, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetQuizAttempt(ctx context.Context, id string) (*repo.QuizAttempt, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GenerateCertificate(ctx context.Context, assignmentID string, generatedBy string) (*repo.Certificate, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetUserCertificates(ctx context.Context, userID string, filters map[string]interface{}) ([]repo.Certificate, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetCertificate(ctx context.Context, id string) (*repo.Certificate, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) ValidateCertificate(ctx context.Context, certificateNumber string) (*repo.Certificate, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) CreateNotification(ctx context.Context, tenantID string, req dto.CreateNotificationRequest, createdBy string) (*repo.TrainingNotification, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetUserNotifications(ctx context.Context, userID string, unreadOnly bool) ([]repo.TrainingNotification, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) MarkNotificationAsRead(ctx context.Context, notificationID string, userID string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetUserProgress(ctx context.Context, userID string) (*repo.TrainingAnalytics, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetCourseProgress(ctx context.Context, courseID string) (*repo.CourseAnalytics, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetOrganizationAnalytics(ctx context.Context, tenantID string) (*repo.OrganizationAnalytics, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) RecordAnalytics(ctx context.Context, tenantID string, req dto.RecordAnalyticsRequest) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) BulkAssignMaterial(ctx context.Context, tenantID string, req dto.BulkAssignMaterialRequest, assignedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) BulkAssignCourse(ctx context.Context, tenantID string, req dto.BulkAssignCourseRequest, assignedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) BulkUpdateProgress(ctx context.Context, req dto.BulkUpdateProgressRequest, updatedBy string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetOverdueAssignments(ctx context.Context, tenantID string) ([]repo.TrainingAssignment, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TrainingService) SendReminderNotifications(ctx context.Context, tenantID string) error {
	return fmt.Errorf("not implemented yet")
}

func (s *TrainingService) GetUpcomingDeadlines(ctx context.Context, tenantID string, days int) ([]repo.TrainingAssignment, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// Helper function to create string pointer
func trainingStringPtr(s string) *string {
	return &s
}
