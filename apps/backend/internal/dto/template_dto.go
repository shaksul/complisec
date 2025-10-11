package dto

import "time"

// DocumentTemplateDTO represents a document template
type DocumentTemplateDTO struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	TemplateType string    `json:"template_type"`
	Content      string    `json:"content"`
	IsSystem     bool      `json:"is_system"`
	IsActive     bool      `json:"is_active"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateTemplateRequest for creating a new template
type CreateTemplateRequest struct {
	Name         string  `json:"name" validate:"required,min=3,max=255"`
	Description  *string `json:"description"`
	TemplateType string  `json:"template_type" validate:"required,oneof=passport_pc passport_monitor passport_device transfer_act writeoff_act repair_log other"`
	Content      string  `json:"content" validate:"required"`
}

// UpdateTemplateRequest for updating an existing template
type UpdateTemplateRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=3,max=255"`
	Description  *string `json:"description"`
	TemplateType *string `json:"template_type" validate:"omitempty,oneof=passport_pc passport_monitor passport_device transfer_act writeoff_act repair_log other"`
	Content      *string `json:"content"`
	IsActive     *bool   `json:"is_active"`
}

// FillTemplateRequest for filling a template with asset data
type FillTemplateRequest struct {
	TemplateID     string                 `json:"template_id" validate:"required,uuid"`
	AssetID        string                 `json:"asset_id" validate:"required,uuid"`
	AdditionalData map[string]interface{} `json:"additional_data"`  // Extra fields not in asset model
	SaveAsDocument bool                   `json:"save_as_document"` // If true, save filled template as document
	DocumentTitle  string                 `json:"document_title"`   // Title for saved document
	GeneratePDF    bool                   `json:"generate_pdf"`     // If true, generate PDF instead of HTML
}

// FillTemplateResponse returns the filled template
type FillTemplateResponse struct {
	HTML       string  `json:"html,omitempty"`         // Filled HTML template
	PDFBase64  *string `json:"pdf_base64,omitempty"`   // PDF file as base64 if generate_pdf is true
	DocumentID *string `json:"document_id,omitempty"`  // ID of saved document if save_as_document is true
}

// InventoryNumberRuleDTO represents a rule for generating inventory numbers
type InventoryNumberRuleDTO struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	AssetType       string    `json:"asset_type"`
	AssetClass      *string   `json:"asset_class,omitempty"`
	Pattern         string    `json:"pattern"`
	CurrentSequence int       `json:"current_sequence"`
	Prefix          *string   `json:"prefix,omitempty"`
	Description     *string   `json:"description,omitempty"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateInventoryRuleRequest for creating a new rule
type CreateInventoryRuleRequest struct {
	AssetType   string  `json:"asset_type" validate:"required,max=50"`
	AssetClass  *string `json:"asset_class" validate:"omitempty,max=50"`
	Pattern     string  `json:"pattern" validate:"required,max=255"`
	Prefix      *string `json:"prefix" validate:"omitempty,max=50"`
	Description *string `json:"description"`
}

// UpdateInventoryRuleRequest for updating an existing rule
type UpdateInventoryRuleRequest struct {
	Pattern         *string `json:"pattern" validate:"omitempty,max=255"`
	CurrentSequence *int    `json:"current_sequence" validate:"omitempty,gte=0"`
	Prefix          *string `json:"prefix" validate:"omitempty,max=50"`
	Description     *string `json:"description"`
	IsActive        *bool   `json:"is_active"`
}

// GenerateInventoryNumberRequest for generating a new inventory number
type GenerateInventoryNumberRequest struct {
	AssetType  string  `json:"asset_type" validate:"required"`
	AssetClass *string `json:"asset_class"`
}

// GenerateInventoryNumberResponse returns the generated number
type GenerateInventoryNumberResponse struct {
	InventoryNumber string `json:"inventory_number"`
	Pattern         string `json:"pattern"`
	Sequence        int    `json:"sequence"`
}

// AssetPassportDataRequest for updating asset passport-specific fields
type AssetPassportDataRequest struct {
	SerialNumber  *string                `json:"serial_number" validate:"omitempty,max=255"`
	PCNumber      *string                `json:"pc_number" validate:"omitempty,max=100"`
	Model         *string                `json:"model" validate:"omitempty,max=255"`
	CPU           *string                `json:"cpu" validate:"omitempty,max=255"`
	RAM           *string                `json:"ram" validate:"omitempty,max=100"`
	HDDInfo       *string                `json:"hdd_info"`
	NetworkCard   *string                `json:"network_card" validate:"omitempty,max=255"`
	OpticalDrive  *string                `json:"optical_drive" validate:"omitempty,max=255"`
	IPAddress     *string                `json:"ip_address" validate:"omitempty,ip"`
	MACAddress    *string                `json:"mac_address" validate:"omitempty,mac"`
	Manufacturer  *string                `json:"manufacturer" validate:"omitempty,max=255"`
	PurchaseYear  *int                   `json:"purchase_year" validate:"omitempty,gte=1900,lte=2100"`
	WarrantyUntil *string                `json:"warranty_until"` // Date string
	Metadata      map[string]interface{} `json:"metadata"`
}

// TemplateVariablesResponse lists available variables for templates
type TemplateVariablesResponse struct {
	Variables []TemplateVariable `json:"variables"`
}

// TemplateVariable describes a single template variable
type TemplateVariable struct {
	Name        string `json:"name"`
	Placeholder string `json:"placeholder"` // e.g. "{{asset_name}}"
	Description string `json:"description"`
	Example     string `json:"example"`
	Category    string `json:"category"` // "asset", "user", "date", "custom"
}

