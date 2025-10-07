package dto

import (
	"time"
)

// Material DTOs

type CreateMaterialRequest struct {
	Title           string         `json:"title" validate:"required,min=1,max=255"`
	Description     *string        `json:"description,omitempty"`
	URI             string         `json:"uri" validate:"required,min=1"`
	Type            string         `json:"type" validate:"required,oneof=file link video"`
	MaterialType    string         `json:"material_type" validate:"required,oneof=document video quiz simulation acknowledgment"`
	DurationMinutes *int           `json:"duration_minutes,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	IsRequired      bool           `json:"is_required"`
	PassingScore    int            `json:"passing_score" validate:"min=0,max=100"`
	AttemptsLimit   *int           `json:"attempts_limit,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

type UpdateMaterialRequest struct {
	Title           *string        `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description     *string        `json:"description,omitempty"`
	URI             *string        `json:"uri,omitempty" validate:"omitempty,min=1"`
	Type            *string        `json:"type,omitempty" validate:"omitempty,oneof=file link video"`
	MaterialType    *string        `json:"material_type,omitempty" validate:"omitempty,oneof=document video quiz simulation acknowledgment"`
	DurationMinutes *int           `json:"duration_minutes,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	IsRequired      *bool          `json:"is_required,omitempty"`
	PassingScore    *int           `json:"passing_score,omitempty" validate:"omitempty,min=0,max=100"`
	AttemptsLimit   *int           `json:"attempts_limit,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

type MaterialResponse struct {
	ID              string         `json:"id"`
	TenantID        string         `json:"tenant_id"`
	Title           string         `json:"title"`
	Description     *string        `json:"description"`
	URI             string         `json:"uri"`
	Type            string         `json:"type"`
	MaterialType    string         `json:"material_type"`
	DurationMinutes *int           `json:"duration_minutes"`
	Tags            []string       `json:"tags"`
	IsRequired      bool           `json:"is_required"`
	PassingScore    int            `json:"passing_score"`
	AttemptsLimit   *int           `json:"attempts_limit"`
	Metadata        map[string]any `json:"metadata"`
	CreatedBy       *string        `json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// Course DTOs

type CreateCourseRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
}

type UpdateCourseRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type CourseResponse struct {
	ID          string             `json:"id"`
	TenantID    string             `json:"tenant_id"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	IsActive    bool               `json:"is_active"`
	CreatedBy   *string            `json:"created_by"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Materials   []MaterialResponse `json:"materials,omitempty"`
}

type CourseMaterialRequest struct {
	OrderIndex int  `json:"order_index"`
	IsRequired bool `json:"is_required"`
}

type CourseMaterialResponse struct {
	ID         string           `json:"id"`
	CourseID   string           `json:"course_id"`
	MaterialID string           `json:"material_id"`
	OrderIndex int              `json:"order_index"`
	IsRequired bool             `json:"is_required"`
	CreatedAt  time.Time        `json:"created_at"`
	Material   MaterialResponse `json:"material,omitempty"`
}

// Assignment DTOs

type AssignMaterialRequest struct {
	MaterialID string         `json:"material_id" validate:"required"`
	UserIDs    []string       `json:"user_ids" validate:"required,min=1"`
	DueAt      *time.Time     `json:"due_at,omitempty"`
	Priority   string         `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type AssignCourseRequest struct {
	CourseID string         `json:"course_id" validate:"required"`
	UserIDs  []string       `json:"user_ids" validate:"required,min=1"`
	DueAt    *time.Time     `json:"due_at,omitempty"`
	Priority string         `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type AssignToRoleRequest struct {
	RoleID     string  `json:"role_id" validate:"required"`
	MaterialID *string `json:"material_id,omitempty"`
	CourseID   *string `json:"course_id,omitempty"`
	IsRequired bool    `json:"is_required"`
	DueDays    *int    `json:"due_days,omitempty"`
}

type UpdateAssignmentRequest struct {
	Status   *string        `json:"status,omitempty" validate:"omitempty,oneof=assigned in_progress completed overdue"`
	DueAt    *time.Time     `json:"due_at,omitempty"`
	Priority *string        `json:"priority,omitempty" validate:"omitempty,oneof=low normal high urgent"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type TrainingAssignmentResponse struct {
	ID                 string                     `json:"id"`
	TenantID           string                     `json:"tenant_id"`
	MaterialID         *string                    `json:"material_id"`
	CourseID           *string                    `json:"course_id"`
	UserID             string                     `json:"user_id"`
	Status             string                     `json:"status"`
	DueAt              *time.Time                 `json:"due_at"`
	CompletedAt        *time.Time                 `json:"completed_at"`
	AssignedBy         *string                    `json:"assigned_by"`
	Priority           string                     `json:"priority"`
	ProgressPercentage int                        `json:"progress_percentage"`
	TimeSpentMinutes   int                        `json:"time_spent_minutes"`
	LastAccessedAt     *time.Time                 `json:"last_accessed_at"`
	ReminderSentAt     *time.Time                 `json:"reminder_sent_at"`
	Metadata           map[string]any             `json:"metadata"`
	CreatedAt          time.Time                  `json:"created_at"`
	Material           *MaterialResponse          `json:"material,omitempty"`
	Course             *CourseResponse            `json:"course,omitempty"`
	User               *UserResponse              `json:"user,omitempty"`
	Progress           []TrainingProgressResponse `json:"progress,omitempty"`
}

// Progress DTOs

type UpdateProgressRequest struct {
	ProgressPercentage int        `json:"progress_percentage" validate:"min=0,max=100"`
	TimeSpentMinutes   int        `json:"time_spent_minutes" validate:"min=0"`
	LastPosition       *int       `json:"last_position,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
}

type TrainingProgressResponse struct {
	ID                 string           `json:"id"`
	AssignmentID       string           `json:"assignment_id"`
	MaterialID         string           `json:"material_id"`
	ProgressPercentage int              `json:"progress_percentage"`
	TimeSpentMinutes   int              `json:"time_spent_minutes"`
	LastPosition       *int             `json:"last_position"`
	CompletedAt        *time.Time       `json:"completed_at"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	Material           MaterialResponse `json:"material,omitempty"`
}

// Quiz DTOs

type CreateQuizQuestionRequest struct {
	Text         string         `json:"text" validate:"required,min=1"`
	OptionsJSON  map[string]any `json:"options_json" validate:"required"`
	CorrectIndex int            `json:"correct_index" validate:"min=0"`
	QuestionType string         `json:"question_type" validate:"required,oneof=multiple_choice single_choice true_false text_input"`
	Points       int            `json:"points" validate:"min=1"`
	Explanation  *string        `json:"explanation,omitempty"`
	OrderIndex   int            `json:"order_index"`
}

type UpdateQuizQuestionRequest struct {
	Text         *string        `json:"text,omitempty" validate:"omitempty,min=1"`
	OptionsJSON  map[string]any `json:"options_json,omitempty"`
	CorrectIndex *int           `json:"correct_index,omitempty" validate:"omitempty,min=0"`
	QuestionType *string        `json:"question_type,omitempty" validate:"omitempty,oneof=multiple_choice single_choice true_false text_input"`
	Points       *int           `json:"points,omitempty" validate:"omitempty,min=1"`
	Explanation  *string        `json:"explanation,omitempty"`
	OrderIndex   *int           `json:"order_index,omitempty"`
}

type QuizQuestionResponse struct {
	ID           string         `json:"id"`
	MaterialID   string         `json:"material_id"`
	Text         string         `json:"text"`
	OptionsJSON  map[string]any `json:"options_json"`
	CorrectIndex int            `json:"correct_index"`
	QuestionType string         `json:"question_type"`
	Points       int            `json:"points"`
	Explanation  *string        `json:"explanation"`
	OrderIndex   int            `json:"order_index"`
	CreatedAt    time.Time      `json:"created_at"`
}

type SubmitQuizAttemptRequest struct {
	AnswersJSON      map[string]any `json:"answers_json" validate:"required"`
	TimeSpentMinutes *int           `json:"time_spent_minutes,omitempty"`
}

type QuizAttemptResponse struct {
	ID               string           `json:"id"`
	UserID           string           `json:"user_id"`
	MaterialID       string           `json:"material_id"`
	AssignmentID     *string          `json:"assignment_id"`
	Score            int              `json:"score"`
	MaxScore         *int             `json:"max_score"`
	Passed           bool             `json:"passed"`
	AnswersJSON      map[string]any   `json:"answers_json"`
	TimeSpentMinutes *int             `json:"time_spent_minutes"`
	AttemptedAt      time.Time        `json:"attempted_at"`
	Material         MaterialResponse `json:"material,omitempty"`
}

// Certificate DTOs

type CertificateResponse struct {
	ID                string            `json:"id"`
	TenantID          string            `json:"tenant_id"`
	AssignmentID      string            `json:"assignment_id"`
	UserID            string            `json:"user_id"`
	MaterialID        *string           `json:"material_id"`
	CourseID          *string           `json:"course_id"`
	CertificateNumber string            `json:"certificate_number"`
	IssuedAt          time.Time         `json:"issued_at"`
	ExpiresAt         *time.Time        `json:"expires_at"`
	IsValid           bool              `json:"is_valid"`
	Metadata          map[string]any    `json:"metadata"`
	CreatedAt         time.Time         `json:"created_at"`
	User              *UserResponse     `json:"user,omitempty"`
	Material          *MaterialResponse `json:"material,omitempty"`
	Course            *CourseResponse   `json:"course,omitempty"`
}

// Notification DTOs

type CreateNotificationRequest struct {
	AssignmentID string `json:"assignment_id" validate:"required"`
	UserID       string `json:"user_id" validate:"required"`
	Type         string `json:"type" validate:"required,oneof=assignment reminder deadline completion"`
	Title        string `json:"title" validate:"required,min=1,max=255"`
	Message      string `json:"message" validate:"required,min=1"`
}

type TrainingNotificationResponse struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	AssignmentID string     `json:"assignment_id"`
	UserID       string     `json:"user_id"`
	Type         string     `json:"type"`
	Title        string     `json:"title"`
	Message      string     `json:"message"`
	IsRead       bool       `json:"is_read"`
	SentAt       time.Time  `json:"sent_at"`
	ReadAt       *time.Time `json:"read_at"`
}

// Analytics DTOs

type RecordAnalyticsRequest struct {
	UserID      *string `json:"user_id,omitempty"`
	MaterialID  *string `json:"material_id,omitempty"`
	CourseID    *string `json:"course_id,omitempty"`
	MetricType  string  `json:"metric_type" validate:"required"`
	MetricValue float64 `json:"metric_value" validate:"required"`
}

type TrainingAnalyticsResponse struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	UserID      *string   `json:"user_id"`
	MaterialID  *string   `json:"material_id"`
	CourseID    *string   `json:"course_id"`
	MetricType  string    `json:"metric_type"`
	MetricValue float64   `json:"metric_value"`
	RecordedAt  time.Time `json:"recorded_at"`
}

type CourseAnalyticsResponse struct {
	CourseID             string  `json:"course_id"`
	TotalAssignments     int     `json:"total_assignments"`
	CompletedAssignments int     `json:"completed_assignments"`
	CompletionRate       float64 `json:"completion_rate"`
	AverageTimeSpent     int     `json:"average_time_spent"`
	AverageScore         float64 `json:"average_score"`
}

type OrganizationAnalyticsResponse struct {
	TenantID             string  `json:"tenant_id"`
	TotalMaterials       int     `json:"total_materials"`
	TotalCourses         int     `json:"total_courses"`
	TotalAssignments     int     `json:"total_assignments"`
	CompletedAssignments int     `json:"completed_assignments"`
	OverdueAssignments   int     `json:"overdue_assignments"`
	CompletionRate       float64 `json:"completion_rate"`
	AverageTimeSpent     int     `json:"average_time_spent"`
}

// Bulk operation DTOs

type BulkAssignMaterialRequest struct {
	MaterialID string     `json:"material_id" validate:"required"`
	UserIDs    []string   `json:"user_ids" validate:"required,min=1"`
	DueAt      *time.Time `json:"due_at,omitempty"`
	Priority   string     `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
}

type BulkAssignCourseRequest struct {
	CourseID string     `json:"course_id" validate:"required"`
	UserIDs  []string   `json:"user_ids" validate:"required,min=1"`
	DueAt    *time.Time `json:"due_at,omitempty"`
	Priority string     `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
}

type BulkUpdateProgressRequest struct {
	AssignmentID       string `json:"assignment_id" validate:"required"`
	MaterialID         string `json:"material_id" validate:"required"`
	ProgressPercentage int    `json:"progress_percentage" validate:"min=0,max=100"`
	TimeSpentMinutes   int    `json:"time_spent_minutes" validate:"min=0"`
}

// List and filter DTOs

type TrainingListRequest struct {
	Page      int                    `json:"page" validate:"min=1"`
	PageSize  int                    `json:"page_size" validate:"min=1,max=100"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	SortBy    string                 `json:"sort_by,omitempty"`
	SortOrder string                 `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
}

type TrainingListResponse struct {
	Items      []interface{} `json:"items"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

