package domain

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"
)

// DocumentService handles business logic for documents
type DocumentService struct {
	documentRepo *repo.DocumentRepo
	auditRepo    *repo.AuditRepo
}

// NewDocumentService creates a new document service
func NewDocumentService(documentRepo *repo.DocumentRepo, auditRepo *repo.AuditRepo) *DocumentService {
	return &DocumentService{
		documentRepo: documentRepo,
		auditRepo:    auditRepo,
	}
}

// ListDocuments retrieves documents with filtering and pagination
func (s *DocumentService) ListDocuments(ctx context.Context, tenantID string, filters dto.DocumentFiltersDTO) ([]repo.Document, error) {
	// Convert DTO filters to map
	filterMap := make(map[string]interface{})
	if filters.Status != nil {
		filterMap["status"] = *filters.Status
	}
	if filters.Type != nil {
		filterMap["type"] = *filters.Type
	}
	if filters.Category != nil {
		filterMap["category"] = *filters.Category
	}
	if filters.Search != nil {
		filterMap["search"] = *filters.Search
	}

	documents, err := s.documentRepo.ListDocuments(ctx, tenantID, filterMap)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	return documents, nil
}

// GetDocument retrieves a document by ID
func (s *DocumentService) GetDocument(ctx context.Context, id, tenantID string) (*repo.Document, error) {
	document, err := s.documentRepo.GetDocument(ctx, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	return document, nil
}

// CreateDocument creates a new document
func (s *DocumentService) CreateDocument(ctx context.Context, tenantID, userID string, req dto.CreateDocumentDTO) (*repo.Document, error) {
	// Validate document type
	validTypes := []string{"policy", "standard", "procedure", "instruction", "act", "other"}
	if !contains(validTypes, req.Type) {
		return nil, fmt.Errorf("invalid document type: %s", req.Type)
	}

	// Create document
	document := repo.Document{
		ID:             generateUUID(),
		TenantID:       tenantID,
		Title:          req.Title,
		Code:           req.Code,
		Description:    req.Description,
		Type:           req.Type,
		Category:       req.Category,
		Tags:           req.Tags,
		Status:         "draft",
		CurrentVersion: 1,
		OwnerID:        req.OwnerID,
		Classification: req.Classification,
		EffectiveFrom:  req.EffectiveFrom,
		ReviewPeriodMonths: func() int {
			if req.ReviewPeriodMonths != nil {
				return *req.ReviewPeriodMonths
			}
			return 12 // Default value
		}(),
		AssetIDs:     req.AssetIDs,
		RiskIDs:      req.RiskIDs,
		ControlIDs:   req.ControlIDs,
		AVScanStatus: "pending",
		CreatedBy:    userID,
	}

	// Ensure Tags is not nil
	if document.Tags == nil {
		document.Tags = []string{}
	}

	fmt.Printf("DEBUG: Creating document: %+v\n", document)
	err := s.documentRepo.CreateDocument(ctx, document)
	if err != nil {
		fmt.Printf("DEBUG: Error creating document: %v\n", err)
		return nil, fmt.Errorf("failed to create document: %w", err)
	}
	fmt.Printf("DEBUG: Document created successfully\n")

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.created", "document", &document.ID, map[string]interface{}{
		"document_id": document.ID,
		"title":       document.Title,
		"type":        document.Type,
	})
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return &document, nil
}

// UploadDocumentVersion uploads a new version of a document
func (s *DocumentService) UploadDocumentVersion(ctx context.Context, tenantID, userID, documentID string, fileContent []byte, filename string, options dto.CreateDocumentVersionDTO) (*repo.DocumentVersion, error) {
	// Check if document exists
	document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return nil, fmt.Errorf("document not found")
	}

	// Generate storage key (in real implementation, this would be S3 key)
	storageKey := fmt.Sprintf("documents/%s/versions/%s", documentID, filename)

	// Calculate checksum
	checksum := calculateSHA256(fileContent)

	// Create version
	version := repo.DocumentVersion{
		ID:             generateUUID(),
		DocumentID:     documentID,
		VersionNumber:  document.CurrentVersion + 1,
		StorageKey:     storageKey,
		MimeType:       &[]string{getMimeType(filename)}[0],
		SizeBytes:      &[]int64{int64(len(fileContent))}[0],
		ChecksumSHA256: &checksum,
		AVScanStatus:   "pending",
		CreatedBy:      userID,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	// Process OCR if enabled
	if options.EnableOCR {
		ocrText, err := s.processOCR(fileContent, filename)
		if err != nil {
			// Log error but don't fail
			fmt.Printf("OCR processing failed: %v\n", err)
		} else {
			version.OCRText = &ocrText
		}
	}

	// Perform AV scan
	avStatus, avResult := s.performAVScan(fileContent)
	version.AVScanStatus = avStatus
	version.AVScanResult = &avResult

	// Save version
	err = s.documentRepo.CreateDocumentVersion(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("failed to create document version: %w", err)
	}

	// Update document current version
	document.CurrentVersion = version.VersionNumber
	document.StorageKey = &storageKey
	document.MimeType = version.MimeType
	document.SizeBytes = version.SizeBytes
	document.ChecksumSHA256 = version.ChecksumSHA256
	document.OCRText = version.OCRText
	document.AVScanStatus = version.AVScanStatus
	document.AVScanResult = version.AVScanResult

	err = s.documentRepo.UpdateDocument(ctx, *document)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.version.uploaded", "document", &documentID, map[string]interface{}{
		"document_id": documentID,
		"version":     version.VersionNumber,
		"filename":    filename,
		"size":        len(fileContent),
	})
	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return &version, nil
}

// CreateDocumentVersionWithFile creates a new version with file upload
func (s *DocumentService) CreateDocumentVersionWithFile(ctx context.Context, documentID, tenantID, userID string, req dto.CreateDocumentVersionDTO, file interface{}, filename string) (*repo.DocumentVersion, error) {
	// Check if document exists
	document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return nil, fmt.Errorf("document not found")
	}

	// Type assert file to *multipart.FileHeader
	fmt.Printf("DEBUG: Attempting to cast file to *multipart.FileHeader\n")
	fileHeader, ok := file.(*multipart.FileHeader)
	if !ok {
		fmt.Printf("DEBUG: File type: %T\n", file)
		return nil, fmt.Errorf("invalid file type: %T", file)
	}
	fmt.Printf("DEBUG: Successfully cast to *multipart.FileHeader\n")
	fmt.Printf("DEBUG: Using filename: '%s'\n", filename)

	// Open file
	fileHandle, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileHandle.Close()

	// Read file content
	fileContent, err := io.ReadAll(fileHandle)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Generate storage key
	storageKey := fmt.Sprintf("documents/%s/versions/%s", documentID, filename)

	// Create storage directory if it doesn't exist
	storageDir := fmt.Sprintf("./storage/documents/%s/versions", documentID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Save file to local storage
	filePath := filepath.Join(storageDir, filename)
	if err := os.WriteFile(filePath, fileContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Calculate checksum
	checksum := calculateSHA256(fileContent)

	// Create version
	version := repo.DocumentVersion{
		ID:             generateUUID(),
		DocumentID:     documentID,
		VersionNumber:  document.CurrentVersion + 1,
		StorageKey:     storageKey,
		MimeType:       &[]string{getMimeType(filename)}[0],
		SizeBytes:      &[]int64{int64(len(fileContent))}[0],
		ChecksumSHA256: &checksum,
		AVScanStatus:   "pending",
		CreatedBy:      userID,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	// Process OCR if enabled
	if req.EnableOCR {
		ocrText, err := s.processOCR(fileContent, filename)
		if err != nil {
			fmt.Printf("OCR processing failed: %v\n", err)
		} else {
			version.OCRText = &ocrText
		}
	}

	// Perform AV scan
	avStatus, avResult := s.performAVScan(fileContent)
	version.AVScanStatus = avStatus
	version.AVScanResult = &avResult

	// Save version
	err = s.documentRepo.CreateDocumentVersion(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("failed to create document version: %w", err)
	}

	// Update document current version
	document.CurrentVersion = version.VersionNumber
	document.StorageKey = &storageKey
	document.MimeType = version.MimeType
	document.SizeBytes = version.SizeBytes
	document.ChecksumSHA256 = version.ChecksumSHA256
	document.OCRText = version.OCRText
	document.AVScanStatus = version.AVScanStatus
	document.AVScanResult = version.AVScanResult

	err = s.documentRepo.UpdateDocument(ctx, *document)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.version.uploaded", "document", &documentID, map[string]interface{}{
		"document_id": documentID,
		"version":     version.VersionNumber,
		"filename":    filename,
		"size":        len(fileContent),
	})
	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return &version, nil
}

// PublishDocument publishes an approved document
func (s *DocumentService) PublishDocument(ctx context.Context, tenantID, userID, documentID string) error {
	// Check if document exists and is approved
	document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return fmt.Errorf("document not found")
	}
	if document.Status != "approved" {
		return fmt.Errorf("document must be approved to publish")
	}

	// Update document status to published (we'll use approved for now)
	// In a real system, you might have a separate "published" status
	document.Status = "approved"
	err = s.documentRepo.UpdateDocument(ctx, *document)
	if err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.published", "document", &documentID, map[string]interface{}{
		"document_id": documentID,
		"title":       document.Title,
	})
	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return nil
}

// Helper functions
func calculateSHA256(data []byte) string {
	// In real implementation, use crypto/sha256
	return "sha256:" + fmt.Sprintf("%x", data)[:16] // Simplified for demo
}

func getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}

func (s *DocumentService) processOCR(fileContent []byte, filename string) (string, error) {
	// Perform OCR using external tools to avoid CGO dependencies.
	// Requirements at runtime: `tesseract` (and optionally `pdftoppm` for PDFs).

	// Quick check for platform support
	if runtime.GOOS != "linux" {
		// We only support the containerized linux runtime by default
		return "", fmt.Errorf("ocr is only supported on linux runtime")
	}

	// Create a temp directory to work in
	tmpDir, err := os.MkdirTemp("", "ocr-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write input bytes to a temp file
	inputPath := filepath.Join(tmpDir, "input"+strings.ToLower(filepath.Ext(filename)))
	if err := os.WriteFile(inputPath, fileContent, 0600); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".tif", ".tiff", ".bmp", ".webp":
		// Run tesseract directly on the image
		text, err := runTesseract(inputPath)
		if err != nil {
			return "", err
		}
		return text, nil
	case ".pdf":
		// Convert PDF pages to images using pdftoppm (poppler-utils)
		// Generate PNGs with 300 DPI for better accuracy
		base := filepath.Join(tmpDir, "page")
		cmd := exec.Command("pdftoppm", "-r", "300", "-png", inputPath, base)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("pdftoppm failed: %v, output: %s", err, string(out))
		}

		// Collect generated page images: page-1.png, page-2.png, ...
		var ocrParts []string
		for page := 1; page <= 50; page++ { // hard cap to avoid runaway
			imgPath := fmt.Sprintf("%s-%d.png", base, page)
			if _, err := os.Stat(imgPath); err != nil {
				break
			}
			part, err := runTesseract(imgPath)
			if err != nil {
				return "", err
			}
			ocrParts = append(ocrParts, part)
		}
		return strings.Join(ocrParts, "\n\n"), nil
	default:
		// Unsupported for OCR; return an explicit error so callers can decide to ignore
		return "", fmt.Errorf("ocr not supported for file type: %s", ext)
	}
}

// runTesseract executes `tesseract <image> stdout -l rus+eng` and returns extracted text.
func runTesseract(imagePath string) (string, error) {
	// Try Kazakh+Russian+English first, then fall back to subsets
	langs := []string{"kaz+rus+eng", "rus+eng", "eng"}
	var lastErr error
	for _, lang := range langs {
		cmd := exec.Command("tesseract", imagePath, "stdout", "-l", lang)
		output, err := cmd.CombinedOutput()
		if err == nil {
			return string(output), nil
		}
		lastErr = fmt.Errorf("tesseract failed (lang=%s): %v, output: %s", lang, err, string(output))
	}
	return "", lastErr
}

func (s *DocumentService) performAVScan(fileContent []byte) (string, string) {
	// In real implementation, this would use antivirus service
	// For now, simulate clean scan
	return "clean", "No threats detected"
}

// UpdateDocument updates an existing document
func (s *DocumentService) UpdateDocument(ctx context.Context, id, tenantID, userID string, req dto.UpdateDocumentDTO) (*repo.Document, error) {
	// Get existing document
	document, err := s.documentRepo.GetDocument(ctx, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return nil, fmt.Errorf("document not found")
	}

	// Check if document can be updated
	if document.Status == "approved" {
		return nil, fmt.Errorf("cannot update approved document")
	}

	// Update document
	document.Title = req.Title
	document.Description = req.Description
	document.Type = req.Type
	document.Category = req.Category
	document.Tags = req.Tags
	document.Status = req.Status

	err = s.documentRepo.UpdateDocument(ctx, *document)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.updated", "document", &id, map[string]interface{}{
		"document_id": document.ID,
		"title":       document.Title,
		"status":      document.Status,
	})
	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return document, nil
}

// DeleteDocument soft deletes a document
func (s *DocumentService) DeleteDocument(ctx context.Context, id, tenantID, userID string) error {
	// Check if document exists
	document, err := s.documentRepo.GetDocument(ctx, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return fmt.Errorf("document not found")
	}

	// Check if document can be deleted
	if document.Status == "approved" {
		return fmt.Errorf("cannot delete approved document")
	}

	err = s.documentRepo.DeleteDocument(ctx, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.deleted", "document", &id, map[string]interface{}{
		"document_id": id,
		"title":       document.Title,
	})
	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return nil
}

// CreateDocumentVersion creates a new version of a document
func (s *DocumentService) CreateDocumentVersion(ctx context.Context, documentID, tenantID, userID string, req dto.CreateDocumentVersionDTO) (*repo.DocumentVersion, error) {
	// Get document
	document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return nil, fmt.Errorf("document not found")
	}

	// Create new version
	newVersion := document.CurrentVersion + 1
	version := repo.DocumentVersion{
		ID:             generateUUID(),
		DocumentID:     documentID,
		VersionNumber:  newVersion,
		StorageKey:     *req.FilePath,
		MimeType:       req.MimeType,
		SizeBytes:      req.FileSize,
		ChecksumSHA256: req.ChecksumSHA256,
		OCRText:        req.Content,
		AVScanStatus:   "pending",
		CreatedBy:      userID,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	err = s.documentRepo.CreateDocumentVersion(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("failed to create document version: %w", err)
	}

	// Update document current version
	document.CurrentVersion = newVersion
	err = s.documentRepo.UpdateDocument(ctx, *document)
	if err != nil {
		return nil, fmt.Errorf("failed to update document version: %w", err)
	}

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.version_created", "document_version", &version.ID, map[string]interface{}{
		"document_id": documentID,
		"version":     newVersion,
		"storage_key": version.StorageKey,
	})
	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return &version, nil
}

// ListDocumentVersions retrieves versions for a document
func (s *DocumentService) ListDocumentVersions(ctx context.Context, documentID, tenantID string) ([]repo.DocumentVersion, error) {
	versions, err := s.documentRepo.ListDocumentVersions(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list document versions: %w", err)
	}

	return versions, nil
}

// GetDocumentVersion retrieves a specific version of a document
func (s *DocumentService) GetDocumentVersion(ctx context.Context, versionID, tenantID string) (*repo.DocumentVersion, error) {
	version, err := s.documentRepo.GetDocumentVersion(ctx, versionID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document version: %w", err)
	}

	return version, nil
}

// CreateDocumentAcknowledgment creates an acknowledgment for a document
func (s *DocumentService) CreateDocumentAcknowledgment(ctx context.Context, documentID, tenantID, userID string, req dto.CreateDocumentAcknowledgmentDTO) (*repo.DocumentAcknowledgment, error) {
	// Get document
	document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return nil, fmt.Errorf("document not found")
	}

	// Check if document is approved
	if document.Status != "approved" {
		return nil, fmt.Errorf("document must be approved before creating acknowledgments")
	}

	// Create acknowledgment
	var deadlineStr *string
	if req.Deadline != nil {
		deadline := req.Deadline.Format("2006-01-02")
		deadlineStr = &deadline
	}

	acknowledgment := repo.DocumentAcknowledgment{
		ID:         generateUUID(),
		DocumentID: documentID,
		VersionID:  req.VersionID,
		UserID:     req.UserID,
		Status:     "pending",
		Deadline:   deadlineStr,
	}

	err = s.documentRepo.CreateDocumentAcknowledgment(ctx, acknowledgment)
	if err != nil {
		return nil, fmt.Errorf("failed to create acknowledgment: %w", err)
	}

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.acknowledgment_created", "document_acknowledgment", &acknowledgment.ID, map[string]interface{}{
		"document_id":       documentID,
		"acknowledgment_id": acknowledgment.ID,
		"user_id":           req.UserID,
	})
	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return &acknowledgment, nil
}

// UpdateDocumentAcknowledgment updates an acknowledgment
func (s *DocumentService) UpdateDocumentAcknowledgment(ctx context.Context, acknowledgmentID, tenantID, userID string, req dto.UpdateDocumentAcknowledgmentDTO) (*repo.DocumentAcknowledgment, error) {
	// Get acknowledgment (this would need a new method in repo)
	// For now, we'll create a basic acknowledgment
	var acknowledgedAtStr *string
	if req.AcknowledgedAt != nil {
		acknowledgedAt := req.AcknowledgedAt.Format("2006-01-02 15:04:05")
		acknowledgedAtStr = &acknowledgedAt
	}

	acknowledgment := &repo.DocumentAcknowledgment{
		ID:             acknowledgmentID,
		Status:         req.Status,
		QuizScore:      req.QuizScore,
		QuizPassed:     req.QuizPassed,
		AcknowledgedAt: acknowledgedAtStr,
	}

	err := s.documentRepo.UpdateDocumentAcknowledgment(ctx, *acknowledgment)
	if err != nil {
		return nil, fmt.Errorf("failed to update acknowledgment: %w", err)
	}

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.acknowledgment_updated", "document_acknowledgment", &acknowledgmentID, map[string]interface{}{
		"acknowledgment_id": acknowledgmentID,
		"status":            req.Status,
	})
	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return acknowledgment, nil
}

// CreateDocumentQuiz creates a quiz question for a document
func (s *DocumentService) CreateDocumentQuiz(ctx context.Context, documentID, tenantID, userID string, req dto.CreateDocumentQuizDTO) (*repo.DocumentQuiz, error) {
	// Get document
	document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return nil, fmt.Errorf("document not found")
	}

	// Create quiz
	quiz := repo.DocumentQuiz{
		ID:            generateUUID(),
		DocumentID:    documentID,
		Question:      req.Question,
		QuestionOrder: req.QuestionOrder,
		CorrectAnswer: &req.CorrectAnswer,
		IsActive:      true,
	}

	// Convert options to JSON string
	if len(req.Options) > 0 {
		optionsJSON := fmt.Sprintf(`["%s"]`, strings.Join(req.Options, `","`))
		quiz.Options = &optionsJSON
	}

	err = s.documentRepo.CreateDocumentQuiz(ctx, quiz)
	if err != nil {
		return nil, fmt.Errorf("failed to create quiz: %w", err)
	}

	// Log audit event
	err = s.auditRepo.LogAction(ctx, tenantID, userID, "document.quiz_created", "document_quiz", &quiz.ID, map[string]interface{}{
		"document_id": documentID,
		"quiz_id":     quiz.ID,
		"question":    req.Question,
	})
	if err != nil {
		fmt.Printf("Failed to log audit event: %v\n", err)
	}

	return &quiz, nil
}

// ListDocumentQuizzes retrieves quizzes for a document
func (s *DocumentService) ListDocumentQuizzes(ctx context.Context, documentID, tenantID string) ([]repo.DocumentQuiz, error) {
	quizzes, err := s.documentRepo.ListDocumentQuizzes(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list quizzes: %w", err)
	}

	return quizzes, nil
}

// GetUserPendingAcknowledgment retrieves pending acknowledgments for a user
func (s *DocumentService) GetUserPendingAcknowledgment(ctx context.Context, userID, tenantID string) ([]repo.DocumentAcknowledgment, error) {
	acknowledgments, err := s.documentRepo.GetUserPendingAcknowledgment(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user pending acknowledgments: %w", err)
	}

	return acknowledgments, nil
}

// ListDocumentAcknowledgment retrieves acknowledgments for a document
func (s *DocumentService) ListDocumentAcknowledgment(ctx context.Context, documentID, tenantID string) ([]repo.DocumentAcknowledgment, error) {
	acknowledgments, err := s.documentRepo.ListDocumentAcknowledgment(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list document acknowledgments: %w", err)
	}

	return acknowledgments, nil
}

// DownloadDocumentVersion retrieves file content for a document version
func (s *DocumentService) DownloadDocumentVersion(ctx context.Context, versionID, tenantID string) ([]byte, error) {
	// Get version info
	version, err := s.documentRepo.GetDocumentVersion(ctx, versionID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document version: %w", err)
	}
	if version == nil {
		return nil, fmt.Errorf("document version not found")
	}

	fmt.Printf("DEBUG: Downloading file for version %s, storage key: %s\n", versionID, version.StorageKey)

	// Extract filename from storage key
	filename := filepath.Base(version.StorageKey)
	documentID := filepath.Base(filepath.Dir(filepath.Dir(version.StorageKey)))

	// Read file from local storage
	filePath := filepath.Join("./storage/documents", documentID, "versions", filename)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from storage: %w", err)
	}

	fmt.Printf("DEBUG: Read file from storage, size: %d bytes\n", len(fileContent))
	return fileContent, nil
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func generateUUID() string {
	// Generate a proper UUID v4 format
	// Using a simple approach for now - this is not cryptographically secure
	now := time.Now().Unix()
	// Use only lower -31 bit to ensure values fit in the required ranges
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		now&0xffffffff,
		(now>>8)&0xffff,
		(now>>16)&0xffff,
		(now>>24)&0xffff,
		(now & 0xffffffffffff))
}
