package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/gofiber/fiber/v2"
)

// DocumentHandler - обработчик для документов
type DocumentHandler struct {
	documentService domain.DocumentStorageServiceInterface
	ragService      *domain.RAGService
}

// NewDocumentHandler создает новый экземпляр DocumentHandler
func NewDocumentHandler(documentService domain.DocumentStorageServiceInterface) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

// SetRAGService устанавливает RAG сервис для автоиндексации
func (h *DocumentHandler) SetRAGService(ragService *domain.RAGService) {
	h.ragService = ragService
}

// RegisterRoutes регистрирует маршруты для документов
func (h *DocumentHandler) RegisterRoutes(router fiber.Router) {
	// Тестовый маршрут
	router.Get("/documents/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Documents test route works"})
	})

	// Папки
	router.Post("/folders", RequirePermission("document.create"), h.CreateFolder)
	router.Get("/folders", RequirePermission("document.read"), h.ListFolders)
	router.Get("/folders/:id", RequirePermission("document.read"), h.GetFolder)
	router.Put("/folders/:id", RequirePermission("document.edit"), h.UpdateFolder)
	router.Delete("/folders/:id", RequirePermission("document.delete"), h.DeleteFolder)

	// Документы
	router.Post("/documents", RequirePermission("document.create"), h.CreateDocument)
	router.Post("/documents/upload", RequirePermission("document.create"), func(c *fiber.Ctx) error {
		log.Printf("DEBUG: DocumentHandler.UploadDocument middleware - before handler")
		return h.UploadDocument(c)
	})
	router.Post("/documents/test-upload", h.TestUploadDocument)
	router.Post("/documents/simple-upload", h.SimpleUploadDocument)
	router.Post("/documents/test-multipart", h.TestMultipartForm)
	router.Get("/documents", RequirePermission("document.read"), h.ListDocuments)
	router.Get("/documents/structured", RequirePermission("document.read"), h.ListStructuredDocuments)
	router.Get("/documents/search", RequirePermission("document.read"), h.SearchDocuments)
	router.Get("/stats", RequirePermission("document.read"), h.GetDocumentStats)
	router.Get("/documents/:id", RequirePermission("document.read"), h.GetDocument)
	router.Put("/documents/:id", RequirePermission("document.edit"), h.UpdateDocument)
	router.Delete("/documents/:id", RequirePermission("document.delete"), h.DeleteDocument)
	router.Get("/documents/:id/download", RequirePermission("document.read"), h.DownloadDocument)
	router.Get("/documents/:id/versions", RequirePermission("document.read"), h.GetDocumentVersions)
	router.Get("/documents/versions/:versionId/download", RequirePermission("document.read"), h.DownloadDocumentVersion)
	router.Post("/documents/:id/versions", RequirePermission("document.edit"), h.UploadDocumentVersion)
}

// CreateFolder создает новую папку
func (h *DocumentHandler) CreateFolder(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	var req dto.CreateFolderDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	folder, err := h.documentService.CreateFolder(ctx, tenantID, req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create folder: %v", err)})
	}

	return c.JSON(folder)
}

// GetFolder получает папку по ID
func (h *DocumentHandler) GetFolder(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	folderID := c.Params("id")

	folder, err := h.documentService.GetFolder(ctx, folderID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get folder: %v", err)})
	}

	return c.JSON(folder)
}

// ListFolders получает список папок
func (h *DocumentHandler) ListFolders(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	parentID := c.Query("parent_id")
	var parentIDPtr *string
	if parentID != "" {
		parentIDPtr = &parentID
	}

	folders, err := h.documentService.ListFolders(ctx, tenantID, parentIDPtr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to list folders: %v", err)})
	}

	return c.JSON(folders)
}

// UpdateFolder обновляет папку
func (h *DocumentHandler) UpdateFolder(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	folderID := c.Params("id")

	var req dto.UpdateFolderDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, err := h.documentService.UpdateFolder(ctx, folderID, tenantID, req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update folder: %v", err)})
	}

	return c.SendStatus(200)
}

// DeleteFolder удаляет папку
func (h *DocumentHandler) DeleteFolder(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	folderID := c.Params("id")

	err := h.documentService.DeleteFolder(ctx, folderID, tenantID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete folder: %v", err)})
	}

	return c.SendStatus(200)
}

// UploadDocument загружает документ
func (h *DocumentHandler) UploadDocument(c *fiber.Ctx) error {
	log.Printf("DEBUG: DocumentHandler.UploadDocument START - handler called")

	// Добавляем обработку panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ERROR: DocumentHandler.UploadDocument panic: %v", r)
		}
	}()

	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: DocumentHandler.UploadDocument called with tenantID=%s, userID=%s", tenantID, userID)
	log.Printf("DEBUG: DocumentHandler.UploadDocument request method: %s", c.Method())
	log.Printf("DEBUG: DocumentHandler.UploadDocument request path: %s", c.Path())
	log.Printf("DEBUG: DocumentHandler.UploadDocument content type: %s", c.Get("Content-Type"))
	log.Printf("DEBUG: DocumentHandler.UploadDocument starting multipart form parsing")

	// Получаем файл из формы
	log.Printf("DEBUG: DocumentHandler.UploadDocument getting form file")
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("ERROR: DocumentHandler.UploadDocument FormFile error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "No file provided"})
	}
	log.Printf("DEBUG: DocumentHandler.UploadDocument got file: %s", file.Filename)

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer src.Close()

	// Создаем заголовок файла
	header := &multipart.FileHeader{
		Filename: file.Filename,
		Size:     file.Size,
		Header:   make(map[string][]string),
	}
	header.Header["Content-Type"] = []string{file.Header.Get("Content-Type")}

	// Получаем метаданные
	req := dto.UploadDocumentDTO{
		Name:        c.FormValue("name"),
		Description: getStringPtr(c.FormValue("description")),
		FolderID:    getStringPtr(c.FormValue("folder_id")),
		EnableOCR:   c.FormValue("enable_ocr") == "true",
		Metadata:    getStringPtr(c.FormValue("metadata")),
	}

	// Парсим теги
	if tagsStr := c.FormValue("tags"); tagsStr != "" {
		req.Tags = strings.Split(tagsStr, ",")
		for i, tag := range req.Tags {
			req.Tags[i] = strings.TrimSpace(tag)
		}
	} else {
		// Если теги не указаны, добавляем автоматический тег #documents
		req.Tags = []string{"#documents"}
	}

	// Парсим связи с другими модулями
	if linkedToStr := c.FormValue("linked_to"); linkedToStr != "" {
		var linkedTo dto.DocumentLinkDTO
		if err := json.Unmarshal([]byte(linkedToStr), &linkedTo); err == nil {
			req.LinkedTo = &linkedTo
		}
	}

	log.Printf("DEBUG: DocumentHandler.UploadDocument calling service with req=%+v", req)
	document, err := h.documentService.UploadDocument(ctx, tenantID, src, header, req, userID)
	if err != nil {
		log.Printf("ERROR: DocumentHandler.UploadDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to upload document: %v", err)})
	}

	log.Printf("DEBUG: DocumentHandler.UploadDocument success, documentID=%s", document.ID)

	// Автоматически запускаем индексацию в RAG в фоне
	if h.ragService != nil {
		go func() {
			ragCtx := context.Background()
			if err := h.ragService.IndexDocument(ragCtx, tenantID, document.ID); err != nil {
				log.Printf("WARNING: Auto-indexing to RAG failed for document %s: %v", document.ID, err)
			} else {
				log.Printf("INFO: Document %s automatically indexed to RAG", document.ID)
			}
		}()
	}

	return c.JSON(fiber.Map{"data": document})
}

// TestUploadDocument - простой тест загрузки файла
func (h *DocumentHandler) TestUploadDocument(c *fiber.Ctx) error {
	log.Printf("DEBUG: TestUploadDocument called")

	// Получаем файл из формы
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("ERROR: TestUploadDocument FormFile error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "No file provided"})
	}

	log.Printf("DEBUG: TestUploadDocument got file: %s", file.Filename)

	return c.JSON(fiber.Map{"message": "File received successfully", "filename": file.Filename})
}

// CreateDocument создает документ без файла
func (h *DocumentHandler) CreateDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: DocumentHandler.CreateDocument called")

	var req dto.CreateDocumentDTO
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: DocumentHandler.CreateDocument BodyParser error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("DEBUG: DocumentHandler.CreateDocument req=%+v", req)

	// Создаем документ без файла
	document, err := h.documentService.CreateDocument(ctx, tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: DocumentHandler.CreateDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create document: %v", err)})
	}

	log.Printf("DEBUG: DocumentHandler.CreateDocument success, documentID=%s", document.ID)
	return c.Status(201).JSON(fiber.Map{"data": document})
}

// UploadDocumentVersion загружает новую версию документа
func (h *DocumentHandler) UploadDocumentVersion(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	documentID := c.Params("id")

	log.Printf("DEBUG: DocumentHandler.UploadDocumentVersion called for documentID=%s", documentID)

	// Получаем файл из формы
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("ERROR: DocumentHandler.UploadDocumentVersion FormFile error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "No file provided"})
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		log.Printf("ERROR: DocumentHandler.UploadDocumentVersion file open error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Cannot open file"})
	}
	defer src.Close()

	log.Printf("DEBUG: DocumentHandler.UploadDocumentVersion calling service with documentID=%s", documentID)

	// Создаем версию документа
	version, err := h.documentService.CreateDocumentVersion(ctx, documentID, tenantID, src, file, userID)
	if err != nil {
		log.Printf("ERROR: DocumentHandler.UploadDocumentVersion service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create document version: %v", err)})
	}

	log.Printf("DEBUG: DocumentHandler.UploadDocumentVersion success, versionID=%s", version.ID)
	return c.JSON(fiber.Map{"data": version})
}

// SimpleUploadDocument - максимально простой тест загрузки файла
func (h *DocumentHandler) SimpleUploadDocument(c *fiber.Ctx) error {
	log.Printf("DEBUG: SimpleUploadDocument called")

	return c.JSON(fiber.Map{"message": "Simple endpoint works"})
}

// TestMultipartForm - тест multipart form parsing
func (h *DocumentHandler) TestMultipartForm(c *fiber.Ctx) error {
	log.Printf("DEBUG: TestMultipartForm called")

	// Попробуем получить форму
	_, err := c.MultipartForm()
	if err != nil {
		log.Printf("ERROR: TestMultipartForm MultipartForm error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("MultipartForm error: %v", err)})
	}

	log.Printf("DEBUG: TestMultipartForm got form successfully")

	// Попробуем получить файл
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("ERROR: TestMultipartForm FormFile error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("FormFile error: %v", err)})
	}

	log.Printf("DEBUG: TestMultipartForm got file: %s", file.Filename)

	return c.JSON(fiber.Map{"message": "Multipart form parsing works", "filename": file.Filename})
}

// GetDocument получает документ по ID
func (h *DocumentHandler) GetDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	documentID := c.Params("id")

	document, err := h.documentService.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get document: %v", err)})
	}

	return c.JSON(fiber.Map{"data": document})
}

// ListDocuments получает список документов
func (h *DocumentHandler) ListDocuments(c *fiber.Ctx) error {
	ctx := c.Context()
	fmt.Printf("DEBUG: ListDocuments called, path=%s\n", c.Path())
	tenantIDRaw := c.Locals("tenant_id")
	fmt.Printf("DEBUG: ListDocuments tenantIDRaw=%v\n", tenantIDRaw)
	if tenantIDRaw == nil {
		fmt.Printf("DEBUG: ListDocuments tenant_id is nil\n")
		return c.Status(400).JSON(fiber.Map{"error": "Missing tenant ID"})
	}
	tenantID := tenantIDRaw.(string)
	fmt.Printf("DEBUG: ListDocuments tenantID=%s\n", tenantID)

	// Парсим параметры запроса
	filters := dto.FileDocumentFiltersDTO{
		Page:  1,
		Limit: 20,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filters.Limit = limit
		}
	}

	if folderID := c.Query("folder_id"); folderID != "" {
		filters.FolderID = &folderID
	}

	if mimeType := c.Query("mime_type"); mimeType != "" {
		filters.MimeType = &mimeType
	}

	if ownerID := c.Query("owner_id"); ownerID != "" {
		filters.OwnerID = &ownerID
	}

	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters.SortBy = &sortBy
	}

	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filters.SortOrder = &sortOrder
	}

	fmt.Printf("DEBUG: DocumentHandler.ListDocuments calling service with tenantID=%s, filters=%+v\n", tenantID, filters)
	documents, err := h.documentService.ListDocuments(ctx, tenantID, filters)
	if err != nil {
		fmt.Printf("ERROR: DocumentHandler.ListDocuments service error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to list documents: %v", err)})
	}
	fmt.Printf("DEBUG: DocumentHandler.ListDocuments got %d documents\n", len(documents))

	return c.JSON(fiber.Map{"data": documents})
}

// UpdateDocument обновляет документ
func (h *DocumentHandler) UpdateDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	documentID := c.Params("id")

	log.Printf("DEBUG: DocumentHandler.UpdateDocument called - documentID=%s, tenantID=%s, userID=%s", documentID, tenantID, userID)
	log.Printf("DEBUG: DocumentHandler.UpdateDocument request body: %s", string(c.Body()))

	var req dto.UpdateDocumentDTO
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: DocumentHandler.UpdateDocument BodyParser error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("DEBUG: DocumentHandler.UpdateDocument parsed req=%+v", req)

	document, err := h.documentService.UpdateDocument(ctx, documentID, tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: DocumentHandler.UpdateDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update document: %v", err)})
	}

	log.Printf("DEBUG: DocumentHandler.UpdateDocument success")
	return c.JSON(fiber.Map{"data": document})
}

// DeleteDocument удаляет документ
func (h *DocumentHandler) DeleteDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	documentID := c.Params("id")

	err := h.documentService.DeleteDocument(ctx, documentID, tenantID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete document: %v", err)})
	}

	return c.SendStatus(200)
}

// DownloadDocument скачивает документ
func (h *DocumentHandler) DownloadDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	documentID := c.Params("id")

	downloadData, err := h.documentService.DownloadDocument(ctx, documentID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to download document: %v", err)})
	}

	// Устанавливаем заголовки для скачивания с UTF-8 кодировкой
	contentType := downloadData.MimeType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Для текстовых файлов принудительно устанавливаем UTF-8
	if strings.HasPrefix(contentType, "text/") {
		contentType += "; charset=utf-8"
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadData.FileName))
	c.Set("Content-Length", strconv.FormatInt(downloadData.FileSize, 10))
	c.Set("Last-Modified", downloadData.LastModified.Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	c.Set("Accept-Charset", "utf-8")

	return c.Send(downloadData.Content)
}

// SearchDocuments выполняет поиск документов
func (h *DocumentHandler) SearchDocuments(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	searchTerm := c.Query("q")

	if searchTerm == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing search term"})
	}

	results, err := h.documentService.SearchDocuments(ctx, tenantID, searchTerm)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to search documents: %v", err)})
	}

	return c.JSON(fiber.Map{"data": results})
}

// GetDocumentStats получает статистику документов
func (h *DocumentHandler) GetDocumentStats(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	stats, err := h.documentService.GetDocumentStats(ctx, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get document stats: %v", err)})
	}

	return c.JSON(fiber.Map{"data": stats})
}

// GetDocumentVersions получает версии документа
func (h *DocumentHandler) GetDocumentVersions(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	documentID := c.Params("id")

	log.Printf("DEBUG: DocumentHandler.GetDocumentVersions called for documentID=%s", documentID)

	// Получаем версии документа
	versions, err := h.documentService.GetDocumentVersions(ctx, documentID, tenantID)
	if err != nil {
		log.Printf("ERROR: DocumentHandler.GetDocumentVersions service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get document versions: %v", err)})
	}

	log.Printf("DEBUG: DocumentHandler.GetDocumentVersions returning %d versions", len(versions))
	return c.JSON(fiber.Map{"data": versions})
}

// ListStructuredDocuments получает документы, организованные по структуре папок
func (h *DocumentHandler) ListStructuredDocuments(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	// Получаем ВСЕ документы включая связанные с модулями (для файлового хранилища)
	documents, err := h.documentService.ListAllDocuments(ctx, tenantID, dto.FileDocumentFiltersDTO{
		Page:  1,
		Limit: 1000, // Большое число для получения всех документов
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to list all documents: %v", err)})
	}

	// Организуем документы по структуре папок
	structuredData := organizeDocumentsByStructure(documents)

	return c.JSON(fiber.Map{"data": structuredData})
}

// organizeDocumentsByStructure организует документы по их структуре папок
func organizeDocumentsByStructure(documents []dto.DocumentDTO) map[string]interface{} {
	structure := make(map[string]interface{})

	for _, doc := range documents {
		// Извлекаем модуль и категорию из пути файла
		module, category := extractModuleAndCategoryFromPath(doc.FilePath)

		// Создаем структуру папок
		if structure["modules"] == nil {
			structure["modules"] = make(map[string]interface{})
		}
		modules := structure["modules"].(map[string]interface{})

		if modules[module] == nil {
			modules[module] = make(map[string]interface{})
		}
		moduleMap := modules[module].(map[string]interface{})

		if moduleMap["categories"] == nil {
			moduleMap["categories"] = make(map[string]interface{})
		}
		categories := moduleMap["categories"].(map[string]interface{})

		if categories[category] == nil {
			categories[category] = make([]dto.DocumentDTO, 0)
		}
		documents := categories[category].([]dto.DocumentDTO)
		categories[category] = append(documents, doc)
	}

	return structure
}

// extractModuleAndCategoryFromPath извлекает модуль и категорию из пути файла
func extractModuleAndCategoryFromPath(filePath string) (string, string) {
	// Пример пути: storage/documents/00000000-0000-0000-0000-000000000001/modules/documents/categories/uncategorized/file.txt
	parts := strings.Split(filePath, "/")

	module := "general"
	category := "uncategorized"

	// Ищем "modules" в пути
	for i, part := range parts {
		if part == "modules" && i+1 < len(parts) {
			module = parts[i+1]
		}
		if part == "categories" && i+1 < len(parts) {
			category = parts[i+1]
		}
	}

	return module, category
}

// Helper function to get string pointer
func getStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// DownloadDocumentVersion скачивает версию документа
func (h *DocumentHandler) DownloadDocumentVersion(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	versionID := c.Params("versionId")

	downloadData, err := h.documentService.DownloadDocumentVersion(ctx, versionID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to download document version: %v", err)})
	}

	// Устанавливаем заголовки для скачивания с UTF-8 кодировкой
	contentType := downloadData.MimeType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Для текстовых файлов принудительно устанавливаем UTF-8
	if strings.HasPrefix(contentType, "text/") {
		contentType += "; charset=utf-8"
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadData.FileName))
	c.Set("Content-Length", strconv.FormatInt(downloadData.FileSize, 10))
	c.Set("Last-Modified", downloadData.LastModified.Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	c.Set("Accept-Charset", "utf-8")

	return c.Send(downloadData.Content)
}
