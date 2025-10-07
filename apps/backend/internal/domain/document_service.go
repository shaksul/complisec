package domain

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

// DocumentService - сервис для работы с документами
type DocumentService struct {
	documentRepo repo.DocumentRepoInterface
	storagePath  string
}

// NewDocumentService создает новый экземпляр DocumentService
func NewDocumentService(documentRepo repo.DocumentRepoInterface, storagePath string) *DocumentService {
	return &DocumentService{
		documentRepo: documentRepo,
		storagePath:  storagePath,
	}
}

// CreateFolder создает новую папку
func (s *DocumentService) CreateFolder(ctx context.Context, tenantID string, req dto.CreateFolderDTO, createdBy string) (*dto.FolderDTO, error) {
	folderID := uuid.New().String()

	folder := repo.Folder{
		ID:          folderID,
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		OwnerID:     createdBy,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
		Metadata:    req.Metadata,
	}

	if err := s.documentRepo.CreateFolder(ctx, folder); err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	// Создаем аудит запись
	auditLog := repo.DocumentAuditLog{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		FolderID:  &folderID,
		UserID:    createdBy,
		Action:    "created",
		Details:   &req.Name,
		CreatedAt: time.Now(),
	}
	s.documentRepo.CreateDocumentAuditLog(ctx, auditLog)

	return &dto.FolderDTO{
		ID:          folder.ID,
		TenantID:    folder.TenantID,
		Name:        folder.Name,
		Description: folder.Description,
		ParentID:    folder.ParentID,
		OwnerID:     folder.OwnerID,
		CreatedBy:   folder.CreatedBy,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
		IsActive:    folder.IsActive,
		Metadata:    folder.Metadata,
	}, nil
}

// GetFolder получает папку по ID
func (s *DocumentService) GetFolder(ctx context.Context, id, tenantID string) (*dto.FolderDTO, error) {
	folder, err := s.documentRepo.GetFolderByID(ctx, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	return &dto.FolderDTO{
		ID:          folder.ID,
		TenantID:    folder.TenantID,
		Name:        folder.Name,
		Description: folder.Description,
		ParentID:    folder.ParentID,
		OwnerID:     folder.OwnerID,
		CreatedBy:   folder.CreatedBy,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
		IsActive:    folder.IsActive,
		Metadata:    folder.Metadata,
	}, nil
}

// ListFolders получает список папок
func (s *DocumentService) ListFolders(ctx context.Context, tenantID string, parentID *string) ([]dto.FolderDTO, error) {
	folders, err := s.documentRepo.ListFolders(ctx, tenantID, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list folders: %w", err)
	}

	result := make([]dto.FolderDTO, 0)
	for _, folder := range folders {
		result = append(result, dto.FolderDTO{
			ID:          folder.ID,
			TenantID:    folder.TenantID,
			Name:        folder.Name,
			Description: folder.Description,
			ParentID:    folder.ParentID,
			OwnerID:     folder.OwnerID,
			CreatedBy:   folder.CreatedBy,
			CreatedAt:   folder.CreatedAt,
			UpdatedAt:   folder.UpdatedAt,
			IsActive:    folder.IsActive,
			Metadata:    folder.Metadata,
		})
	}

	return result, nil
}

// UpdateFolder обновляет папку
func (s *DocumentService) UpdateFolder(ctx context.Context, id, tenantID string, req dto.UpdateFolderDTO, updatedBy string) error {
	folder, err := s.documentRepo.GetFolderByID(ctx, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get folder: %w", err)
	}

	folder.Name = req.Name
	folder.Description = req.Description
	folder.Metadata = req.Metadata
	folder.UpdatedAt = time.Now()

	if err := s.documentRepo.UpdateFolder(ctx, *folder); err != nil {
		return fmt.Errorf("failed to update folder: %w", err)
	}

	// Создаем аудит запись
	auditLog := repo.DocumentAuditLog{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		FolderID:  &id,
		UserID:    updatedBy,
		Action:    "updated",
		Details:   &req.Name,
		CreatedAt: time.Now(),
	}
	s.documentRepo.CreateDocumentAuditLog(ctx, auditLog)

	return nil
}

// DeleteFolder удаляет папку
func (s *DocumentService) DeleteFolder(ctx context.Context, id, tenantID string, deletedBy string) error {
	if err := s.documentRepo.DeleteFolder(ctx, id, tenantID); err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	// Создаем аудит запись
	auditLog := repo.DocumentAuditLog{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		FolderID:  &id,
		UserID:    deletedBy,
		Action:    "deleted",
		CreatedAt: time.Now(),
	}
	s.documentRepo.CreateDocumentAuditLog(ctx, auditLog)

	return nil
}

// UploadDocument загружает документ
func (s *DocumentService) UploadDocument(ctx context.Context, tenantID string, file multipart.File, header *multipart.FileHeader, req dto.UploadDocumentDTO, createdBy string) (*dto.DocumentDTO, error) {
	// Создаем уникальное имя файла
	fileExt := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	filePath := filepath.Join(s.storagePath, "documents", tenantID, fileName)

	// Создаем директорию если не существует
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Создаем файл
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Копируем содержимое файла
	fileSize, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Вычисляем хеш файла
	hash := sha256.New()
	file.Seek(0, 0) // Возвращаемся к началу файла
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}
	fileHash := fmt.Sprintf("%x", hash.Sum(nil))

	// Создаем документ в БД
	documentID := uuid.New().String()
	document := repo.Document{
		ID:           documentID,
		TenantID:     tenantID,
		Name:         req.Name,
		OriginalName: header.Filename,
		Description:  req.Description,
		FilePath:     filePath,
		FileSize:     fileSize,
		MimeType:     header.Header.Get("Content-Type"),
		FileHash:     fileHash,
		FolderID:     req.FolderID,
		OwnerID:      createdBy,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
		Version:      "1",
		Metadata:     req.Metadata,
	}

	if err := s.documentRepo.CreateDocument(ctx, document); err != nil {
		// Удаляем файл в случае ошибки
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Добавляем теги
	for _, tag := range req.Tags {
		if err := s.documentRepo.AddDocumentTag(ctx, documentID, tag); err != nil {
			// Логируем ошибку, но не прерываем процесс
			fmt.Printf("Failed to add tag %s: %v\n", tag, err)
		}
	}

	// Добавляем связи с другими модулями
	if req.LinkedTo != nil {
		link := repo.DocumentLink{
			ID:         uuid.New().String(),
			DocumentID: documentID,
			Module:     req.LinkedTo.Module,
			EntityID:   req.LinkedTo.EntityID,
			CreatedBy:  createdBy,
			CreatedAt:  time.Now(),
		}
		if err := s.documentRepo.AddDocumentLink(ctx, link); err != nil {
			fmt.Printf("Failed to add document link: %v\n", err)
		}
	}

	// Создаем аудит запись
	auditLog := repo.DocumentAuditLog{
		ID:         uuid.New().String(),
		TenantID:   tenantID,
		DocumentID: &documentID,
		UserID:     createdBy,
		Action:     "uploaded",
		Details:    &req.Name,
		CreatedAt:  time.Now(),
	}
	s.documentRepo.CreateDocumentAuditLog(ctx, auditLog)

	// Получаем теги и связи для ответа
	tags, _ := s.documentRepo.GetDocumentTags(ctx, documentID)
	links, _ := s.documentRepo.GetDocumentLinks(ctx, documentID)

	var documentLinks []dto.DocumentLinkDTO
	for _, link := range links {
		documentLinks = append(documentLinks, dto.DocumentLinkDTO{
			Module:   link.Module,
			EntityID: link.EntityID,
		})
	}

	return &dto.DocumentDTO{
		ID:           document.ID,
		TenantID:     document.TenantID,
		Name:         document.Name,
		OriginalName: document.OriginalName,
		Description:  document.Description,
		FilePath:     document.FilePath,
		FileSize:     document.FileSize,
		MimeType:     document.MimeType,
		FileHash:     document.FileHash,
		FolderID:     document.FolderID,
		OwnerID:      document.OwnerID,
		CreatedBy:    document.CreatedBy,
		CreatedAt:    document.CreatedAt,
		UpdatedAt:    document.UpdatedAt,
		IsActive:     document.IsActive,
		Version:      document.Version,
		Metadata:     document.Metadata,
		Tags:         tags,
		Links:        documentLinks,
	}, nil
}

// GetDocument получает документ по ID
func (s *DocumentService) GetDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDTO, error) {
	document, err := s.documentRepo.GetDocumentByID(ctx, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Получаем теги и связи
	tags, _ := s.documentRepo.GetDocumentTags(ctx, id)
	links, _ := s.documentRepo.GetDocumentLinks(ctx, id)

	var documentLinks []dto.DocumentLinkDTO
	for _, link := range links {
		documentLinks = append(documentLinks, dto.DocumentLinkDTO{
			Module:   link.Module,
			EntityID: link.EntityID,
		})
	}

	// Получаем OCR текст если есть
	var ocrText *string
	ocr, err := s.documentRepo.GetOCRText(ctx, id)
	if err == nil && ocr != nil {
		ocrText = &ocr.Content
	}

	return &dto.DocumentDTO{
		ID:           document.ID,
		TenantID:     document.TenantID,
		Name:         document.Name,
		OriginalName: document.OriginalName,
		Description:  document.Description,
		FilePath:     document.FilePath,
		FileSize:     document.FileSize,
		MimeType:     document.MimeType,
		FileHash:     document.FileHash,
		FolderID:     document.FolderID,
		OwnerID:      document.OwnerID,
		CreatedBy:    document.CreatedBy,
		CreatedAt:    document.CreatedAt,
		UpdatedAt:    document.UpdatedAt,
		IsActive:     document.IsActive,
		Version:      document.Version,
		Metadata:     document.Metadata,
		Tags:         tags,
		Links:        documentLinks,
		OCRText:      ocrText,
	}, nil
}

// ListDocuments получает список документов
func (s *DocumentService) ListDocuments(ctx context.Context, tenantID string, filters dto.FileDocumentFiltersDTO) ([]dto.DocumentDTO, error) {
	filterMap := make(map[string]interface{})

	if filters.FolderID != nil {
		filterMap["folder_id"] = *filters.FolderID
	}
	if filters.MimeType != nil {
		filterMap["mime_type"] = *filters.MimeType
	}
	if filters.OwnerID != nil {
		filterMap["owner_id"] = *filters.OwnerID
	}
	if filters.Search != nil {
		filterMap["search"] = *filters.Search
	}
	if filters.SortBy != nil {
		filterMap["sort_by"] = *filters.SortBy
	}
	if filters.SortOrder != nil {
		filterMap["sort_order"] = *filters.SortOrder
	}
	if filters.Module != nil {
		filterMap["module"] = *filters.Module
	}
	if filters.EntityID != nil {
		filterMap["entity_id"] = *filters.EntityID
	}
	filterMap["page"] = filters.Page
	filterMap["limit"] = filters.Limit

	fmt.Printf("DEBUG: DocumentService.ListDocuments calling repo with tenantID=%s, filters=%v\n", tenantID, filterMap)
	documents, err := s.documentRepo.ListDocuments(ctx, tenantID, filterMap)
	if err != nil {
		fmt.Printf("ERROR: DocumentService.ListDocuments repo error: %v\n", err)
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	fmt.Printf("DEBUG: DocumentService.ListDocuments got %d documents\n", len(documents))

	result := make([]dto.DocumentDTO, 0)
	for _, document := range documents {
		// Получаем теги и связи для каждого документа
		tags, _ := s.documentRepo.GetDocumentTags(ctx, document.ID)
		links, _ := s.documentRepo.GetDocumentLinks(ctx, document.ID)

		var documentLinks []dto.DocumentLinkDTO
		for _, link := range links {
			documentLinks = append(documentLinks, dto.DocumentLinkDTO{
				Module:   link.Module,
				EntityID: link.EntityID,
			})
		}

		result = append(result, dto.DocumentDTO{
			ID:           document.ID,
			TenantID:     document.TenantID,
			Name:         document.Name,
			OriginalName: document.OriginalName,
			Description:  document.Description,
			FilePath:     document.FilePath,
			FileSize:     document.FileSize,
			MimeType:     document.MimeType,
			FileHash:     document.FileHash,
			FolderID:     document.FolderID,
			OwnerID:      document.OwnerID,
			CreatedBy:    document.CreatedBy,
			CreatedAt:    document.CreatedAt,
			UpdatedAt:    document.UpdatedAt,
			IsActive:     document.IsActive,
			Version:      document.Version,
			Metadata:     document.Metadata,
			Tags:         tags,
			Links:        documentLinks,
		})
	}

	return result, nil
}

// UpdateDocument обновляет документ
func (s *DocumentService) UpdateDocument(ctx context.Context, id, tenantID string, req dto.UpdateFileDocumentDTO, updatedBy string) error {
	document, err := s.documentRepo.GetDocumentByID(ctx, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	document.Name = req.Name
	document.Description = req.Description
	document.FolderID = req.FolderID
	document.Metadata = req.Metadata
	document.UpdatedAt = time.Now()

	if err := s.documentRepo.UpdateDocument(ctx, *document); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	// Обновляем теги
	currentTags, _ := s.documentRepo.GetDocumentTags(ctx, id)

	// Удаляем старые теги
	for _, tag := range currentTags {
		if !contains(req.Tags, tag) {
			s.documentRepo.RemoveDocumentTag(ctx, id, tag)
		}
	}

	// Добавляем новые теги
	for _, tag := range req.Tags {
		if !contains(currentTags, tag) {
			s.documentRepo.AddDocumentTag(ctx, id, tag)
		}
	}

	// Создаем аудит запись
	auditLog := repo.DocumentAuditLog{
		ID:         uuid.New().String(),
		TenantID:   tenantID,
		DocumentID: &id,
		UserID:     updatedBy,
		Action:     "updated",
		Details:    &req.Name,
		CreatedAt:  time.Now(),
	}
	s.documentRepo.CreateDocumentAuditLog(ctx, auditLog)

	return nil
}

// DeleteDocument удаляет документ
func (s *DocumentService) DeleteDocument(ctx context.Context, id, tenantID string, deletedBy string) error {
	document, err := s.documentRepo.GetDocumentByID(ctx, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Удаляем файл с диска
	if err := os.Remove(document.FilePath); err != nil {
		fmt.Printf("Failed to remove file %s: %v\n", document.FilePath, err)
	}

	if err := s.documentRepo.DeleteDocument(ctx, id, tenantID); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Создаем аудит запись
	auditLog := repo.DocumentAuditLog{
		ID:         uuid.New().String(),
		TenantID:   tenantID,
		DocumentID: &id,
		UserID:     deletedBy,
		Action:     "deleted",
		Details:    &document.Name,
		CreatedAt:  time.Now(),
	}
	s.documentRepo.CreateDocumentAuditLog(ctx, auditLog)

	return nil
}

// DownloadDocument возвращает путь к файлу для скачивания
func (s *DocumentService) DownloadDocument(ctx context.Context, id, tenantID string) (*dto.DocumentDownloadDTO, error) {
	document, err := s.documentRepo.GetDocumentByID(ctx, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Проверяем существование файла
	if _, err := os.Stat(document.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", document.FilePath)
	}

	// Читаем файл
	file, err := os.Open(document.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &dto.DocumentDownloadDTO{
		Content:      content,
		FileName:     document.OriginalName,
		MimeType:     document.MimeType,
		FileSize:     document.FileSize,
		LastModified: fileInfo.ModTime(),
	}, nil
}

// SearchDocuments выполняет поиск документов
func (s *DocumentService) SearchDocuments(ctx context.Context, tenantID, searchTerm string) ([]dto.FileDocumentSearchResultDTO, error) {
	documents, err := s.documentRepo.SearchDocuments(ctx, tenantID, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	var result []dto.FileDocumentSearchResultDTO
	for _, document := range documents {
		// Получаем OCR текст для поиска
		var ocrText *string
		ocr, err := s.documentRepo.GetOCRText(ctx, document.ID)
		if err == nil && ocr != nil {
			ocrText = &ocr.Content
		}

		result = append(result, dto.FileDocumentSearchResultDTO{
			DocumentID:     document.ID,
			Name:           document.Name,
			Description:    document.Description,
			MimeType:       document.MimeType,
			FileSize:       document.FileSize,
			CreatedAt:      document.CreatedAt.Format(time.RFC3339),
			RelevanceScore: 0.0, // TODO: implement relevance scoring
			OCRText:        ocrText,
		})
	}

	return result, nil
}

// GetDocumentStats получает статистику документов
func (s *DocumentService) GetDocumentStats(ctx context.Context, tenantID string) (*dto.FileDocumentStatsDTO, error) {
	// Получаем общее количество документов
	documents, err := s.documentRepo.ListDocuments(ctx, tenantID, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	// Получаем папки
	folders, err := s.documentRepo.ListFolders(ctx, tenantID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}

	// Подсчитываем статистику
	totalSize := int64(0)
	documentsByType := make(map[string]int)

	for _, doc := range documents {
		totalSize += doc.FileSize
		mimeType := strings.Split(doc.MimeType, "/")[0]
		documentsByType[mimeType]++
	}

	// Получаем последние документы
	recentDocs := make([]dto.DocumentDTO, 0, 5)
	for i, doc := range documents {
		if i >= 5 {
			break
		}
		recentDocs = append(recentDocs, dto.DocumentDTO{
			ID:           doc.ID,
			TenantID:     doc.TenantID,
			Name:         doc.Name,
			OriginalName: doc.OriginalName,
			Description:  doc.Description,
			FilePath:     doc.FilePath,
			FileSize:     doc.FileSize,
			MimeType:     doc.MimeType,
			FileHash:     doc.FileHash,
			FolderID:     doc.FolderID,
			OwnerID:      doc.OwnerID,
			CreatedBy:    doc.CreatedBy,
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			IsActive:     doc.IsActive,
			Version:      doc.Version,
			Metadata:     doc.Metadata,
		})
	}

	return &dto.FileDocumentStatsDTO{
		TotalDocuments:  len(documents),
		TotalFolders:    len(folders),
		TotalSize:       totalSize,
		DocumentsByType: documentsByType,
		RecentDocuments: recentDocs,
		StorageUsage:    totalSize,
	}, nil
}

// AddDocumentLink добавляет связь документа с другим модулем
func (s *DocumentService) AddDocumentLink(ctx context.Context, documentID string, link dto.CreateDocumentLinkDTO) error {
	documentLink := repo.DocumentLink{
		ID:         uuid.New().String(),
		DocumentID: documentID,
		Module:     link.Module,
		EntityID:   link.EntityID,
		CreatedBy:  "system", // TODO: получать из контекста
		CreatedAt:  time.Now(),
	}

	if err := s.documentRepo.AddDocumentLink(ctx, documentLink); err != nil {
		return fmt.Errorf("failed to add document link: %w", err)
	}

	return nil
}

// RemoveDocumentLink удаляет связь документа с другим модулем
func (s *DocumentService) RemoveDocumentLink(ctx context.Context, documentID, module, entityID string) error {
	// Получаем все связи документа
	links, err := s.documentRepo.GetDocumentLinks(ctx, documentID)
	if err != nil {
		return fmt.Errorf("failed to get document links: %w", err)
	}

	// Находим нужную связь
	var linkToRemove *repo.DocumentLink
	for _, link := range links {
		if link.Module == module && link.EntityID == entityID {
			linkToRemove = &link
			break
		}
	}

	if linkToRemove == nil {
		return fmt.Errorf("document link not found")
	}

	// Удаляем связь
	if err := s.documentRepo.DeleteDocumentLink(ctx, documentID, module, entityID); err != nil {
		return fmt.Errorf("failed to delete document link: %w", err)
	}

	return nil
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
