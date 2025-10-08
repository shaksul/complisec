package dto

import "time"

// CreateDocumentDTO represents the request to create a document
type CreateDocumentDTO struct {
	Title              string   `json:"title" validate:"required,min=1,max=255"`
	Code               *string  `json:"code"`
	Description        *string  `json:"description"`
	Type               string   `json:"type" validate:"required,oneof=policy standard procedure instruction act other"`
	Category           *string  `json:"category"`
	Tags               []string `json:"tags"`
	OwnerID            *string  `json:"owner_id"`
	Classification     string   `json:"classification" validate:"oneof=Public Internal Confidential"`
	EffectiveFrom      *string  `json:"effective_from"`
	ReviewPeriodMonths *int     `json:"review_period_months" validate:"omitempty,min=1,max=120"`
	AssetIDs           []string `json:"asset_ids"`
	RiskIDs            []string `json:"risk_ids"`
	ControlIDs         []string `json:"control_ids"`
}

// UpdateDocumentDTO represents the request to update a document
type UpdateDocumentDTO struct {
	Title       string   `json:"title" validate:"required,min=1,max=255"`
	Description *string  `json:"description"`
	Type        string   `json:"type" validate:"required,oneof=policy standard procedure instruction act other"`
	Category    *string  `json:"category"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status" validate:"oneof=draft in_review approved obsolete"`
	FolderID    *string  `json:"folder_id"`
	Metadata    *string  `json:"metadata"`
}

// CreateDocumentVersionDTO represents the request to create a document version
type CreateDocumentVersionDTO struct {
	Title          string  `json:"title" validate:"required,min=1,max=255"`
	Content        *string `json:"content"`
	FilePath       *string `json:"file_path"`
	FileSize       *int64  `json:"file_size"`
	MimeType       *string `json:"mime_type"`
	ChecksumSHA256 *string `json:"checksum_sha256"`
	EnableOCR      bool    `json:"enable_ocr"`
}

// CreateDocumentAcknowledgmentDTO represents the request to create an acknowledgment
type CreateDocumentAcknowledgmentDTO struct {
	UserID    string     `json:"user_id" validate:"required"`
	VersionID *string    `json:"version_id"`
	Deadline  *time.Time `json:"deadline"`
}

// UpdateDocumentAcknowledgmentDTO represents the request to update an acknowledgment
type UpdateDocumentAcknowledgmentDTO struct {
	Status         string     `json:"status" validate:"oneof=pending completed failed"`
	QuizScore      *int       `json:"quiz_score"`
	QuizPassed     bool       `json:"quiz_passed"`
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
}

// CreateDocumentQuizDTO represents the request to create a quiz question
type CreateDocumentQuizDTO struct {
	Question      string   `json:"question" validate:"required,min=1"`
	QuestionOrder int      `json:"question_order" validate:"min=1"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correct_answer" validate:"required"`
}

// AnswerDocumentQuizDTO represents the request to answer a quiz question
type AnswerDocumentQuizDTO struct {
	QuizID string `json:"quiz_id" validate:"required"`
	Answer string `json:"answer" validate:"required"`
}

// DocumentFiltersDTO represents filters for listing documents
type DocumentFiltersDTO struct {
	Status   *string `json:"status"`
	Type     *string `json:"type"`
	Category *string `json:"category"`
	Search   *string `json:"search"`
	Page     int     `json:"page" validate:"min=1"`
	Limit    int     `json:"limit" validate:"min=1,max=100"`
}

// DocumentApprovalRouteDTO represents an approval route
type DocumentApprovalRouteDTO struct {
	RouteName string                    `json:"route_name" validate:"required,min=1,max=255"`
	Steps     []DocumentApprovalStepDTO `json:"steps" validate:"required,min=1"`
}

// DocumentApprovalStepDTO represents a step in an approval route
type DocumentApprovalStepDTO struct {
	StepOrder      int     `json:"step_order" validate:"min=1"`
	ApproverRoleID *string `json:"approver_role_id"`
	ApproverUserID *string `json:"approver_user_id"`
	IsRequired     bool    `json:"is_required"`
}

// DocumentApprovalActionDTO represents an approval action
type DocumentApprovalActionDTO struct {
	Status  string  `json:"status" validate:"required,oneof=approved rejected"`
	Comment *string `json:"comment"`
}

// DocumentStatsDTO represents document statistics
type DocumentStatsDTO struct {
	TotalDocuments    int            `json:"total_documents"`
	PendingApproval   int            `json:"pending_approval"`
	PendingAck        int            `json:"pending_ack"`
	OverdueAck        int            `json:"overdue_ack"`
	DocumentsByType   map[string]int `json:"documents_by_type"`
	DocumentsByStatus map[string]int `json:"documents_by_status"`
}

// DocumentSearchResultDTO represents a search result
type DocumentSearchResultDTO struct {
	DocumentID     string  `json:"document_id"`
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	Type           string  `json:"type"`
	Category       *string `json:"category"`
	Status         string  `json:"status"`
	Version        string  `json:"version"`
	CreatedAt      string  `json:"created_at"`
	RelevanceScore float64 `json:"relevance_score,omitempty"`
}

// SubmitDocumentDTO represents the request to submit document for approval
type SubmitDocumentDTO struct {
	WorkflowType string            `json:"workflow_type" validate:"oneof=sequential parallel"`
	Steps        []ApprovalStepDTO `json:"steps" validate:"required,min=1"`
}

// ApprovalStepDTO represents an approval step
type ApprovalStepDTO struct {
	StepOrder  int     `json:"step_order" validate:"min=1"`
	ApproverID string  `json:"approver_id" validate:"required"`
	Deadline   *string `json:"deadline"`
}

// ApprovalActionDTO represents an approval action
type ApprovalActionDTO struct {
	Action  string  `json:"action" validate:"oneof=approve reject"`
	Comment *string `json:"comment"`
}

// CreateACKCampaignDTO represents the request to create an ACK campaign
type CreateACKCampaignDTO struct {
	Title        string   `json:"title" validate:"required,min=1,max=255"`
	Description  *string  `json:"description"`
	AudienceType string   `json:"audience_type" validate:"oneof=all role department custom"`
	AudienceIDs  []string `json:"audience_ids"`
	Deadline     *string  `json:"deadline"`
	QuizID       *string  `json:"quiz_id"`
}

// CreateTrainingMaterialDTO represents the request to create a training material
type CreateTrainingMaterialDTO struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description *string `json:"description"`
	Type        string  `json:"type" validate:"oneof=document video presentation other"`
}

// CreateTrainingAssignmentDTO represents the request to assign training material
type CreateTrainingAssignmentDTO struct {
	UserIDs  []string `json:"user_ids" validate:"required,min=1"`
	Deadline *string  `json:"deadline"`
	QuizID   *string  `json:"quiz_id"`
}

// CreateQuizDTO represents the request to create a quiz
type CreateQuizDTO struct {
	Title            string            `json:"title" validate:"required,min=1,max=255"`
	Description      *string           `json:"description"`
	Questions        []QuizQuestionDTO `json:"questions" validate:"required,min=1"`
	PassingScore     int               `json:"passing_score" validate:"min=0,max=100"`
	TimeLimitMinutes *int              `json:"time_limit_minutes" validate:"min=1"`
}

// QuizQuestionDTO represents a quiz question
type QuizQuestionDTO struct {
	Question      string   `json:"question" validate:"required,min=1"`
	Options       []string `json:"options" validate:"required,min=2"`
	CorrectAnswer int      `json:"correct_answer" validate:"min=0"`
	Explanation   *string  `json:"explanation"`
}

// SubmitQuizAnswerDTO represents the request to submit quiz answers
type SubmitQuizAnswerDTO struct {
	Answers []QuizAnswerDTO `json:"answers" validate:"required,min=1"`
}

// QuizAnswerDTO represents a quiz answer
type QuizAnswerDTO struct {
	QuestionID string `json:"question_id" validate:"required"`
	Answer     int    `json:"answer" validate:"min=0"`
}
