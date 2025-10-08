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

// DocumentStorageServiceInterface - интерфейс для универсального сервиса хранения документов
type DocumentStorageServiceInterface interface {
	// Универсальные методы для работы с документами
	CreateDocument(ctx context.Context, tenantID string, req dto.CreateDocumentDTO, createdBy string) (*dto.DocumentDTO, error)
	UploadDocument(ctx context.Context, tenantID string, file multipart.File, header *multipart.FileHeader, req dto.UploadDocumentDTO, createdBy string) (*dto.DocumentDTO, error)
	GetDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDTO, error)
	DeleteDocument(ctx context.Context, id, tenantID, deletedBy string) error
	ListDocuments(ctx context.Context, tenantID string, filters dto.FileDocumentFiltersDTO) ([]dto.DocumentDTO, error)

	// Методы для связывания документов с модулями
	LinkDocumentToModule(ctx context.Context, documentID, module, entityID, linkType, description string, linkedBy string) error
	UnlinkDocumentFromModule(ctx context.Context, documentID, module, entityID string, unlinkedBy string) error
	GetModuleDocuments(ctx context.Context, module, entityID, tenantID string) ([]dto.DocumentDTO, error)

	// Методы для работы с папками
	CreateFolder(ctx context.Context, tenantID string, req dto.CreateFolderDTO, createdBy string) (*dto.FolderDTO, error)
	GetFolder(ctx context.Context, id, tenantID string) (*dto.FolderDTO, error)
	ListFolders(ctx context.Context, tenantID string, parentID *string) ([]dto.FolderDTO, error)
	UpdateFolder(ctx context.Context, id, tenantID string, req dto.UpdateFolderDTO, updatedBy string) (*dto.FolderDTO, error)
	DeleteFolder(ctx context.Context, id, tenantID string, deletedBy string) error

	// Дополнительные методы для документов
	UpdateDocument(ctx context.Context, id, tenantID string, req dto.UpdateDocumentDTO, updatedBy string) (*dto.DocumentDTO, error)
	DownloadDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDownloadDTO, error)
	SearchDocuments(ctx context.Context, tenantID string, query string) ([]dto.DocumentDTO, error)
	GetDocumentStats(ctx context.Context, tenantID string) (*dto.DocumentStatsDTO, error)

	// Методы для работы с версиями документов
	CreateDocumentVersion(ctx context.Context, documentID, tenantID string, file io.ReadSeeker, header *multipart.FileHeader, createdBy string) (*dto.DocumentVersionDTO, error)
	GetDocumentVersions(ctx context.Context, documentID, tenantID string) ([]dto.DocumentVersionDTO, error)

	// Методы для миграции существующих документов
	MigrateAssetDocument(ctx context.Context, assetDoc *repo.AssetDocument, tenantID, migratedBy string) (*dto.DocumentDTO, error)
	MigrateRiskAttachment(ctx context.Context, riskAttachment *repo.RiskAttachment, tenantID, migratedBy string) (*dto.DocumentDTO, error)
}

// DocumentStorageService - универсальный сервис для работы с документами
type DocumentStorageService struct {
	documentService DocumentServiceInterface
}

// NewDocumentStorageService создает новый экземпляр DocumentStorageService
func NewDocumentStorageService(documentService DocumentServiceInterface) *DocumentStorageService {
	return &DocumentStorageService{
		documentService: documentService,
	}
}

// CreateDocument создает документ без файла
func (s *DocumentStorageService) CreateDocument(ctx context.Context, tenantID string, req dto.CreateDocumentDTO, createdBy string) (*dto.DocumentDTO, error) {
	return s.documentService.CreateDocument(ctx, tenantID, req, createdBy)
}

// UploadDocument загружает документ в централизованное хранилище
func (s *DocumentStorageService) UploadDocument(ctx context.Context, tenantID string, file multipart.File, header *multipart.FileHeader, req dto.UploadDocumentDTO, createdBy string) (*dto.DocumentDTO, error) {
	return s.documentService.UploadDocument(ctx, tenantID, file, header, req, createdBy)
}

// GetDocument получает документ по ID
func (s *DocumentStorageService) GetDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDTO, error) {
	return s.documentService.GetDocument(ctx, id, tenantID)
}

// DeleteDocument удаляет документ
func (s *DocumentStorageService) DeleteDocument(ctx context.Context, id, tenantID, deletedBy string) error {
	return s.documentService.DeleteDocument(ctx, id, tenantID, deletedBy)
}

// ListDocuments получает список документов с фильтрами
func (s *DocumentStorageService) ListDocuments(ctx context.Context, tenantID string, filters dto.FileDocumentFiltersDTO) ([]dto.DocumentDTO, error) {
	return s.documentService.ListDocuments(ctx, tenantID, filters)
}

// LinkDocumentToModule связывает документ с модулем
func (s *DocumentStorageService) LinkDocumentToModule(ctx context.Context, documentID, module, entityID, linkType, description string, linkedBy string) error {
	link := dto.CreateDocumentLinkDTO{
		DocumentID:  documentID,
		Module:      module,
		EntityID:    entityID,
		LinkType:    linkType,
		Description: &description,
	}
	return s.documentService.AddDocumentLink(ctx, documentID, link)
}

// UnlinkDocumentFromModule отвязывает документ от модуля
func (s *DocumentStorageService) UnlinkDocumentFromModule(ctx context.Context, documentID, module, entityID string, unlinkedBy string) error {
	return s.documentService.RemoveDocumentLink(ctx, documentID, module, entityID)
}

// GetModuleDocuments получает документы, связанные с конкретным модулем
func (s *DocumentStorageService) GetModuleDocuments(ctx context.Context, module, entityID, tenantID string) ([]dto.DocumentDTO, error) {
	filters := dto.FileDocumentFiltersDTO{
		Module:   &module,
		EntityID: &entityID,
	}
	return s.documentService.ListDocuments(ctx, tenantID, filters)
}

// CreateFolder создает папку
func (s *DocumentStorageService) CreateFolder(ctx context.Context, tenantID string, req dto.CreateFolderDTO, createdBy string) (*dto.FolderDTO, error) {
	return s.documentService.CreateFolder(ctx, tenantID, req, createdBy)
}

// GetFolder получает папку по ID
func (s *DocumentStorageService) GetFolder(ctx context.Context, id, tenantID string) (*dto.FolderDTO, error) {
	return s.documentService.GetFolder(ctx, id, tenantID)
}

// ListFolders получает список папок
func (s *DocumentStorageService) ListFolders(ctx context.Context, tenantID string, parentID *string) ([]dto.FolderDTO, error) {
	return s.documentService.ListFolders(ctx, tenantID, parentID)
}

// MigrateAssetDocument мигрирует документ из модуля активов в централизованное хранилище
func (s *DocumentStorageService) MigrateAssetDocument(ctx context.Context, assetDoc *repo.AssetDocument, tenantID, migratedBy string) (*dto.DocumentDTO, error) {
	// Проверяем, что файл существует
	if _, err := os.Stat(assetDoc.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("source file not found: %s", assetDoc.FilePath)
	}

	// Читаем файл
	file, err := os.Open(assetDoc.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Создаем заголовок файла для multipart
	header := &multipart.FileHeader{
		Filename: filepath.Base(assetDoc.FilePath),
		Size:     fileInfo.Size(),
	}

	// Определяем MIME тип
	mimeType := assetDoc.Mime
	if mimeType == "" {
		// Пытаемся определить MIME тип по расширению
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

	// Создаем DTO для загрузки
	description := assetDoc.DocumentType
	uploadDTO := dto.UploadDocumentDTO{
		Name:        assetDoc.Title,
		Description: &description, // Используем тип документа как описание
		Tags:        []string{"#активы", "#миграция"},
	}

	// Создаем документ через сервис
	document, err := s.documentService.UploadDocument(ctx, tenantID, file, header, uploadDTO, migratedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create document in storage: %w", err)
	}

	return document, nil
}

// MigrateRiskAttachment мигрирует вложение из модуля рисков в централизованное хранилище
func (s *DocumentStorageService) MigrateRiskAttachment(ctx context.Context, riskAttachment *repo.RiskAttachment, tenantID, migratedBy string) (*dto.DocumentDTO, error) {
	// Проверяем, что файл существует
	if _, err := os.Stat(riskAttachment.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("source file not found: %s", riskAttachment.FilePath)
	}

	// Читаем файл
	file, err := os.Open(riskAttachment.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Создаем заголовок файла для multipart
	header := &multipart.FileHeader{
		Filename: riskAttachment.FileName,
		Size:     fileInfo.Size(),
	}

	// Определяем MIME тип
	mimeType := riskAttachment.MimeType
	if mimeType == "" {
		// Пытаемся определить MIME тип по расширению
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

	// Создаем заголовок с правильным MIME типом
	header.Header = make(map[string][]string)
	header.Header["Content-Type"] = []string{mimeType}

	// Создаем DTO для загрузки
	uploadDTO := dto.UploadDocumentDTO{
		Name:        riskAttachment.FileName,
		Description: riskAttachment.Description,
		Tags:        []string{"#риски", "#миграция"},
	}

	// Создаем документ через сервис
	document, err := s.documentService.UploadDocument(ctx, tenantID, file, header, uploadDTO, migratedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create document in storage: %w", err)
	}

	return document, nil
}

// UpdateFolder обновляет папку
func (s *DocumentStorageService) UpdateFolder(ctx context.Context, id, tenantID string, req dto.UpdateFolderDTO, updatedBy string) (*dto.FolderDTO, error) {
	err := s.documentService.UpdateFolder(ctx, id, tenantID, req, updatedBy)
	if err != nil {
		return nil, err
	}

	// Получаем обновленную папку
	return s.documentService.GetFolder(ctx, id, tenantID)
}

// DeleteFolder удаляет папку
func (s *DocumentStorageService) DeleteFolder(ctx context.Context, id, tenantID string, deletedBy string) error {
	return s.documentService.DeleteFolder(ctx, id, tenantID, deletedBy)
}

// UpdateDocument обновляет документ
func (s *DocumentStorageService) UpdateDocument(ctx context.Context, id, tenantID string, req dto.UpdateDocumentDTO, updatedBy string) (*dto.DocumentDTO, error) {
	// Конвертируем UpdateDocumentDTO в UpdateFileDocumentDTO
	updateReq := dto.UpdateFileDocumentDTO{
		Name:        req.Title,
		Description: req.Description,
		FolderID:    req.FolderID, // Исправлено: передаем FolderID
		Metadata:    req.Metadata, // Исправлено: передаем Metadata
		Tags:        req.Tags,
	}

	err := s.documentService.UpdateDocument(ctx, id, tenantID, updateReq, updatedBy)
	if err != nil {
		return nil, err
	}

	// Получаем обновленный документ
	return s.documentService.GetDocument(ctx, id, tenantID)
}

// DownloadDocument скачивает документ
func (s *DocumentStorageService) DownloadDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDownloadDTO, error) {
	return s.documentService.DownloadDocument(ctx, id, tenantID)
}

// SearchDocuments ищет документы
func (s *DocumentStorageService) SearchDocuments(ctx context.Context, tenantID string, query string) ([]dto.DocumentDTO, error) {
	results, err := s.documentService.SearchDocuments(ctx, tenantID, query)
	if err != nil {
		return nil, err
	}

	// Собираем ID документов для батчевой загрузки
	documentIDs := make([]string, 0, len(results))
	for _, result := range results {
		documentIDs = append(documentIDs, result.DocumentID)
	}

	// Оптимизация: получаем документы батчем вместо N+1 запросов
	documents := make([]dto.DocumentDTO, 0, len(results))
	if len(documentIDs) > 0 {
		// Используем батчевую загрузку для оптимизации производительности
		documents, err = s.documentService.GetDocumentsByIDs(ctx, documentIDs, tenantID)
		if err != nil {
			return nil, fmt.Errorf("failed to get documents by IDs: %w", err)
		}
	}

	return documents, nil
}

// GetDocumentStats получает статистику документов
func (s *DocumentStorageService) GetDocumentStats(ctx context.Context, tenantID string) (*dto.DocumentStatsDTO, error) {
	stats, err := s.documentService.GetDocumentStats(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Конвертируем FileDocumentStatsDTO в DocumentStatsDTO
	return &dto.DocumentStatsDTO{
		TotalDocuments:    stats.TotalDocuments,
		PendingApproval:   0, // TODO: implement approval tracking
		PendingAck:        0, // TODO: implement acknowledgment tracking
		OverdueAck:        0, // TODO: implement overdue tracking
		DocumentsByType:   stats.DocumentsByType,
		DocumentsByStatus: make(map[string]int), // TODO: implement status tracking
	}, nil
}

// CreateDocumentVersion создает новую версию документа
func (s *DocumentStorageService) CreateDocumentVersion(ctx context.Context, documentID, tenantID string, file io.ReadSeeker, header *multipart.FileHeader, createdBy string) (*dto.DocumentVersionDTO, error) {
	return s.documentService.CreateDocumentVersion(ctx, documentID, tenantID, file, header, createdBy)
}

// GetDocumentVersions получает версии документа
func (s *DocumentStorageService) GetDocumentVersions(ctx context.Context, documentID, tenantID string) ([]dto.DocumentVersionDTO, error) {
	return s.documentService.GetDocumentVersions(ctx, documentID, tenantID)
}

// Helper function to create string pointer
func documentStorageStringPtr(s string) *string {
	return &s
}
