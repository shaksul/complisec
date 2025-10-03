package dto

import (
	"time"
)

// CalculateRiskLevel - вычисляет уровень риска на основе likelihood и impact (1-4 шкала)
func CalculateRiskLevel(likelihood, impact int) (int, string) {
	level := likelihood * impact

	var label string
	switch {
	case level <= 2:
		label = RiskLevelLabelLow
	case level <= 4:
		label = RiskLevelLabelMedium
	case level <= 6:
		label = RiskLevelLabelHigh
	default: // 7-8
		label = RiskLevelLabelCritical
	}

	return level, label
}

// RiskRequest - базовый запрос для риска
type RiskRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Category    *string `json:"category,omitempty" validate:"omitempty,max=100"`
	Likelihood  int     `json:"likelihood" validate:"required,min=1,max=4"`
	Impact      int     `json:"impact" validate:"required,min=1,max=4"`
	OwnerUserID *string `json:"owner_user_id,omitempty" validate:"omitempty,uuid4"`
	AssetID     *string `json:"asset_id,omitempty" validate:"omitempty,uuid4"`
	Methodology *string `json:"methodology,omitempty" validate:"omitempty,oneof=ISO27005 NIST COSO Custom"`
	Strategy    *string `json:"strategy,omitempty" validate:"omitempty,oneof=accept mitigate transfer avoid"`
	DueDate     *string `json:"due_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// CreateRiskRequest - запрос на создание риска
type CreateRiskRequest struct {
	RiskRequest
}

// UpdateRiskRequest - запрос на обновление риска
type UpdateRiskRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Category    *string `json:"category,omitempty" validate:"omitempty,max=100"`
	Likelihood  *int    `json:"likelihood,omitempty" validate:"omitempty,min=1,max=4"`
	Impact      *int    `json:"impact,omitempty" validate:"omitempty,min=1,max=4"`
	OwnerUserID *string `json:"owner_user_id,omitempty" validate:"omitempty,uuid4"`
	AssetID     *string `json:"asset_id,omitempty" validate:"omitempty,uuid4"`
	Methodology *string `json:"methodology,omitempty" validate:"omitempty,oneof=ISO27005 NIST COSO Custom"`
	Strategy    *string `json:"strategy,omitempty" validate:"omitempty,oneof=accept mitigate transfer avoid"`
	DueDate     *string `json:"due_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// RiskResponse - ответ с данными риска
type RiskResponse struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Category    *string    `json:"category"`
	Likelihood  *int       `json:"likelihood"`
	Impact      *int       `json:"impact"`
	Level       *int       `json:"level"`
	Status      string     `json:"status"`
	OwnerUserID *string    `json:"owner_user_id"`
	AssetID     *string    `json:"asset_id"`
	Methodology *string    `json:"methodology"`
	Strategy    *string    `json:"strategy"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	OwnerName   *string    `json:"owner_name,omitempty"`
	AssetName   *string    `json:"asset_name,omitempty"`
	LevelLabel  *string    `json:"level_label,omitempty"`
}

// RiskListRequest - запрос на получение списка рисков
type RiskListRequest struct {
	Page        int    `query:"page" validate:"min=1"`
	PageSize    int    `query:"page_size" validate:"min=1,max=100"`
	Status      string `query:"status" validate:"omitempty,oneof=new in_analysis in_treatment accepted transferred mitigated closed"`
	Category    string `query:"category" validate:"omitempty,max=100"`
	AssetID     string `query:"asset_id" validate:"omitempty,uuid4"`
	OwnerUserID string `query:"owner_user_id" validate:"omitempty,uuid4"`
	Methodology string `query:"methodology" validate:"omitempty,oneof=ISO27005 NIST COSO Custom"`
	Strategy    string `query:"strategy" validate:"omitempty,oneof=accept mitigate transfer avoid"`
	Level       string `query:"level" validate:"omitempty,oneof=Low Medium High Critical"`
	Search      string `query:"search" validate:"omitempty,max=255"`
}

// RiskStatusUpdateRequest - запрос на обновление статуса риска
type RiskStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=new in_analysis in_treatment accepted transferred mitigated closed"`
}

// Risk level constants - updated for 1-4 scale
const (
	RiskLevelLow      = 1 // 1-2
	RiskLevelMedium   = 2 // 3-4
	RiskLevelHigh     = 3 // 5-6
	RiskLevelCritical = 4 // 7-8
)

// Risk level labels
const (
	RiskLevelLabelLow      = "Low"
	RiskLevelLabelMedium   = "Medium"
	RiskLevelLabelHigh     = "High"
	RiskLevelLabelCritical = "Critical"
)

// Risk status constants - new statuses
const (
	RiskStatusNew         = "new"
	RiskStatusInAnalysis  = "in_analysis"
	RiskStatusInTreatment = "in_treatment"
	RiskStatusAccepted    = "accepted"
	RiskStatusTransferred = "transferred"
	RiskStatusMitigated   = "mitigated"
	RiskStatusClosed      = "closed"
)

// Risk category constants
const (
	RiskCategoryTechnical    = "technical"
	RiskCategoryOperational  = "operational"
	RiskCategoryCompliance   = "compliance"
	RiskCategoryFinancial    = "financial"
	RiskCategoryReputational = "reputational"
)

// Risk methodology constants
const (
	RiskMethodologyISO27005 = "ISO27005"
	RiskMethodologyNIST     = "NIST"
	RiskMethodologyCOSO     = "COSO"
	RiskMethodologyCustom   = "Custom"
)

// Risk strategy constants
const (
	RiskStrategyAccept   = "accept"
	RiskStrategyMitigate = "mitigate"
	RiskStrategyTransfer = "transfer"
	RiskStrategyAvoid    = "avoid"
)
