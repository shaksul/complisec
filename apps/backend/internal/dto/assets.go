package dto

import "time"

// Asset types
const (
	AssetTypeServer        = "server"
	AssetTypeWorkstation   = "workstation"
	AssetTypeApplication   = "application"
	AssetTypeDatabase      = "database"
	AssetTypeDocument      = "document"
	AssetTypeNetworkDevice = "network_device"
	AssetTypeOther         = "other"
)

// Asset classes
const (
	AssetClassHardware = "hardware"
	AssetClassSoftware = "software"
	AssetClassData     = "data"
	AssetClassService  = "service"
)

// Asset criticality levels
const (
	CriticalityLow    = "low"
	CriticalityMedium = "medium"
	CriticalityHigh   = "high"
)

// Asset statuses
const (
	AssetStatusActive         = "active"
	AssetStatusInRepair       = "in_repair"
	AssetStatusStorage        = "storage"
	AssetStatusDecommissioned = "decommissioned"
)

// Document types
const (
	DocumentTypePassport    = "passport"
	DocumentTypeTransferAct = "transfer_act"
	DocumentTypeWriteoffAct = "writeoff_act"
	DocumentTypeRepairLog   = "repair_log"
	DocumentTypeOther       = "other"
)

// CreateAssetRequest represents the request to create a new asset
type CreateAssetRequest struct {
	Name              string `json:"name" validate:"required,min=1,max=255"`
	Type              string `json:"type" validate:"required,oneof=server workstation application database document network_device other"`
	Class             string `json:"class" validate:"required,oneof=hardware software data service"`
	OwnerID           string `json:"owner_id" validate:"omitempty,uuid"`
	ResponsibleUserID string `json:"responsible_user_id" validate:"omitempty,uuid"`
	Location          string `json:"location,omitempty"`
	Criticality       string `json:"criticality" validate:"required,oneof=low medium high"`
	Confidentiality   string `json:"confidentiality" validate:"required,oneof=low medium high"`
	Integrity         string `json:"integrity" validate:"required,oneof=low medium high"`
	Availability      string `json:"availability" validate:"required,oneof=low medium high"`
	Status            string `json:"status,omitempty" validate:"omitempty,oneof=active in_repair storage decommissioned"`
	// Passport fields
	SerialNumber  *string `json:"serial_number,omitempty" validate:"omitempty,max=255"`
	PCNumber      *string `json:"pc_number,omitempty" validate:"omitempty,max=100"`
	Model         *string `json:"model,omitempty" validate:"omitempty,max=255"`
	CPU           *string `json:"cpu,omitempty" validate:"omitempty,max=255"`
	RAM           *string `json:"ram,omitempty" validate:"omitempty,max=100"`
	HDDInfo       *string `json:"hdd_info,omitempty"`
	NetworkCard   *string `json:"network_card,omitempty" validate:"omitempty,max=255"`
	OpticalDrive  *string `json:"optical_drive,omitempty" validate:"omitempty,max=255"`
	IPAddress     *string `json:"ip_address,omitempty"`
	MACAddress    *string `json:"mac_address,omitempty"`
	Manufacturer  *string `json:"manufacturer,omitempty" validate:"omitempty,max=255"`
	PurchaseYear  *int    `json:"purchase_year,omitempty" validate:"omitempty,gte=1900,lte=2100"`
	WarrantyUntil *string `json:"warranty_until,omitempty"` // Date string
}

// UpdateAssetRequest represents the request to update an existing asset
type UpdateAssetRequest struct {
	Name              *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Type              *string `json:"type,omitempty" validate:"omitempty,oneof=server workstation application database document network_device other"`
	Class             *string `json:"class,omitempty" validate:"omitempty,oneof=hardware software data service"`
	OwnerID           *string `json:"owner_id,omitempty" validate:"omitempty,uuid"`
	ResponsibleUserID *string `json:"responsible_user_id,omitempty" validate:"omitempty,uuid"`
	Location          *string `json:"location,omitempty"`
	Criticality       *string `json:"criticality,omitempty" validate:"omitempty,oneof=low medium high"`
	Confidentiality   *string `json:"confidentiality,omitempty" validate:"omitempty,oneof=low medium high"`
	Integrity         *string `json:"integrity,omitempty" validate:"omitempty,oneof=low medium high"`
	Availability      *string `json:"availability,omitempty" validate:"omitempty,oneof=low medium high"`
	Status            *string `json:"status,omitempty" validate:"omitempty,oneof=active in_repair storage decommissioned"`
	// Passport fields
	SerialNumber  *string `json:"serial_number,omitempty" validate:"omitempty,max=255"`
	PCNumber      *string `json:"pc_number,omitempty" validate:"omitempty,max=100"`
	Model         *string `json:"model,omitempty" validate:"omitempty,max=255"`
	CPU           *string `json:"cpu,omitempty" validate:"omitempty,max=255"`
	RAM           *string `json:"ram,omitempty" validate:"omitempty,max=100"`
	HDDInfo       *string `json:"hdd_info,omitempty"`
	NetworkCard   *string `json:"network_card,omitempty" validate:"omitempty,max=255"`
	OpticalDrive  *string `json:"optical_drive,omitempty" validate:"omitempty,max=255"`
	IPAddress     *string `json:"ip_address,omitempty"`
	MACAddress    *string `json:"mac_address,omitempty"`
	Manufacturer  *string `json:"manufacturer,omitempty" validate:"omitempty,max=255"`
	PurchaseYear  *int    `json:"purchase_year,omitempty" validate:"omitempty,gte=1900,lte=2100"`
	WarrantyUntil *string `json:"warranty_until,omitempty"` // Date string
}

// AssetResponse represents the response for asset data
type AssetResponse struct {
	ID                  string    `json:"id"`
	InventoryNumber     string    `json:"inventory_number"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Class               string    `json:"class"`
	OwnerID             *string   `json:"owner_id,omitempty"`
	OwnerName           *string   `json:"owner_name,omitempty"`
	ResponsibleUserID   *string   `json:"responsible_user_id,omitempty"`
	ResponsibleUserName *string   `json:"responsible_user_name,omitempty"`
	Location            *string   `json:"location,omitempty"`
	Criticality         string    `json:"criticality"`
	Confidentiality     string    `json:"confidentiality"`
	Integrity           string    `json:"integrity"`
	Availability        string    `json:"availability"`
	Status              string    `json:"status"`
	// Passport fields
	SerialNumber  *string `json:"serial_number,omitempty"`
	PCNumber      *string `json:"pc_number,omitempty"`
	Model         *string `json:"model,omitempty"`
	CPU           *string `json:"cpu,omitempty"`
	RAM           *string `json:"ram,omitempty"`
	HDDInfo       *string `json:"hdd_info,omitempty"`
	NetworkCard   *string `json:"network_card,omitempty"`
	OpticalDrive  *string `json:"optical_drive,omitempty"`
	IPAddress     *string `json:"ip_address,omitempty"`
	MACAddress    *string `json:"mac_address,omitempty"`
	Manufacturer  *string `json:"manufacturer,omitempty"`
	PurchaseYear  *int    `json:"purchase_year,omitempty"`
	WarrantyUntil *string `json:"warranty_until,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AssetDocumentRequest represents the request to add a document to an asset
type AssetDocumentRequest struct {
	DocumentType string `json:"document_type" validate:"required,oneof=passport transfer_act writeoff_act repair_log other"`
	FilePath     string `json:"file_path" validate:"required"`
}

// AssetDocumentUploadRequest represents the request to upload a new document file
type AssetDocumentUploadRequest struct {
	DocumentType string `form:"document_type" validate:"required,oneof=passport transfer_act writeoff_act repair_log other"`
	Title        string `form:"title" validate:"omitempty,max=255"`
	File         []byte `form:"file" validate:"required"`
}

// AssetDocumentLinkRequest represents the request to link an existing document to an asset
type AssetDocumentLinkRequest struct {
	DocumentID   string `json:"document_id" validate:"required,uuid"`
	DocumentType string `json:"document_type" validate:"required,oneof=passport transfer_act writeoff_act repair_log other"`
}

// AssetDocumentResponse represents the response for asset document data
type AssetDocumentResponse struct {
	ID           string    `json:"id"`
	AssetID      string    `json:"asset_id"`
	Title        string    `json:"title"`
	DocumentType string    `json:"document_type"`
	Mime         string    `json:"mime"`
	SizeBytes    int64     `json:"size_bytes"`
	DownloadURL  string    `json:"download_url"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// DocumentStorageRequest represents the request to list documents from storage
type DocumentStorageRequest struct {
	Query    string `json:"query" validate:"omitempty,max=255"`
	Type     string `json:"type" validate:"omitempty,oneof=passport transfer_act writeoff_act repair_log other"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

// DocumentStorageResponse represents the response for document storage data
type DocumentStorageResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	DocumentType string    `json:"document_type"`
	Version      string    `json:"version"`
	SizeBytes    int64     `json:"size_bytes"`
	Mime         string    `json:"mime"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// AssetSoftwareRequest represents the request to add software to an asset
type AssetSoftwareRequest struct {
	SoftwareName string     `json:"software_name" validate:"required,min=1,max=255"`
	Version      *string    `json:"version,omitempty" validate:"omitempty,max=100"`
	InstalledAt  *time.Time `json:"installed_at,omitempty"`
}

// AssetSoftwareResponse represents the response for asset software data
type AssetSoftwareResponse struct {
	ID           string     `json:"id"`
	AssetID      string     `json:"asset_id"`
	SoftwareName string     `json:"software_name"`
	Version      *string    `json:"version,omitempty"`
	InstalledAt  *time.Time `json:"installed_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// AssetHistoryResponse represents the response for asset history data
type AssetHistoryResponse struct {
	ID           string    `json:"id"`
	AssetID      string    `json:"asset_id"`
	FieldChanged string    `json:"field_changed"`
	OldValue     *string   `json:"old_value,omitempty"`
	NewValue     string    `json:"new_value"`
	ChangedBy    string    `json:"changed_by"`
	ChangedAt    time.Time `json:"changed_at"`
}

// AssetListRequest represents the request parameters for listing assets
type AssetListRequest struct {
	Page        int    `json:"page" validate:"min=1"`
	PageSize    int    `json:"page_size" validate:"min=1,max=1000"`
	Type        string `json:"type,omitempty" validate:"omitempty,oneof=server workstation application database document network_device other"`
	Class       string `json:"class,omitempty" validate:"omitempty,oneof=hardware software data service"`
	Status      string `json:"status,omitempty" validate:"omitempty,oneof=active in_repair storage decommissioned"`
	Criticality string `json:"criticality,omitempty" validate:"omitempty,oneof=low medium high"`
	OwnerID     string `json:"owner_id,omitempty" validate:"omitempty,uuid"`
	Search      string `json:"search,omitempty"`
}

// AssetInventoryRequest represents the request to perform inventory
type AssetInventoryRequest struct {
	AssetIDs []string `json:"asset_ids" validate:"required,min=1,dive,uuid"`
	Action   string   `json:"action" validate:"required,oneof=verify update_status"`
	Status   *string  `json:"status,omitempty" validate:"omitempty,oneof=active in_repair storage decommissioned"`
	Notes    *string  `json:"notes,omitempty"`
}

// AssetHistoryFiltersRequest represents the request to filter asset history
type AssetHistoryFiltersRequest struct {
	ChangedBy *string `json:"changed_by,omitempty" validate:"omitempty,uuid"`
	FromDate  *string `json:"from_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	ToDate    *string `json:"to_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// BulkUpdateStatusRequest represents the request to update status for multiple assets
type BulkUpdateStatusRequest struct {
	AssetIDs []string `json:"asset_ids" validate:"required,min=1,dive,uuid"`
	Status   string   `json:"status" validate:"required,oneof=active in_repair storage decommissioned"`
}

// BulkUpdateOwnerRequest represents the request to update owner for multiple assets
type BulkUpdateOwnerRequest struct {
	AssetIDs []string `json:"asset_ids" validate:"required,min=1,dive,uuid"`
	OwnerID  string   `json:"owner_id" validate:"omitempty,uuid"`
}
