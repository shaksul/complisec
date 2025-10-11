package domain

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"
)

// DocumentStorageServiceInterface - РёРЅС‚РµСЂС„РµР№СЃ РґР»СЏ СѓРЅРёРІРµСЂСЃР°Р»СЊРЅРѕРіРѕ СЃРµСЂРІРёСЃР° С…СЂР°РЅРµРЅРёСЏ РґРѕРєСѓРјРµРЅС‚РѕРІ
type DocumentStorageServiceInterface interface {
	// РЈРЅРёРІРµСЂСЃР°Р»СЊРЅС‹Рµ РјРµС‚РѕРґС‹ РґР»СЏ СЂР°Р±РѕС‚С‹ СЃ РґРѕРєСѓРјРµРЅС‚Р°РјРё
	CreateDocument(ctx context.Context, tenantID string, req dto.CreateDocumentDTO, createdBy string) (*dto.DocumentDTO, error)
	UploadDocument(ctx context.Context, tenantID string, file multipart.File, header *multipart.FileHeader, req dto.UploadDocumentDTO, createdBy string) (*dto.DocumentDTO, error)
	SaveGeneratedDocument(ctx context.Context, tenantID string, content []byte, fileName, mimeType string, req dto.UploadDocumentDTO, createdBy string) (*dto.DocumentDTO, error)
	GetDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDTO, error)
	DeleteDocument(ctx context.Context, id, tenantID, deletedBy string) error
	ListDocuments(ctx context.Context, tenantID string, filters dto.FileDocumentFiltersDTO) ([]dto.DocumentDTO, error)
	ListAllDocuments(ctx context.Context, tenantID string, filters dto.FileDocumentFiltersDTO) ([]dto.DocumentDTO, error)

	// РњРµС‚РѕРґС‹ РґР»СЏ СЃРІСЏР·С‹РІР°РЅРёСЏ РґРѕРєСѓРјРµРЅС‚РѕРІ СЃ РјРѕРґСѓР»СЏРјРё
	LinkDocumentToModule(ctx context.Context, documentID, module, entityID, linkType, description string, linkedBy string) error
	UnlinkDocumentFromModule(ctx context.Context, documentID, module, entityID string, unlinkedBy string) error
	GetModuleDocuments(ctx context.Context, module, entityID, tenantID string) ([]dto.DocumentDTO, error)

	// РњРµС‚РѕРґС‹ РґР»СЏ СЂР°Р±РѕС‚С‹ СЃ РїР°РїРєР°РјРё
	CreateFolder(ctx context.Context, tenantID string, req dto.CreateFolderDTO, createdBy string) (*dto.FolderDTO, error)
	GetFolder(ctx context.Context, id, tenantID string) (*dto.FolderDTO, error)
	ListFolders(ctx context.Context, tenantID string, parentID *string) ([]dto.FolderDTO, error)
	UpdateFolder(ctx context.Context, id, tenantID string, req dto.UpdateFolderDTO, updatedBy string) (*dto.FolderDTO, error)
	DeleteFolder(ctx context.Context, id, tenantID string, deletedBy string) error

	// Р”РѕРїРѕР»РЅРёС‚РµР»СЊРЅС‹Рµ РјРµС‚РѕРґС‹ РґР»СЏ РґРѕРєСѓРјРµРЅС‚РѕРІ
	UpdateDocument(ctx context.Context, id, tenantID string, req dto.UpdateDocumentDTO, updatedBy string) (*dto.DocumentDTO, error)
	DownloadDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDownloadDTO, error)
	SearchDocuments(ctx context.Context, tenantID string, query string) ([]dto.DocumentDTO, error)
	GetDocumentStats(ctx context.Context, tenantID string) (*dto.DocumentStatsDTO, error)

	// РњРµС‚РѕРґС‹ РґР»СЏ СЂР°Р±РѕС‚С‹ СЃ РІРµСЂСЃРёСЏРјРё РґРѕРєСѓРјРµРЅС‚РѕРІ
	CreateDocumentVersion(ctx context.Context, documentID, tenantID string, file io.ReadSeeker, header *multipart.FileHeader, createdBy string) (*dto.DocumentVersionDTO, error)
	GetDocumentVersions(ctx context.Context, documentID, tenantID string) ([]dto.DocumentVersionDTO, error)
	DownloadDocumentVersion(ctx context.Context, versionID, tenantID string) (*dto.DocumentDownloadDTO, error)

	// РњРµС‚РѕРґС‹ РґР»СЏ РјРёРіСЂР°С†РёРё СЃСѓС‰РµСЃС‚РІСѓСЋС‰РёС… РґРѕРєСѓРјРµРЅС‚РѕРІ
	MigrateAssetDocument(ctx context.Context, assetDoc *repo.AssetDocument, tenantID, migratedBy string) (*dto.DocumentDTO, error)
	MigrateRiskAttachment(ctx context.Context, riskAttachment *repo.RiskAttachment, tenantID, migratedBy string) (*dto.DocumentDTO, error)
}

// DocumentStorageService - СѓРЅРёРІРµСЂСЃР°Р»СЊРЅС‹Р№ СЃРµСЂРІРёСЃ РґР»СЏ СЂР°Р±РѕС‚С‹ СЃ РґРѕРєСѓРјРµРЅС‚Р°РјРё
type DocumentStorageService struct {
	documentService DocumentServiceInterface
}

// NewDocumentStorageService СЃРѕР·РґР°РµС‚ РЅРѕРІС‹Р№ СЌРєР·РµРјРїР»СЏСЂ DocumentStorageService
func NewDocumentStorageService(documentService DocumentServiceInterface) *DocumentStorageService {
	return &DocumentStorageService{
		documentService: documentService,
	}
}

// CreateDocument СЃРѕР·РґР°РµС‚ РґРѕРєСѓРјРµРЅС‚ Р±РµР· С„Р°Р№Р»Р°
func (s *DocumentStorageService) CreateDocument(ctx context.Context, tenantID string, req dto.CreateDocumentDTO, createdBy string) (*dto.DocumentDTO, error) {
	return s.documentService.CreateDocument(ctx, tenantID, req, createdBy)
}

// UploadDocument Р·Р°РіСЂСѓР¶Р°РµС‚ РґРѕРєСѓРјРµРЅС‚ РІ С†РµРЅС‚СЂР°Р»РёР·РѕРІР°РЅРЅРѕРµ С…СЂР°РЅРёР»РёС‰Рµ
func (s *DocumentStorageService) UploadDocument(ctx context.Context, tenantID string, file multipart.File, header *multipart.FileHeader, req dto.UploadDocumentDTO, createdBy string) (*dto.DocumentDTO, error) {
	return s.documentService.UploadDocument(ctx, tenantID, file, header, req, createdBy)
}

// SaveGeneratedDocument �?�?���?����' ����'��?�?��?�?�?���?�?�<�� ����'�������? ��� �?�?��?�?��?�'
func (s *DocumentStorageService) SaveGeneratedDocument(ctx context.Context, tenantID string, content []byte, fileName, mimeType string, req dto.UploadDocumentDTO, createdBy string) (*dto.DocumentDTO, error) {
	return s.documentService.SaveGeneratedDocument(ctx, tenantID, content, fileName, mimeType, req, createdBy)
}

// GetDocument РїРѕР»СѓС‡Р°РµС‚ РґРѕРєСѓРјРµРЅС‚ РїРѕ ID
func (s *DocumentStorageService) GetDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDTO, error) {
	return s.documentService.GetDocument(ctx, id, tenantID)
}

// DeleteDocument СѓРґР°Р»СЏРµС‚ РґРѕРєСѓРјРµРЅС‚
func (s *DocumentStorageService) DeleteDocument(ctx context.Context, id, tenantID, deletedBy string) error {
	return s.documentService.DeleteDocument(ctx, id, tenantID, deletedBy)
}

// ListDocuments РїРѕР»СѓС‡Р°РµС‚ СЃРїРёСЃРѕРє РґРѕРєСѓРјРµРЅС‚РѕРІ СЃ С„РёР»СЊС‚СЂР°РјРё
func (s *DocumentStorageService) ListDocuments(ctx context.Context, tenantID string, filters dto.FileDocumentFiltersDTO) ([]dto.DocumentDTO, error) {
	return s.documentService.ListDocuments(ctx, tenantID, filters)
}

// LinkDocumentToModule СЃРІСЏР·С‹РІР°РµС‚ РґРѕРєСѓРјРµРЅС‚ СЃ РјРѕРґСѓР»РµРј
func (s *DocumentStorageService) LinkDocumentToModule(ctx context.Context, documentID, module, entityID, linkType, description string, linkedBy string) error {
	link := dto.CreateDocumentLinkDTO{
		DocumentID:  documentID,
		Module:      module,
		EntityID:    entityID,
		LinkType:    linkType,
		Description: &description,
		LinkedBy:    linkedBy, // РџРµСЂРµРґР°РµРј ID РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ
	}
	return s.documentService.AddDocumentLink(ctx, documentID, link)
}

// UnlinkDocumentFromModule РѕС‚РІСЏР·С‹РІР°РµС‚ РґРѕРєСѓРјРµРЅС‚ РѕС‚ РјРѕРґСѓР»СЏ
func (s *DocumentStorageService) UnlinkDocumentFromModule(ctx context.Context, documentID, module, entityID string, unlinkedBy string) error {
	return s.documentService.RemoveDocumentLink(ctx, documentID, module, entityID)
}

// GetModuleDocuments РїРѕР»СѓС‡Р°РµС‚ РґРѕРєСѓРјРµРЅС‚С‹, СЃРІСЏР·Р°РЅРЅС‹Рµ СЃ РєРѕРЅРєСЂРµС‚РЅС‹Рј РјРѕРґСѓР»РµРј
func (s *DocumentStorageService) GetModuleDocuments(ctx context.Context, module, entityID, tenantID string) ([]dto.DocumentDTO, error) {
	fmt.Printf("DEBUG: GetModuleDocuments called with module=%s entityID=%s tenantID=%s\n", module, entityID, tenantID)
	module = normalizeModuleName(module)
	filters := dto.FileDocumentFiltersDTO{
		Module:   &module,
		EntityID: &entityID,
	}
	documents, err := s.documentService.ListDocuments(ctx, tenantID, filters)
	fmt.Printf("DEBUG: GetModuleDocuments returned %d documents (err=%v)\n", len(documents), err)
	return documents, err
}

// CreateFolder СЃРѕР·РґР°РµС‚ РїР°РїРєСѓ
func (s *DocumentStorageService) CreateFolder(ctx context.Context, tenantID string, req dto.CreateFolderDTO, createdBy string) (*dto.FolderDTO, error) {
	return s.documentService.CreateFolder(ctx, tenantID, req, createdBy)
}

// GetFolder РїРѕР»СѓС‡Р°РµС‚ РїР°РїРєСѓ РїРѕ ID
func (s *DocumentStorageService) GetFolder(ctx context.Context, id, tenantID string) (*dto.FolderDTO, error) {
	return s.documentService.GetFolder(ctx, id, tenantID)
}

// ListFolders РїРѕР»СѓС‡Р°РµС‚ СЃРїРёСЃРѕРє РїР°РїРѕРє
func (s *DocumentStorageService) ListFolders(ctx context.Context, tenantID string, parentID *string) ([]dto.FolderDTO, error) {
	return s.documentService.ListFolders(ctx, tenantID, parentID)
}

// MigrateAssetDocument РјРёРіСЂРёСЂСѓРµС‚ РґРѕРєСѓРјРµРЅС‚ РёР· РјРѕРґСѓР»СЏ Р°РєС‚РёРІРѕРІ РІ С†РµРЅС‚СЂР°Р»РёР·РѕРІР°РЅРЅРѕРµ С…СЂР°РЅРёР»РёС‰Рµ
func (s *DocumentStorageService) MigrateAssetDocument(ctx context.Context, assetDoc *repo.AssetDocument, tenantID, migratedBy string) (*dto.DocumentDTO, error) {
	// РџСЂРѕРІРµСЂСЏРµРј, С‡С‚Рѕ С„Р°Р№Р» СЃСѓС‰РµСЃС‚РІСѓРµС‚
	if _, err := os.Stat(assetDoc.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("source file not found: %s", assetDoc.FilePath)
	}

	// Р§РёС‚Р°РµРј С„Р°Р№Р»
	file, err := os.Open(assetDoc.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// РџРѕР»СѓС‡Р°РµРј РёРЅС„РѕСЂРјР°С†РёСЋ Рѕ С„Р°Р№Р»Рµ
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// РЎРѕР·РґР°РµРј Р·Р°РіРѕР»РѕРІРѕРє С„Р°Р№Р»Р° РґР»СЏ multipart
	header := &multipart.FileHeader{
		Filename: filepath.Base(assetDoc.FilePath),
		Size:     fileInfo.Size(),
	}

	// РћРїСЂРµРґРµР»СЏРµРј MIME С‚РёРї
	mimeType := assetDoc.Mime
	if mimeType == "" {
		// РџС‹С‚Р°РµРјСЃСЏ РѕРїСЂРµРґРµР»РёС‚СЊ MIME С‚РёРї РїРѕ СЂР°СЃС€РёСЂРµРЅРёСЋ
		ext := strings.ToLower(filepath.Ext(assetDoc.FilePath))
		switch ext {
		case ".pdf":
			mimeType = "application/pdf"
		case ".doc":
			mimeType = "application/msword"
		case ".docx":
			mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		case ".txt":
			mimeType = "text/plain"
		default:
			mimeType = "application/octet-stream"
		}
	}

	// РЎРѕР·РґР°РµРј DTO РґР»СЏ Р·Р°РіСЂСѓР·РєРё
	description := assetDoc.DocumentType
	uploadDTO := dto.UploadDocumentDTO{
		Name:        assetDoc.Title,
		Description: &description, // РСЃРїРѕР»СЊР·СѓРµРј С‚РёРї РґРѕРєСѓРјРµРЅС‚Р° РєР°Рє РѕРїРёСЃР°РЅРёРµ
		Tags:        []string{"#Р°РєС‚РёРІС‹", "#РјРёРіСЂР°С†РёСЏ"},
	}

	// РЎРѕР·РґР°РµРј РґРѕРєСѓРјРµРЅС‚ С‡РµСЂРµР· СЃРµСЂРІРёСЃ
	document, err := s.documentService.UploadDocument(ctx, tenantID, file, header, uploadDTO, migratedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create document in storage: %w", err)
	}

	return document, nil
}

// MigrateRiskAttachment РјРёРіСЂРёСЂСѓРµС‚ РІР»РѕР¶РµРЅРёРµ РёР· РјРѕРґСѓР»СЏ СЂРёСЃРєРѕРІ РІ С†РµРЅС‚СЂР°Р»РёР·РѕРІР°РЅРЅРѕРµ С…СЂР°РЅРёР»РёС‰Рµ
func (s *DocumentStorageService) MigrateRiskAttachment(ctx context.Context, riskAttachment *repo.RiskAttachment, tenantID, migratedBy string) (*dto.DocumentDTO, error) {
	// РџСЂРѕРІРµСЂСЏРµРј, С‡С‚Рѕ С„Р°Р№Р» СЃСѓС‰РµСЃС‚РІСѓРµС‚
	if _, err := os.Stat(riskAttachment.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("source file not found: %s", riskAttachment.FilePath)
	}

	// Р§РёС‚Р°РµРј С„Р°Р№Р»
	file, err := os.Open(riskAttachment.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// РџРѕР»СѓС‡Р°РµРј РёРЅС„РѕСЂРјР°С†РёСЋ Рѕ С„Р°Р№Р»Рµ
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// РЎРѕР·РґР°РµРј Р·Р°РіРѕР»РѕРІРѕРє С„Р°Р№Р»Р° РґР»СЏ multipart
	header := &multipart.FileHeader{
		Filename: riskAttachment.FileName,
		Size:     fileInfo.Size(),
	}

	// РћРїСЂРµРґРµР»СЏРµРј MIME С‚РёРї
	mimeType := riskAttachment.MimeType
	if mimeType == "" {
		// РџС‹С‚Р°РµРјСЃСЏ РѕРїСЂРµРґРµР»РёС‚СЊ MIME С‚РёРї РїРѕ СЂР°СЃС€РёСЂРµРЅРёСЋ
		ext := strings.ToLower(filepath.Ext(riskAttachment.FileName))
		switch ext {
		case ".pdf":
			mimeType = "application/pdf"
		case ".doc":
			mimeType = "application/msword"
		case ".docx":
			mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		case ".txt":
			mimeType = "text/plain"
		default:
			mimeType = "application/octet-stream"
		}
	}

	// РЎРѕР·РґР°РµРј Р·Р°РіРѕР»РѕРІРѕРє СЃ РїСЂР°РІРёР»СЊРЅС‹Рј MIME С‚РёРїРѕРј
	header.Header = make(map[string][]string)
	header.Header["Content-Type"] = []string{mimeType}

	// РЎРѕР·РґР°РµРј DTO РґР»СЏ Р·Р°РіСЂСѓР·РєРё
	uploadDTO := dto.UploadDocumentDTO{
		Name:        riskAttachment.FileName,
		Description: riskAttachment.Description,
		Tags:        []string{"#СЂРёСЃРєРё", "#РјРёРіСЂР°С†РёСЏ"},
	}

	// РЎРѕР·РґР°РµРј РґРѕРєСѓРјРµРЅС‚ С‡РµСЂРµР· СЃРµСЂРІРёСЃ
	document, err := s.documentService.UploadDocument(ctx, tenantID, file, header, uploadDTO, migratedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create document in storage: %w", err)
	}

	return document, nil
}

// UpdateFolder РѕР±РЅРѕРІР»СЏРµС‚ РїР°РїРєСѓ
func (s *DocumentStorageService) UpdateFolder(ctx context.Context, id, tenantID string, req dto.UpdateFolderDTO, updatedBy string) (*dto.FolderDTO, error) {
	err := s.documentService.UpdateFolder(ctx, id, tenantID, req, updatedBy)
	if err != nil {
		return nil, err
	}

	// РџРѕР»СѓС‡Р°РµРј РѕР±РЅРѕРІР»РµРЅРЅСѓСЋ РїР°РїРєСѓ
	return s.documentService.GetFolder(ctx, id, tenantID)
}

// DeleteFolder СѓРґР°Р»СЏРµС‚ РїР°РїРєСѓ
func (s *DocumentStorageService) DeleteFolder(ctx context.Context, id, tenantID string, deletedBy string) error {
	return s.documentService.DeleteFolder(ctx, id, tenantID, deletedBy)
}

// UpdateDocument РѕР±РЅРѕРІР»СЏРµС‚ РґРѕРєСѓРјРµРЅС‚
func (s *DocumentStorageService) UpdateDocument(ctx context.Context, id, tenantID string, req dto.UpdateDocumentDTO, updatedBy string) (*dto.DocumentDTO, error) {
	// РљРѕРЅРІРµСЂС‚РёСЂСѓРµРј UpdateDocumentDTO РІ UpdateFileDocumentDTO
	updateReq := dto.UpdateFileDocumentDTO{
		Name:        req.Title,
		Description: req.Description,
		FolderID:    req.FolderID, // РСЃРїСЂР°РІР»РµРЅРѕ: РїРµСЂРµРґР°РµРј FolderID
		Metadata:    req.Metadata, // РСЃРїСЂР°РІР»РµРЅРѕ: РїРµСЂРµРґР°РµРј Metadata
		Tags:        req.Tags,
	}

	err := s.documentService.UpdateDocument(ctx, id, tenantID, updateReq, updatedBy)
	if err != nil {
		return nil, err
	}

	// РџРѕР»СѓС‡Р°РµРј РѕР±РЅРѕРІР»РµРЅРЅС‹Р№ РґРѕРєСѓРјРµРЅС‚
	return s.documentService.GetDocument(ctx, id, tenantID)
}

// DownloadDocument СЃРєР°С‡РёРІР°РµС‚ РґРѕРєСѓРјРµРЅС‚
func (s *DocumentStorageService) DownloadDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDownloadDTO, error) {
	return s.documentService.DownloadDocument(ctx, id, tenantID)
}

// SearchDocuments РёС‰РµС‚ РґРѕРєСѓРјРµРЅС‚С‹
func (s *DocumentStorageService) SearchDocuments(ctx context.Context, tenantID string, query string) ([]dto.DocumentDTO, error) {
	results, err := s.documentService.SearchDocuments(ctx, tenantID, query)
	if err != nil {
		return nil, err
	}

	// РЎРѕР±РёСЂР°РµРј ID РґРѕРєСѓРјРµРЅС‚РѕРІ РґР»СЏ Р±Р°С‚С‡РµРІРѕР№ Р·Р°РіСЂСѓР·РєРё
	documentIDs := make([]string, 0, len(results))
	for _, result := range results {
		documentIDs = append(documentIDs, result.DocumentID)
	}

	// РћРїС‚РёРјРёР·Р°С†РёСЏ: РїРѕР»СѓС‡Р°РµРј РґРѕРєСѓРјРµРЅС‚С‹ Р±Р°С‚С‡РµРј РІРјРµСЃС‚Рѕ N+1 Р·Р°РїСЂРѕСЃРѕРІ
	documents := make([]dto.DocumentDTO, 0, len(results))
	if len(documentIDs) > 0 {
		// РСЃРїРѕР»СЊР·СѓРµРј Р±Р°С‚С‡РµРІСѓСЋ Р·Р°РіСЂСѓР·РєСѓ РґР»СЏ РѕРїС‚РёРјРёР·Р°С†РёРё РїСЂРѕРёР·РІРѕРґРёС‚РµР»СЊРЅРѕСЃС‚Рё
		documents, err = s.documentService.GetDocumentsByIDs(ctx, documentIDs, tenantID)
		if err != nil {
			return nil, fmt.Errorf("failed to get documents by IDs: %w", err)
		}
	}

	return documents, nil
}

// GetDocumentStats РїРѕР»СѓС‡Р°РµС‚ СЃС‚Р°С‚РёСЃС‚РёРєСѓ РґРѕРєСѓРјРµРЅС‚РѕРІ
func (s *DocumentStorageService) GetDocumentStats(ctx context.Context, tenantID string) (*dto.DocumentStatsDTO, error) {
	stats, err := s.documentService.GetDocumentStats(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// РљРѕРЅРІРµСЂС‚РёСЂСѓРµРј FileDocumentStatsDTO РІ DocumentStatsDTO
	return &dto.DocumentStatsDTO{
		TotalDocuments:    stats.TotalDocuments,
		PendingApproval:   0, // TODO: implement approval tracking
		PendingAck:        0, // TODO: implement acknowledgment tracking
		OverdueAck:        0, // TODO: implement overdue tracking
		DocumentsByType:   stats.DocumentsByType,
		DocumentsByStatus: make(map[string]int), // TODO: implement status tracking
	}, nil
}

// CreateDocumentVersion СЃРѕР·РґР°РµС‚ РЅРѕРІСѓСЋ РІРµСЂСЃРёСЋ РґРѕРєСѓРјРµРЅС‚Р°
func (s *DocumentStorageService) CreateDocumentVersion(ctx context.Context, documentID, tenantID string, file io.ReadSeeker, header *multipart.FileHeader, createdBy string) (*dto.DocumentVersionDTO, error) {
	return s.documentService.CreateDocumentVersion(ctx, documentID, tenantID, file, header, createdBy)
}

// GetDocumentVersions РїРѕР»СѓС‡Р°РµС‚ РІРµСЂСЃРёРё РґРѕРєСѓРјРµРЅС‚Р°
func (s *DocumentStorageService) GetDocumentVersions(ctx context.Context, documentID, tenantID string) ([]dto.DocumentVersionDTO, error) {
	return s.documentService.GetDocumentVersions(ctx, documentID, tenantID)
}

// DownloadDocumentVersion РїРѕР»СѓС‡Р°РµС‚ РІРµСЂСЃРёСЋ РґРѕРєСѓРјРµРЅС‚Р° РґР»СЏ СЃРєР°С‡РёРІР°РЅРёСЏ
func (s *DocumentStorageService) DownloadDocumentVersion(ctx context.Context, versionID, tenantID string) (*dto.DocumentDownloadDTO, error) {
	return s.documentService.DownloadDocumentVersion(ctx, versionID, tenantID)
}

// Helper function to create string pointer
func documentStorageStringPtr(s string) *string {
	return &s
}
