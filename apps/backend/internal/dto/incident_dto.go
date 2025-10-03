package dto

import "time"

// IncidentRequest - базовый запрос для инцидента
type IncidentRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=2000"`
	Category    string     `json:"category" validate:"required,oneof=technical_failure data_breach unauthorized_access physical malware social_engineering other"`
	Criticality string     `json:"criticality" validate:"required,oneof=low medium high critical"`
	Source      string     `json:"source" validate:"required,oneof=user_report automatic_agent admin_manual monitoring siem"`
	AssetIDs    []string   `json:"asset_ids,omitempty" validate:"omitempty,dive,uuid4"`
	RiskIDs     []string   `json:"risk_ids,omitempty" validate:"omitempty,dive,uuid4"`
	AssignedTo  *string    `json:"assigned_to,omitempty" validate:"omitempty,uuid4"`
	DetectedAt  *time.Time `json:"detected_at,omitempty"`
}

// CreateIncidentRequest - запрос на создание инцидента
type CreateIncidentRequest struct {
	IncidentRequest
}

// UpdateIncidentRequest - запрос на обновление инцидента
type UpdateIncidentRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=2000"`
	Category    *string    `json:"category,omitempty" validate:"omitempty,oneof=technical_failure data_breach unauthorized_access physical malware social_engineering other"`
	Criticality *string    `json:"criticality,omitempty" validate:"omitempty,oneof=low medium high critical"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=new assigned in_progress resolved closed"`
	AssetIDs    []string   `json:"asset_ids,omitempty" validate:"omitempty,dive,uuid4"`
	RiskIDs     []string   `json:"risk_ids,omitempty" validate:"omitempty,dive,uuid4"`
	AssignedTo  *string    `json:"assigned_to,omitempty" validate:"omitempty,uuid4"`
	DetectedAt  *time.Time `json:"detected_at,omitempty"`
}

// IncidentResponse - ответ с данными инцидента
type IncidentResponse struct {
	ID           string      `json:"id"`
	TenantID     string      `json:"tenant_id"`
	Title        string      `json:"title"`
	Description  *string     `json:"description"`
	Category     string      `json:"category"`
	Status       string      `json:"status"`
	Criticality  string      `json:"criticality"`
	Source       string      `json:"source"`
	ReportedBy   string      `json:"reported_by"`
	AssignedTo   *string     `json:"assigned_to"`
	DetectedAt   time.Time   `json:"detected_at"`
	ResolvedAt   *time.Time  `json:"resolved_at,omitempty"`
	ClosedAt     *time.Time  `json:"closed_at,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Assets       []AssetInfo `json:"assets,omitempty"`
	Risks        []RiskInfo  `json:"risks,omitempty"`
	ReportedName *string     `json:"reported_name,omitempty"`
	AssignedName *string     `json:"assigned_name,omitempty"`
}

// AssetInfo - информация об активе
type AssetInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RiskInfo - информация о риске
type RiskInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IncidentListRequest - запрос на получение списка инцидентов
type IncidentListRequest struct {
	Page        int    `query:"page" validate:"min=1"`
	PageSize    int    `query:"page_size" validate:"min=1,max=100"`
	Status      string `query:"status" validate:"omitempty,oneof=new assigned in_progress resolved closed"`
	Criticality string `query:"criticality" validate:"omitempty,oneof=low medium high critical"`
	Category    string `query:"category" validate:"omitempty,oneof=technical_failure data_breach unauthorized_access physical malware social_engineering other"`
	AssetID     string `query:"asset_id" validate:"omitempty,uuid4"`
	RiskID      string `query:"risk_id" validate:"omitempty,uuid4"`
	AssignedTo  string `query:"assigned_to" validate:"omitempty,uuid4"`
	Search      string `query:"search" validate:"omitempty,max=100"`
}

// IncidentStatusUpdateRequest - запрос на обновление статуса инцидента
type IncidentStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=new assigned in_progress resolved closed"`
}

// IncidentCommentRequest - запрос на добавление комментария
type IncidentCommentRequest struct {
	Comment    string `json:"comment" validate:"required,min=1,max=2000"`
	IsInternal bool   `json:"is_internal"`
}

// IncidentCommentResponse - ответ с комментарием
type IncidentCommentResponse struct {
	ID         string    `json:"id"`
	IncidentID string    `json:"incident_id"`
	UserID     string    `json:"user_id"`
	Comment    string    `json:"comment"`
	IsInternal bool      `json:"is_internal"`
	CreatedAt  time.Time `json:"created_at"`
	UserName   *string   `json:"user_name,omitempty"`
}

// IncidentActionRequest - запрос на создание корректирующего действия
type IncidentActionRequest struct {
	ActionType  string     `json:"action_type" validate:"required,oneof=investigation containment eradication recovery prevention"`
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=2000"`
	AssignedTo  *string    `json:"assigned_to,omitempty" validate:"omitempty,uuid4"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// IncidentActionResponse - ответ с корректирующим действием
type IncidentActionResponse struct {
	ID           string     `json:"id"`
	IncidentID   string     `json:"incident_id"`
	ActionType   string     `json:"action_type"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	AssignedTo   *string    `json:"assigned_to"`
	DueDate      *time.Time `json:"due_date"`
	CompletedAt  *time.Time `json:"completed_at"`
	Status       string     `json:"status"`
	CreatedBy    string     `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	AssignedName *string    `json:"assigned_name,omitempty"`
	CreatedName  *string    `json:"created_name,omitempty"`
}

// IncidentMetricsResponse - ответ с метриками инцидента
type IncidentMetricsResponse struct {
	TotalIncidents  int            `json:"total_incidents"`
	OpenIncidents   int            `json:"open_incidents"`
	ClosedIncidents int            `json:"closed_incidents"`
	AverageMTTR     float64        `json:"average_mttr_hours"`
	AverageMTTD     float64        `json:"average_mttd_hours"`
	ByCriticality   map[string]int `json:"by_criticality"`
	ByCategory      map[string]int `json:"by_category"`
	ByStatus        map[string]int `json:"by_status"`
}

// Incident constants
const (
	// Categories
	IncidentCategoryTechnicalFailure   = "technical_failure"
	IncidentCategoryDataBreach         = "data_breach"
	IncidentCategoryUnauthorizedAccess = "unauthorized_access"
	IncidentCategoryPhysical           = "physical"
	IncidentCategoryMalware            = "malware"
	IncidentCategorySocialEngineering  = "social_engineering"
	IncidentCategoryOther              = "other"

	// Criticality levels
	IncidentCriticalityLow      = "low"
	IncidentCriticalityMedium   = "medium"
	IncidentCriticalityHigh     = "high"
	IncidentCriticalityCritical = "critical"

	// Status values
	IncidentStatusNew        = "new"
	IncidentStatusAssigned   = "assigned"
	IncidentStatusInProgress = "in_progress"
	IncidentStatusResolved   = "resolved"
	IncidentStatusClosed     = "closed"

	// Sources
	IncidentSourceUserReport     = "user_report"
	IncidentSourceAutomaticAgent = "automatic_agent"
	IncidentSourceAdminManual    = "admin_manual"
	IncidentSourceMonitoring     = "monitoring"
	IncidentSourceSIEM           = "siem"

	// Action types
	ActionTypeInvestigation = "investigation"
	ActionTypeContainment   = "containment"
	ActionTypeEradication   = "eradication"
	ActionTypeRecovery      = "recovery"
	ActionTypePrevention    = "prevention"

	// Action status
	ActionStatusPending    = "pending"
	ActionStatusInProgress = "in_progress"
	ActionStatusCompleted  = "completed"
	ActionStatusCancelled  = "cancelled"
)
