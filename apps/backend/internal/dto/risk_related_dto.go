package dto

import "time"

// RiskControlRequest - запрос для добавления контроля к риску
type RiskControlRequest struct {
	ControlID            string  `json:"control_id" validate:"required,uuid4"`
	ControlName          string  `json:"control_name" validate:"required,min=1,max=255"`
	ControlType          string  `json:"control_type" validate:"required,oneof=preventive detective corrective"`
	ImplementationStatus string  `json:"implementation_status" validate:"required,oneof=planned in_progress implemented not_applicable"`
	Effectiveness        *string `json:"effectiveness,omitempty" validate:"omitempty,oneof=high medium low"`
	Description          *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

// RiskControlResponse - ответ с данными контроля риска
type RiskControlResponse struct {
	ID                   string    `json:"id"`
	RiskID               string    `json:"risk_id"`
	ControlID            string    `json:"control_id"`
	ControlName          string    `json:"control_name"`
	ControlType          string    `json:"control_type"`
	ImplementationStatus string    `json:"implementation_status"`
	Effectiveness        *string   `json:"effectiveness"`
	Description          *string   `json:"description"`
	CreatedBy            *string   `json:"created_by"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// RiskCommentRequest - запрос для добавления комментария к риску
type RiskCommentRequest struct {
	Comment    string `json:"comment" validate:"required,min=1,max=2000"`
	IsInternal *bool  `json:"is_internal,omitempty"`
}

// RiskCommentResponse - ответ с данными комментария риска
type RiskCommentResponse struct {
	ID         string    `json:"id"`
	RiskID     string    `json:"risk_id"`
	UserID     string    `json:"user_id"`
	Comment    string    `json:"comment"`
	IsInternal bool      `json:"is_internal"`
	UserName   *string   `json:"user_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// RiskHistoryResponse - ответ с данными истории риска
type RiskHistoryResponse struct {
	ID            string    `json:"id"`
	RiskID        string    `json:"risk_id"`
	FieldChanged  string    `json:"field_changed"`
	OldValue      *string   `json:"old_value"`
	NewValue      *string   `json:"new_value"`
	ChangeReason  *string   `json:"change_reason"`
	ChangedBy     string    `json:"changed_by"`
	ChangedAt     time.Time `json:"changed_at"`
	ChangedByName *string   `json:"changed_by_name"`
}

// RiskAttachmentRequest - запрос для добавления вложения к риску
type RiskAttachmentRequest struct {
	FileName    string  `json:"file_name" validate:"required,min=1,max=255"`
	FilePath    string  `json:"file_path" validate:"required,min=1,max=500"`
	FileSize    int64   `json:"file_size" validate:"required,min=1"`
	MimeType    string  `json:"mime_type" validate:"required,min=1,max=100"`
	FileHash    *string `json:"file_hash,omitempty" validate:"omitempty,len=64"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// RiskAttachmentResponse - ответ с данными вложения риска
type RiskAttachmentResponse struct {
	ID             string    `json:"id"`
	RiskID         string    `json:"risk_id"`
	FileName       string    `json:"file_name"`
	FilePath       string    `json:"file_path"`
	FileSize       int64     `json:"file_size"`
	MimeType       string    `json:"mime_type"`
	FileHash       *string   `json:"file_hash"`
	Description    *string   `json:"description"`
	UploadedBy     string    `json:"uploaded_by"`
	UploadedAt     time.Time `json:"uploaded_at"`
	UploadedByName *string   `json:"uploaded_by_name"`
}

// RiskTagRequest - запрос для добавления тега к риску
type RiskTagRequest struct {
	TagName  string  `json:"tag_name" validate:"required,min=1,max=100"`
	TagColor *string `json:"tag_color,omitempty" validate:"omitempty,len=7"`
}

// RiskTagResponse - ответ с данными тега риска
type RiskTagResponse struct {
	ID        string    `json:"id"`
	RiskID    string    `json:"risk_id"`
	TagName   string    `json:"tag_name"`
	TagColor  string    `json:"tag_color"`
	CreatedBy *string   `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// RiskDetailResponse - полный ответ с данными риска и связанными сущностями
type RiskDetailResponse struct {
	RiskResponse
	Controls    []RiskControlResponse    `json:"controls"`
	Comments    []RiskCommentResponse    `json:"comments"`
	History     []RiskHistoryResponse    `json:"history"`
	Attachments []RiskAttachmentResponse `json:"attachments"`
	Tags        []RiskTagResponse        `json:"tags"`
}

// RiskExportRequest - запрос для экспорта рисков
type RiskExportRequest struct {
	Format         string   `query:"format" validate:"required,oneof=csv xlsx pdf"`
	Status         []string `query:"status" validate:"omitempty,dive,oneof=new in_analysis in_treatment accepted transferred mitigated closed"`
	Category       []string `query:"category" validate:"omitempty,dive,max=100"`
	OwnerID        []string `query:"owner_id" validate:"omitempty,dive,uuid4"`
	Level          []string `query:"level" validate:"omitempty,dive,oneof=Low Medium High Critical"`
	IncludeDetails bool     `query:"include_details" validate:"omitempty"`
}

// RiskExportResponse - ответ с данными экспорта
type RiskExportResponse struct {
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	MimeType    string `json:"mime_type"`
	DownloadURL string `json:"download_url"`
}

// Constants for risk related entities
const (
	// Control types
	ControlTypePreventive = "preventive"
	ControlTypeDetective  = "detective"
	ControlTypeCorrective = "corrective"

	// Implementation status
	ImplementationStatusPlanned       = "planned"
	ImplementationStatusInProgress    = "in_progress"
	ImplementationStatusImplemented   = "implemented"
	ImplementationStatusNotApplicable = "not_applicable"

	// Effectiveness levels
	EffectivenessHigh   = "high"
	EffectivenessMedium = "medium"
	EffectivenessLow    = "low"

	// Export formats
	ExportFormatCSV  = "csv"
	ExportFormatXLSX = "xlsx"
	ExportFormatPDF  = "pdf"

	// Default colors for tags
	TagColorDefault   = "#007bff"
	TagColorSuccess   = "#28a745"
	TagColorWarning   = "#ffc107"
	TagColorDanger    = "#dc3545"
	TagColorInfo      = "#17a2b8"
	TagColorSecondary = "#6c757d"
)

