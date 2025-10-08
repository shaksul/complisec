package dto

import "time"

// CreateFolderDTO represents the request to create a folder
type CreateFolderDTO struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
	Metadata    *string `json:"metadata"`
}

// UpdateFolderDTO represents the request to update a folder
type UpdateFolderDTO struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description"`
	Metadata    *string `json:"metadata"`
}

// FolderDTO represents a folder response
type FolderDTO struct {
	ID          string      `json:"id"`
	TenantID    string      `json:"tenant_id"`
	Name        string      `json:"name"`
	Description *string     `json:"description"`
	ParentID    *string     `json:"parent_id"`
	OwnerID     string      `json:"owner_id"`
	CreatedBy   string      `json:"created_by"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	IsActive    bool        `json:"is_active"`
	Metadata    *string     `json:"metadata"`
	Children    []FolderDTO `json:"children,omitempty"`
}

// UploadDocumentDTO represents the request to upload a document
type UploadDocumentDTO struct {
	Name        string           `json:"name" validate:"required,min=1,max=255"`
	Description *string          `json:"description"`
	FolderID    *string          `json:"folder_id"`
	Tags        []string         `json:"tags"`
	LinkedTo    *DocumentLinkDTO `json:"linked_to"`
	EnableOCR   bool             `json:"enable_ocr"`
	Metadata    *string          `json:"metadata"`
}

// UpdateFileDocumentDTO represents the request to update a file document
type UpdateFileDocumentDTO struct {
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	Description *string  `json:"description"`
	FolderID    *string  `json:"folder_id"`
	Tags        []string `json:"tags"`
	Metadata    *string  `json:"metadata"`
}

// DocumentLinkDTO represents a link to another module entity
type DocumentLinkDTO struct {
	Module   string `json:"module" validate:"required,oneof=risk asset incident training compliance"`
	EntityID string `json:"entity_id" validate:"required"`
}

// DocumentDTO represents a document response
type DocumentDTO struct {
	ID           string            `json:"id"`
	TenantID     string            `json:"tenant_id"`
	Title        string            `json:"title"`
	OriginalName string            `json:"original_name"`
	Description  *string           `json:"description"`
	Type         string            `json:"type"`
	Category     *string           `json:"category"`
	FilePath     string            `json:"file_path"`
	FileSize     int64             `json:"file_size"`
	MimeType     string            `json:"mime_type"`
	FileHash     string            `json:"file_hash"`
	FolderID     *string           `json:"folder_id"`
	OwnerID      string            `json:"owner_id"`
	CreatedBy    string            `json:"created_by"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	IsActive     bool              `json:"is_active"`
	Version      string            `json:"version"`
	Metadata     *string           `json:"metadata"`
	Tags         []string          `json:"tags"`
	Links        []DocumentLinkDTO `json:"links"`
	OCRText      *string           `json:"ocr_text,omitempty"`
}

// DocumentPermissionDTO represents document permission
type DocumentPermissionDTO struct {
	ID          string     `json:"id"`
	SubjectType string     `json:"subject_type"`
	SubjectID   string     `json:"subject_id"`
	ObjectType  string     `json:"object_type"`
	ObjectID    string     `json:"object_id"`
	Permission  string     `json:"permission"`
	GrantedBy   string     `json:"granted_by"`
	GrantedAt   time.Time  `json:"granted_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsActive    bool       `json:"is_active"`
}

// CreateDocumentPermissionDTO represents the request to create document permission
type CreateDocumentPermissionDTO struct {
	SubjectType string     `json:"subject_type" validate:"required,oneof=user role"`
	SubjectID   string     `json:"subject_id" validate:"required"`
	ObjectType  string     `json:"object_type" validate:"required,oneof=document folder"`
	ObjectID    string     `json:"object_id" validate:"required"`
	Permission  string     `json:"permission" validate:"required,oneof=view edit delete share"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// DocumentVersionDTO represents a document version
type DocumentVersionDTO struct {
	ID                string    `json:"id"`
	DocumentID        string    `json:"document_id"`
	VersionNumber     int       `json:"version_number"`
	FilePath          string    `json:"file_path"`
	FileSize          int64     `json:"file_size"`
	FileHash          string    `json:"file_hash"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	ChangeDescription *string   `json:"change_description"`
}

// DocumentAuditLogDTO represents document audit log entry
type DocumentAuditLogDTO struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	DocumentID *string   `json:"document_id"`
	FolderID   *string   `json:"folder_id"`
	UserID     string    `json:"user_id"`
	Action     string    `json:"action"`
	Details    *string   `json:"details"`
	IPAddress  *string   `json:"ip_address"`
	UserAgent  *string   `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// FileDocumentFiltersDTO represents filters for listing file documents
type FileDocumentFiltersDTO struct {
	FolderID  *string  `json:"folder_id"`
	Tags      []string `json:"tags"`
	MimeType  *string  `json:"mime_type"`
	OwnerID   *string  `json:"owner_id"`
	Search    *string  `json:"search"`
	DateFrom  *string  `json:"date_from"`
	DateTo    *string  `json:"date_to"`
	Module    *string  `json:"module"`
	EntityID  *string  `json:"entity_id"`
	Page      int      `json:"page" validate:"min=1"`
	Limit     int      `json:"limit" validate:"min=1,max=100"`
	SortBy    *string  `json:"sort_by"`
	SortOrder *string  `json:"sort_order"`
}

// FileDocumentSearchResultDTO represents a search result for file documents
type FileDocumentSearchResultDTO struct {
	DocumentID     string  `json:"document_id"`
	Name           string  `json:"name"`
	Description    *string `json:"description"`
	MimeType       string  `json:"mime_type"`
	FileSize       int64   `json:"file_size"`
	CreatedAt      string  `json:"created_at"`
	RelevanceScore float64 `json:"relevance_score,omitempty"`
	OCRText        *string `json:"ocr_text,omitempty"`
}

// FileDocumentStatsDTO represents file document statistics
type FileDocumentStatsDTO struct {
	TotalDocuments  int            `json:"total_documents"`
	TotalFolders    int            `json:"total_folders"`
	TotalSize       int64          `json:"total_size"`
	DocumentsByType map[string]int `json:"documents_by_type"`
	RecentDocuments []DocumentDTO  `json:"recent_documents"`
	StorageUsage    int64          `json:"storage_usage"`
}

// FolderTreeDTO represents folder tree structure
type FolderTreeDTO struct {
	Folder    FolderDTO       `json:"folder"`
	Children  []FolderTreeDTO `json:"children"`
	Documents []DocumentDTO   `json:"documents"`
}

// OCRTextDTO represents OCR text data
type OCRTextDTO struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	Language   string    `json:"language"`
	Confidence *float64  `json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DocumentDownloadDTO represents document download response
type DocumentDownloadDTO struct {
	Content      []byte    `json:"content"`
	FileName     string    `json:"file_name"`
	MimeType     string    `json:"mime_type"`
	FileSize     int64     `json:"file_size"`
	LastModified time.Time `json:"last_modified"`
}

// CreateDocumentLinkDTO represents the request to create a document link
type CreateDocumentLinkDTO struct {
	DocumentID  string  `json:"document_id" validate:"required"`
	Module      string  `json:"module" validate:"required,oneof=risks assets incidents training compliance audits"`
	EntityID    string  `json:"entity_id" validate:"required"`
	LinkType    string  `json:"link_type" validate:"required,oneof=attachment reference evidence"`
	Description *string `json:"description"`
	LinkedBy    string  `json:"linked_by"` // ID пользователя, который создал связь
}
