package http

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/gofiber/fiber/v2"
)

// DocumentHandler - обработчик для документов
type DocumentHandler struct {
	documentService domain.DocumentServiceInterface
}

// NewDocumentHandler создает новый экземпляр DocumentHandler
func NewDocumentHandler(documentService domain.DocumentServiceInterface) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
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
	router.Post("/documents/upload", RequirePermission("document.create"), h.UploadDocument)
	router.Get("/documents", RequirePermission("document.read"), h.ListDocuments)
	router.Get("/documents/search", RequirePermission("document.read"), h.SearchDocuments)
	router.Get("/stats", RequirePermission("document.read"), h.GetDocumentStats)
	router.Get("/documents/:id", RequirePermission("document.read"), h.GetDocument)
	router.Put("/documents/:id", RequirePermission("document.edit"), h.UpdateDocument)
	router.Delete("/documents/:id", RequirePermission("document.delete"), h.DeleteDocument)
	router.Get("/documents/:id/download", RequirePermission("document.read"), h.DownloadDocument)
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

	err := h.documentService.UpdateFolder(ctx, folderID, tenantID, req, userID)
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
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	// Получаем файл из формы
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file provided"})
	}

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
	}

	// Парсим связи с другими модулями
	if linkedToStr := c.FormValue("linked_to"); linkedToStr != "" {
		var linkedTo dto.DocumentLinkDTO
		if err := json.Unmarshal([]byte(linkedToStr), &linkedTo); err == nil {
			req.LinkedTo = &linkedTo
		}
	}

	document, err := h.documentService.UploadDocument(ctx, tenantID, src, header, req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to upload document: %v", err)})
	}

	return c.JSON(document)
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

	return c.JSON(document)
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

	return c.JSON(documents)
}

// UpdateDocument обновляет документ
func (h *DocumentHandler) UpdateDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	documentID := c.Params("id")

	var req dto.UpdateFileDocumentDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := h.documentService.UpdateDocument(ctx, documentID, tenantID, req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update document: %v", err)})
	}

	return c.SendStatus(200)
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

	// Устанавливаем заголовки для скачивания
	c.Set("Content-Type", downloadData.MimeType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadData.FileName))
	c.Set("Content-Length", strconv.FormatInt(downloadData.FileSize, 10))
	c.Set("Last-Modified", downloadData.LastModified.Format("Mon, 02 Jan 2006 15:04:05 GMT"))

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

	return c.JSON(results)
}

// GetDocumentStats получает статистику документов
func (h *DocumentHandler) GetDocumentStats(c *fiber.Ctx) error {
	ctx := c.Context()
	tenantID := c.Locals("tenant_id").(string)

	stats, err := h.documentService.GetDocumentStats(ctx, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get document stats: %v", err)})
	}

	return c.JSON(stats)
}

// Helper function to get string pointer
func getStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
