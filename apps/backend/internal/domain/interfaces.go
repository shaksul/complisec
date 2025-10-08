package domain

import (
	"context"
	"io"
	"mime/multipart"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"
)

// AssetServiceInterface - интерфейс для AssetService
type AssetServiceInterface interface {
	CreateAsset(ctx context.Context, tenantID string, req dto.CreateAssetRequest, createdBy string) (*repo.Asset, error)
	GetAsset(ctx context.Context, id string) (*repo.Asset, error)
	GetAssetWithDetails(ctx context.Context, id string) (*repo.AssetWithDetails, error)
	ListAssets(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.Asset, error)
	ListAssetsPaginated(ctx context.Context, tenantID string, page, pageSize int, filters map[string]interface{}) ([]repo.Asset, int64, error)
	UpdateAsset(ctx context.Context, id string, req dto.UpdateAssetRequest, updatedBy string) error
	DeleteAsset(ctx context.Context, id string, deletedBy string) error
	AddDocument(ctx context.Context, assetID string, req dto.AssetDocumentRequest, createdBy string) error
	UploadDocument(ctx context.Context, assetID string, req dto.AssetDocumentUploadRequest, createdBy, tenantID string) (*dto.AssetDocumentResponse, error)
	LinkDocument(ctx context.Context, assetID string, req dto.AssetDocumentLinkRequest, createdBy string) (*dto.AssetDocumentResponse, error)
	GetDocumentDownloadPath(ctx context.Context, documentID, userID string) (string, string, error)
	GetDocumentStorage(ctx context.Context, tenantID string, req dto.DocumentStorageRequest) ([]dto.DocumentStorageResponse, int64, error)
	GetAssetDocuments(ctx context.Context, assetID string) ([]repo.AssetDocument, error)
	GetAssetDocumentsFromStorage(ctx context.Context, assetID, tenantID string) ([]dto.DocumentDTO, error)
	DeleteDocument(ctx context.Context, documentID string, deletedBy string) error
	GetDocumentByID(ctx context.Context, documentID string) (*repo.AssetDocument, error)
	AddSoftware(ctx context.Context, assetID string, req dto.AssetSoftwareRequest, addedBy string) error
	GetAssetSoftware(ctx context.Context, assetID string) ([]repo.AssetSoftware, error)
	GetAssetHistory(ctx context.Context, assetID string) ([]repo.AssetHistory, error)
	GetAssetHistoryWithFilters(ctx context.Context, assetID string, filters map[string]interface{}) ([]repo.AssetHistory, error)
	GetAssetRisks(ctx context.Context, assetID string) ([]repo.Risk, error)
	GetAssetIncidents(ctx context.Context, assetID string) ([]repo.Incident, error)
	CanAddRisk(ctx context.Context, assetID string) error
	CanAddIncident(ctx context.Context, assetID string) error
	GetAssetsWithoutOwner(ctx context.Context, tenantID string) ([]repo.Asset, error)
	GetAssetsWithoutPassport(ctx context.Context, tenantID string) ([]repo.Asset, error)
	GetAssetsWithoutCriticality(ctx context.Context, tenantID string) ([]repo.Asset, error)
	BulkUpdateStatus(ctx context.Context, assetIDs []string, newStatus string, updatedBy string) error
	BulkUpdateOwner(ctx context.Context, assetIDs []string, newOwnerID string, updatedBy string) error
	PerformInventory(ctx context.Context, tenantID string, req dto.AssetInventoryRequest, performedBy string) error

	// New centralized document methods
	UploadAssetDocument(ctx context.Context, assetID, tenantID string, file multipart.File, header *multipart.FileHeader, req dto.AssetDocumentUploadRequest, createdBy string) (*dto.DocumentDTO, error)
	LinkExistingDocumentToAsset(ctx context.Context, assetID, documentID, tenantID, linkedBy string) error
	UnlinkDocumentFromAsset(ctx context.Context, assetID, documentID, tenantID, unlinkedBy string) error
	DeleteAssetDocument(ctx context.Context, assetID, documentID, tenantID, deletedBy string) error
	DeleteDocumentLink(ctx context.Context, documentID, tenantID, deletedBy string) error
}

// RiskServiceInterface - интерфейс для RiskService
type RiskServiceInterface interface {
	CreateRisk(ctx context.Context, tenantID, title string, description, category *string, likelihood, impact int, ownerID, assetID *string) (*repo.Risk, error)
	GetRisk(ctx context.Context, id string) (*repo.Risk, error)
	ListRisks(ctx context.Context, tenantID string) ([]repo.Risk, error)
	UpdateRisk(ctx context.Context, id, title string, description, category *string, likelihood, impact int, ownerID, assetID *string) error
	UpdateRiskStatus(ctx context.Context, id, status string) error
	DeleteRisk(ctx context.Context, id string) error
}

// IncidentServiceInterface - интерфейс для IncidentService
type IncidentServiceInterface interface {
	CreateIncident(ctx context.Context, tenantID string, req dto.CreateIncidentRequest, reportedBy string) (*repo.Incident, error)
	GetIncident(ctx context.Context, id, tenantID string) (*repo.Incident, error)
	UpdateIncident(ctx context.Context, id, tenantID string, req dto.UpdateIncidentRequest, updatedBy string) (*repo.Incident, error)
	DeleteIncident(ctx context.Context, id, tenantID string) error
	ListIncidents(ctx context.Context, tenantID string, req dto.IncidentListRequest) ([]*repo.Incident, int, error)
	AddComment(ctx context.Context, incidentID, tenantID string, req dto.IncidentCommentRequest, userID string) (*repo.IncidentComment, error)
	GetComments(ctx context.Context, incidentID, tenantID string) ([]*repo.IncidentComment, error)
	AddAction(ctx context.Context, incidentID, tenantID string, req dto.IncidentActionRequest, createdBy string) (*repo.IncidentAction, error)
	UpdateAction(ctx context.Context, actionID, tenantID string, req dto.IncidentActionRequest, updatedBy string) (*repo.IncidentAction, error)
	GetActions(ctx context.Context, incidentID, tenantID string) ([]*repo.IncidentAction, error)
	GetIncidentMetrics(ctx context.Context, tenantID string) (*repo.IncidentMetricsSummary, error)
	UpdateIncidentStatus(ctx context.Context, id, tenantID string, req dto.IncidentStatusUpdateRequest, updatedBy string) (*repo.Incident, error)
}

// AssetRepoInterface - интерфейс для AssetRepo
type AssetRepoInterface interface {
	Create(ctx context.Context, asset repo.Asset) error
	GetByID(ctx context.Context, id string) (*repo.Asset, error)
	List(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.Asset, error)
	ListPaginated(ctx context.Context, tenantID string, page, pageSize int, filters map[string]interface{}) ([]repo.Asset, int64, error)
	Update(ctx context.Context, asset repo.Asset) error
	SoftDelete(ctx context.Context, id string) error
	GetWithDetails(ctx context.Context, id string) (*repo.AssetWithDetails, error)
	AddDocument(ctx context.Context, assetID, documentType, filePath, createdBy string) error
	GetAssetDocuments(ctx context.Context, assetID string) ([]repo.AssetDocument, error)
	DeleteDocument(ctx context.Context, documentID string) error
	GetDocumentByID(ctx context.Context, documentID string) (*repo.AssetDocument, error)
	AddSoftware(ctx context.Context, assetID, softwareName, version string, installedAt *time.Time) error
	GetAssetSoftware(ctx context.Context, assetID string) ([]repo.AssetSoftware, error)
	AddHistory(ctx context.Context, assetID, fieldChanged, oldValue, newValue, changedBy string) error
	GetAssetHistory(ctx context.Context, assetID string) ([]repo.AssetHistory, error)
	GetAssetHistoryWithFilters(ctx context.Context, assetID string, filters map[string]interface{}) ([]repo.AssetHistory, error)
	GetAssetRisks(ctx context.Context, assetID string) ([]repo.Risk, error)
	GetAssetIncidents(ctx context.Context, assetID string) ([]repo.Incident, error)
	GetAssetsWithoutOwner(ctx context.Context, tenantID string) ([]repo.Asset, error)
	GetAssetsWithoutPassport(ctx context.Context, tenantID string) ([]repo.Asset, error)
	GetAssetsWithoutCriticality(ctx context.Context, tenantID string) ([]repo.Asset, error)

	// Document methods that are missing
	AddDocumentWithFile(ctx context.Context, assetID, documentID, documentType, filePath, fileName, mimeType string, fileSize int64, createdBy string) error
	GetDocumentFromStorage(ctx context.Context, documentID string) (*repo.AssetDocument, error)
	LinkDocumentToAsset(ctx context.Context, assetID, documentID, storageDocumentID, documentType, createdBy string) error
	GetDocumentStorage(ctx context.Context, tenantID string, req dto.DocumentStorageRequest) ([]dto.DocumentStorageResponse, int64, error)
}

// UserRepoInterface - интерфейс для UserRepo
type UserRepoInterface interface {
	GetByID(ctx context.Context, id string) (*repo.User, error)
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
}

// IncidentRepoInterface - интерфейс для IncidentRepo
type IncidentRepoInterface interface {
	Create(ctx context.Context, incident *repo.Incident) error
	GetByID(ctx context.Context, id, tenantID string) (*repo.Incident, error)
	Update(ctx context.Context, incident *repo.Incident) error
	Delete(ctx context.Context, id, tenantID string) error
	List(ctx context.Context, tenantID string, filters map[string]interface{}, limit, offset int) ([]*repo.Incident, int, error)
	AddAsset(ctx context.Context, incidentID, assetID string) error
	RemoveAsset(ctx context.Context, incidentID, assetID string) error
	GetAssets(ctx context.Context, incidentID string) ([]*repo.Asset, error)
	AddRisk(ctx context.Context, incidentID, riskID string) error
	RemoveRisk(ctx context.Context, incidentID, riskID string) error
	GetRisks(ctx context.Context, incidentID string) ([]*repo.Risk, error)
	AddComment(ctx context.Context, comment *repo.IncidentComment) error
	GetComments(ctx context.Context, incidentID string) ([]*repo.IncidentComment, error)
	AddAttachment(ctx context.Context, attachment *repo.IncidentAttachment) error
	GetAttachments(ctx context.Context, incidentID string) ([]*repo.IncidentAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID string) error
	AddAction(ctx context.Context, action *repo.IncidentAction) error
	UpdateAction(ctx context.Context, action *repo.IncidentAction) error
	GetActions(ctx context.Context, incidentID string) ([]*repo.IncidentAction, error)
	DeleteAction(ctx context.Context, actionID string) error
	AddMetric(ctx context.Context, metric *repo.IncidentMetrics) error
	GetMetrics(ctx context.Context, incidentID string) ([]*repo.IncidentMetrics, error)
	GetIncidentMetrics(ctx context.Context, tenantID string) (*repo.IncidentMetricsSummary, error)
}

// RiskRepoInterface - интерфейс для RiskRepo
type RiskRepoInterface interface {
	GetByIDWithTenant(ctx context.Context, id, tenantID string) (*repo.Risk, error)
}

// DocumentServiceInterface - интерфейс для DocumentService
type DocumentServiceInterface interface {
	// Folders
	CreateFolder(ctx context.Context, tenantID string, req dto.CreateFolderDTO, createdBy string) (*dto.FolderDTO, error)
	GetFolder(ctx context.Context, id, tenantID string) (*dto.FolderDTO, error)
	ListFolders(ctx context.Context, tenantID string, parentID *string) ([]dto.FolderDTO, error)
	UpdateFolder(ctx context.Context, id, tenantID string, req dto.UpdateFolderDTO, updatedBy string) error
	DeleteFolder(ctx context.Context, id, tenantID string, deletedBy string) error

	// Documents
	CreateDocument(ctx context.Context, tenantID string, req dto.CreateDocumentDTO, createdBy string) (*dto.DocumentDTO, error)
	UploadDocument(ctx context.Context, tenantID string, file multipart.File, header *multipart.FileHeader, req dto.UploadDocumentDTO, createdBy string) (*dto.DocumentDTO, error)
	GetDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDTO, error)
	GetDocumentsByIDs(ctx context.Context, ids []string, tenantID string) ([]dto.DocumentDTO, error)
	ListDocuments(ctx context.Context, tenantID string, filters dto.FileDocumentFiltersDTO) ([]dto.DocumentDTO, error)
	UpdateDocument(ctx context.Context, id, tenantID string, req dto.UpdateFileDocumentDTO, updatedBy string) error
	DeleteDocument(ctx context.Context, id, tenantID string, deletedBy string) error
	DownloadDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDownloadDTO, error)

	// Search and Stats
	SearchDocuments(ctx context.Context, tenantID, searchTerm string) ([]dto.FileDocumentSearchResultDTO, error)
	GetDocumentStats(ctx context.Context, tenantID string) (*dto.FileDocumentStatsDTO, error)

	// Document Versions
	CreateDocumentVersion(ctx context.Context, documentID, tenantID string, file io.ReadSeeker, header *multipart.FileHeader, createdBy string) (*dto.DocumentVersionDTO, error)
	GetDocumentVersions(ctx context.Context, documentID, tenantID string) ([]dto.DocumentVersionDTO, error)

	// Document Links
	AddDocumentLink(ctx context.Context, documentID string, link dto.CreateDocumentLinkDTO) error
	RemoveDocumentLink(ctx context.Context, documentID, module, entityID string) error
}
