package repo

import (
	"context"
	"database/sql"
)

// DBInterface - интерфейс для базы данных
type DBInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// RiskRepoInterface - интерфейс для RiskRepo
type RiskRepoInterface interface {
	GetByIDWithTenant(ctx context.Context, id, tenantID string) (*Risk, error)
}

// TrainingRepoInterface - интерфейс для TrainingRepo
type TrainingRepoInterface interface {
	// Materials
	CreateMaterial(ctx context.Context, material Material) error
	GetMaterialByID(ctx context.Context, id string) (*Material, error)
	ListMaterials(ctx context.Context, tenantID string, filters map[string]interface{}) ([]Material, error)
	UpdateMaterial(ctx context.Context, material Material) error
	DeleteMaterial(ctx context.Context, id string) error

	// Courses
	CreateCourse(ctx context.Context, course TrainingCourse) error
	GetCourseByID(ctx context.Context, id string) (*TrainingCourse, error)
	ListCourses(ctx context.Context, tenantID string, filters map[string]interface{}) ([]TrainingCourse, error)
	UpdateCourse(ctx context.Context, course TrainingCourse) error
	DeleteCourse(ctx context.Context, id string) error

	// Course Materials
	AddMaterialToCourse(ctx context.Context, courseMaterial CourseMaterial) error
	RemoveMaterialFromCourse(ctx context.Context, courseID, materialID string) error
	GetCourseMaterials(ctx context.Context, courseID string) ([]CourseMaterial, error)

	// Assignments
	CreateAssignment(ctx context.Context, assignment TrainingAssignment) error
	GetAssignmentByID(ctx context.Context, id string) (*TrainingAssignment, error)
	GetUserAssignments(ctx context.Context, userID string, filters map[string]interface{}) ([]TrainingAssignment, error)
	UpdateAssignment(ctx context.Context, assignment TrainingAssignment) error
	DeleteAssignment(ctx context.Context, id string) error
	GetOverdueAssignments(ctx context.Context, tenantID string) ([]TrainingAssignment, error)
	GetUpcomingDeadlines(ctx context.Context, tenantID string, days int) ([]TrainingAssignment, error)

	// Progress
	CreateProgress(ctx context.Context, progress TrainingProgress) error
	UpdateProgress(ctx context.Context, progress TrainingProgress) error
	GetProgressByAssignment(ctx context.Context, assignmentID string) ([]TrainingProgress, error)
	GetProgressByAssignmentAndMaterial(ctx context.Context, assignmentID, materialID string) (*TrainingProgress, error)

	// Quiz Questions
	CreateQuizQuestion(ctx context.Context, question QuizQuestion) error
	GetQuizQuestionByID(ctx context.Context, id string) (*QuizQuestion, error)
	ListQuizQuestions(ctx context.Context, materialID string) ([]QuizQuestion, error)
	UpdateQuizQuestion(ctx context.Context, question QuizQuestion) error
	DeleteQuizQuestion(ctx context.Context, id string) error

	// Quiz Attempts
	CreateQuizAttempt(ctx context.Context, attempt QuizAttempt) error
	GetQuizAttemptByID(ctx context.Context, id string) (*QuizAttempt, error)
	GetQuizAttempts(ctx context.Context, assignmentID, materialID string) ([]QuizAttempt, error)

	// Certificates
	CreateCertificate(ctx context.Context, certificate Certificate) error
	GetCertificateByID(ctx context.Context, id string) (*Certificate, error)
	GetCertificateByNumber(ctx context.Context, certificateNumber string) (*Certificate, error)
	GetUserCertificates(ctx context.Context, userID string, filters map[string]interface{}) ([]Certificate, error)

	// Notifications
	CreateNotification(ctx context.Context, notification TrainingNotification) error
	GetUserNotifications(ctx context.Context, userID string, unreadOnly bool) ([]TrainingNotification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID, userID string) error

	// Analytics
	CreateAnalytics(ctx context.Context, analytics TrainingAnalytics) error
	GetUserAnalytics(ctx context.Context, userID string) (*TrainingAnalytics, error)
	GetCourseAnalytics(ctx context.Context, courseID string) (*CourseAnalytics, error)
	GetOrganizationAnalytics(ctx context.Context, tenantID string) (*OrganizationAnalytics, error)

	// Role assignments
	CreateRoleAssignment(ctx context.Context, assignment RoleTrainingAssignment) error
	GetRoleAssignments(ctx context.Context, roleID string) ([]RoleTrainingAssignment, error)
	DeleteRoleAssignment(ctx context.Context, id string) error
}

// DocumentRepoInterface - интерфейс для DocumentRepo
type DocumentRepoInterface interface {
	// Folders
	CreateFolder(ctx context.Context, folder Folder) error
	GetFolderByID(ctx context.Context, id, tenantID string) (*Folder, error)
	ListFolders(ctx context.Context, tenantID string, parentID *string) ([]Folder, error)
	UpdateFolder(ctx context.Context, folder Folder) error
	DeleteFolder(ctx context.Context, id, tenantID string) error

	// Documents
	CreateDocument(ctx context.Context, document Document) error
	GetDocumentByID(ctx context.Context, id, tenantID string) (*Document, error)
	GetDocumentsByIDs(ctx context.Context, ids []string, tenantID string) ([]Document, error)
	ListDocuments(ctx context.Context, tenantID string, filters map[string]interface{}) ([]Document, error)
	ListAllDocuments(ctx context.Context, tenantID string, filters map[string]interface{}) ([]Document, error)
	UpdateDocument(ctx context.Context, document Document) error
	DeleteDocument(ctx context.Context, id, tenantID string) error

	// Document Tags
	AddDocumentTag(ctx context.Context, documentID, tag string) error
	RemoveDocumentTag(ctx context.Context, documentID, tag string) error
	GetDocumentTags(ctx context.Context, documentID string) ([]string, error)
	GetDocumentsTags(ctx context.Context, documentIDs []string) (map[string][]string, error)

	// Document Links
	AddDocumentLink(ctx context.Context, link DocumentLink) error
	GetDocumentLinks(ctx context.Context, documentID string) ([]DocumentLink, error)
	GetDocumentsLinks(ctx context.Context, documentIDs []string) (map[string][]DocumentLink, error)
	DeleteDocumentLink(ctx context.Context, documentID, module, entityID string) error
	HasModuleLinks(ctx context.Context, documentID string) (bool, error) // Проверяет, есть ли у документа связи с модулями

	// OCR Text
	CreateOCRText(ctx context.Context, ocrText OCRText) error
	GetOCRText(ctx context.Context, documentID string) (*OCRText, error)
	GetDocumentsOCRTexts(ctx context.Context, documentIDs []string) (map[string]*string, error)

	// Document Permissions
	CreateDocumentPermission(ctx context.Context, permission DocumentPermission) error
	GetDocumentPermissions(ctx context.Context, objectType, objectID, tenantID string) ([]DocumentPermission, error)

	// Document Versions
	CreateDocumentVersion(ctx context.Context, version DocumentVersion) error
	GetDocumentVersions(ctx context.Context, documentID string) ([]DocumentVersion, error)
	GetDocumentVersion(ctx context.Context, versionID, tenantID string) (DocumentVersion, error)

	// Document Audit Log
	CreateDocumentAuditLog(ctx context.Context, log DocumentAuditLog) error
	GetDocumentAuditLog(ctx context.Context, tenantID string, filters map[string]interface{}) ([]DocumentAuditLog, error)

	// Search
	SearchDocuments(ctx context.Context, tenantID, searchTerm string) ([]Document, error)
}
