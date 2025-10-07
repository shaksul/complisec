package repo

import (
	"context"
	"database/sql"
	"time"
)

type DB struct {
	*sql.DB
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.Query(query, args...)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.DB.ExecContext(ctx, query, args...)
}

// Training-related types

// Material represents a training material
type Material struct {
	ID              string         `json:"id" db:"id"`
	TenantID        string         `json:"tenant_id" db:"tenant_id"`
	Title           string         `json:"title" db:"title"`
	Description     *string        `json:"description" db:"description"`
	URI             string         `json:"uri" db:"uri"`
	Type            string         `json:"type" db:"type"`
	MaterialType    string         `json:"material_type" db:"material_type"`
	DurationMinutes *int           `json:"duration_minutes" db:"duration_minutes"`
	Tags            []string       `json:"tags" db:"tags"`
	IsRequired      bool           `json:"is_required" db:"is_required"`
	PassingScore    int            `json:"passing_score" db:"passing_score"`
	AttemptsLimit   *int           `json:"attempts_limit" db:"attempts_limit"`
	Metadata        map[string]any `json:"metadata" db:"metadata"`
	CreatedBy       *string        `json:"created_by" db:"created_by"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// TrainingCourse represents a training course
type TrainingCourse struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedBy   *string   `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CourseMaterial represents a junction between course and material
type CourseMaterial struct {
	ID         string    `json:"id" db:"id"`
	CourseID   string    `json:"course_id" db:"course_id"`
	MaterialID string    `json:"material_id" db:"material_id"`
	OrderIndex int       `json:"order_index" db:"order_index"`
	IsRequired bool      `json:"is_required" db:"is_required"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// TrainingAssignment represents a training assignment
type TrainingAssignment struct {
	ID                 string         `json:"id" db:"id"`
	TenantID           string         `json:"tenant_id" db:"tenant_id"`
	MaterialID         *string        `json:"material_id" db:"material_id"`
	CourseID           *string        `json:"course_id" db:"course_id"`
	UserID             string         `json:"user_id" db:"user_id"`
	Status             string         `json:"status" db:"status"`
	DueAt              *time.Time     `json:"due_at" db:"due_at"`
	CompletedAt        *time.Time     `json:"completed_at" db:"completed_at"`
	AssignedBy         *string        `json:"assigned_by" db:"assigned_by"`
	Priority           string         `json:"priority" db:"priority"`
	ProgressPercentage int            `json:"progress_percentage" db:"progress_percentage"`
	TimeSpentMinutes   int            `json:"time_spent_minutes" db:"time_spent_minutes"`
	LastAccessedAt     *time.Time     `json:"last_accessed_at" db:"last_accessed_at"`
	ReminderSentAt     *time.Time     `json:"reminder_sent_at" db:"reminder_sent_at"`
	Metadata           map[string]any `json:"metadata" db:"metadata"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
}

// TrainingProgress represents progress tracking for a material
type TrainingProgress struct {
	ID                 string     `json:"id" db:"id"`
	AssignmentID       string     `json:"assignment_id" db:"assignment_id"`
	MaterialID         string     `json:"material_id" db:"material_id"`
	ProgressPercentage int        `json:"progress_percentage" db:"progress_percentage"`
	TimeSpentMinutes   int        `json:"time_spent_minutes" db:"time_spent_minutes"`
	LastPosition       *int       `json:"last_position" db:"last_position"`
	CompletedAt        *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// QuizQuestion represents a quiz question
type QuizQuestion struct {
	ID           string         `json:"id" db:"id"`
	MaterialID   string         `json:"material_id" db:"material_id"`
	Text         string         `json:"text" db:"text"`
	OptionsJSON  map[string]any `json:"options_json" db:"options_json"`
	CorrectIndex int            `json:"correct_index" db:"correct_index"`
	QuestionType string         `json:"question_type" db:"question_type"`
	Points       int            `json:"points" db:"points"`
	Explanation  *string        `json:"explanation" db:"explanation"`
	OrderIndex   int            `json:"order_index" db:"order_index"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
}

// QuizAttempt represents a quiz attempt
type QuizAttempt struct {
	ID               string         `json:"id" db:"id"`
	UserID           string         `json:"user_id" db:"user_id"`
	MaterialID       string         `json:"material_id" db:"material_id"`
	AssignmentID     *string        `json:"assignment_id" db:"assignment_id"`
	Score            int            `json:"score" db:"score"`
	MaxScore         *int           `json:"max_score" db:"max_score"`
	Passed           bool           `json:"passed" db:"passed"`
	AnswersJSON      map[string]any `json:"answers_json" db:"answers_json"`
	TimeSpentMinutes *int           `json:"time_spent_minutes" db:"time_spent_minutes"`
	AttemptedAt      time.Time      `json:"attempted_at" db:"attempted_at"`
}

// Certificate represents a training certificate
type Certificate struct {
	ID                string         `json:"id" db:"id"`
	TenantID          string         `json:"tenant_id" db:"tenant_id"`
	AssignmentID      string         `json:"assignment_id" db:"assignment_id"`
	UserID            string         `json:"user_id" db:"user_id"`
	MaterialID        *string        `json:"material_id" db:"material_id"`
	CourseID          *string        `json:"course_id" db:"course_id"`
	CertificateNumber string         `json:"certificate_number" db:"certificate_number"`
	IssuedAt          time.Time      `json:"issued_at" db:"issued_at"`
	ExpiresAt         *time.Time     `json:"expires_at" db:"expires_at"`
	IsValid           bool           `json:"is_valid" db:"is_valid"`
	Metadata          map[string]any `json:"metadata" db:"metadata"`
	CreatedAt         time.Time      `json:"created_at" db:"created_at"`
}

// TrainingNotification represents a training notification
type TrainingNotification struct {
	ID           string     `json:"id" db:"id"`
	TenantID     string     `json:"tenant_id" db:"tenant_id"`
	AssignmentID string     `json:"assignment_id" db:"assignment_id"`
	UserID       string     `json:"user_id" db:"user_id"`
	Type         string     `json:"type" db:"type"`
	Title        string     `json:"title" db:"title"`
	Message      string     `json:"message" db:"message"`
	IsRead       bool       `json:"is_read" db:"is_read"`
	SentAt       time.Time  `json:"sent_at" db:"sent_at"`
	ReadAt       *time.Time `json:"read_at" db:"read_at"`
}

// TrainingAnalytics represents training analytics
type TrainingAnalytics struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	UserID      *string   `json:"user_id" db:"user_id"`
	MaterialID  *string   `json:"material_id" db:"material_id"`
	CourseID    *string   `json:"course_id" db:"course_id"`
	MetricType  string    `json:"metric_type" db:"metric_type"`
	MetricValue float64   `json:"metric_value" db:"metric_value"`
	RecordedAt  time.Time `json:"recorded_at" db:"recorded_at"`
}

// RoleTrainingAssignment represents role-based training assignment
type RoleTrainingAssignment struct {
	ID         string    `json:"id" db:"id"`
	TenantID   string    `json:"tenant_id" db:"tenant_id"`
	RoleID     string    `json:"role_id" db:"role_id"`
	MaterialID *string   `json:"material_id" db:"material_id"`
	CourseID   *string   `json:"course_id" db:"course_id"`
	IsRequired bool      `json:"is_required" db:"is_required"`
	DueDays    *int      `json:"due_days" db:"due_days"`
	AssignedBy *string   `json:"assigned_by" db:"assigned_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Analytics summary types
type CourseAnalytics struct {
	CourseID             string  `json:"course_id"`
	TotalAssignments     int     `json:"total_assignments"`
	CompletedAssignments int     `json:"completed_assignments"`
	CompletionRate       float64 `json:"completion_rate"`
	AverageTimeSpent     int     `json:"average_time_spent"`
	AverageScore         float64 `json:"average_score"`
}

type OrganizationAnalytics struct {
	TenantID             string  `json:"tenant_id"`
	TotalMaterials       int     `json:"total_materials"`
	TotalCourses         int     `json:"total_courses"`
	TotalAssignments     int     `json:"total_assignments"`
	CompletedAssignments int     `json:"completed_assignments"`
	OverdueAssignments   int     `json:"overdue_assignments"`
	CompletionRate       float64 `json:"completion_rate"`
	AverageTimeSpent     int     `json:"average_time_spent"`
}
