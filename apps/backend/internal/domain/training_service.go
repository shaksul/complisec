package domain

import (
	"context"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"
)

// TrainingServiceInterface - интерфейс для TrainingService
type TrainingServiceInterface interface {
	// Materials management
	CreateMaterial(ctx context.Context, tenantID string, req dto.CreateMaterialRequest, createdBy string) (*repo.Material, error)
	GetMaterial(ctx context.Context, id string) (*repo.Material, error)
	ListMaterials(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.Material, error)
	UpdateMaterial(ctx context.Context, id string, req dto.UpdateMaterialRequest, updatedBy string) error
	DeleteMaterial(ctx context.Context, id string, deletedBy string) error

	// Courses management
	CreateCourse(ctx context.Context, tenantID string, req dto.CreateCourseRequest, createdBy string) (*repo.TrainingCourse, error)
	GetCourse(ctx context.Context, id string) (*repo.TrainingCourse, error)
	ListCourses(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.TrainingCourse, error)
	UpdateCourse(ctx context.Context, id string, req dto.UpdateCourseRequest, updatedBy string) error
	DeleteCourse(ctx context.Context, id string, deletedBy string) error
	AddMaterialToCourse(ctx context.Context, courseID, materialID string, req dto.CourseMaterialRequest, addedBy string) error
	RemoveMaterialFromCourse(ctx context.Context, courseID, materialID string, removedBy string) error
	GetCourseMaterials(ctx context.Context, courseID string) ([]repo.CourseMaterial, error)

	// Quiz management
	CreateQuizQuestion(ctx context.Context, materialID string, req dto.CreateQuizQuestionRequest, createdBy string) (*repo.QuizQuestion, error)
	GetQuizQuestion(ctx context.Context, id string) (*repo.QuizQuestion, error)
	ListQuizQuestions(ctx context.Context, materialID string) ([]repo.QuizQuestion, error)
	UpdateQuizQuestion(ctx context.Context, id string, req dto.UpdateQuizQuestionRequest, updatedBy string) error
	DeleteQuizQuestion(ctx context.Context, id string, deletedBy string) error

	// Assignments
	AssignMaterial(ctx context.Context, tenantID string, req dto.AssignMaterialRequest, assignedBy string) (*repo.TrainingAssignment, error)
	AssignCourse(ctx context.Context, tenantID string, req dto.AssignCourseRequest, assignedBy string) (*repo.TrainingAssignment, error)
	AssignToRole(ctx context.Context, tenantID string, req dto.AssignToRoleRequest, assignedBy string) error
	GetUserAssignments(ctx context.Context, userID string, filters map[string]interface{}) ([]repo.TrainingAssignment, error)
	GetAssignment(ctx context.Context, id string) (*repo.TrainingAssignment, error)
	UpdateAssignment(ctx context.Context, id string, req dto.UpdateAssignmentRequest, updatedBy string) error
	DeleteAssignment(ctx context.Context, id string, deletedBy string) error

	// Progress tracking
	UpdateProgress(ctx context.Context, assignmentID, materialID string, req dto.UpdateProgressRequest, updatedBy string) error
	GetProgress(ctx context.Context, assignmentID string) ([]repo.TrainingProgress, error)
	MarkAsCompleted(ctx context.Context, assignmentID, materialID string, completedBy string) error

	// Quiz attempts
	SubmitQuizAttempt(ctx context.Context, assignmentID, materialID string, req dto.SubmitQuizAttemptRequest, submittedBy string) (*repo.QuizAttempt, error)
	GetQuizAttempts(ctx context.Context, assignmentID, materialID string) ([]repo.QuizAttempt, error)
	GetQuizAttempt(ctx context.Context, id string) (*repo.QuizAttempt, error)

	// Certificates
	GenerateCertificate(ctx context.Context, assignmentID string, generatedBy string) (*repo.Certificate, error)
	GetUserCertificates(ctx context.Context, userID string, filters map[string]interface{}) ([]repo.Certificate, error)
	GetCertificate(ctx context.Context, id string) (*repo.Certificate, error)
	ValidateCertificate(ctx context.Context, certificateNumber string) (*repo.Certificate, error)

	// Notifications
	CreateNotification(ctx context.Context, tenantID string, req dto.CreateNotificationRequest, createdBy string) (*repo.TrainingNotification, error)
	GetUserNotifications(ctx context.Context, userID string, unreadOnly bool) ([]repo.TrainingNotification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string, userID string) error

	// Analytics
	GetUserProgress(ctx context.Context, userID string) (*repo.TrainingAnalytics, error)
	GetCourseProgress(ctx context.Context, courseID string) (*repo.CourseAnalytics, error)
	GetOrganizationAnalytics(ctx context.Context, tenantID string) (*repo.OrganizationAnalytics, error)
	RecordAnalytics(ctx context.Context, tenantID string, req dto.RecordAnalyticsRequest) error

	// Bulk operations
	BulkAssignMaterial(ctx context.Context, tenantID string, req dto.BulkAssignMaterialRequest, assignedBy string) error
	BulkAssignCourse(ctx context.Context, tenantID string, req dto.BulkAssignCourseRequest, assignedBy string) error
	BulkUpdateProgress(ctx context.Context, req dto.BulkUpdateProgressRequest, updatedBy string) error

	// Deadline management
	GetOverdueAssignments(ctx context.Context, tenantID string) ([]repo.TrainingAssignment, error)
	SendReminderNotifications(ctx context.Context, tenantID string) error
	GetUpcomingDeadlines(ctx context.Context, tenantID string, days int) ([]repo.TrainingAssignment, error)
}
