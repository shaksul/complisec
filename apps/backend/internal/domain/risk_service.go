package domain

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

type RiskService struct {
	riskRepo        *repo.RiskRepo
	auditRepo       *repo.AuditRepo
	documentService DocumentServiceInterface
}

func NewRiskService(riskRepo *repo.RiskRepo, auditRepo *repo.AuditRepo, documentService DocumentServiceInterface) *RiskService {
	return &RiskService{
		riskRepo:        riskRepo,
		auditRepo:       auditRepo,
		documentService: documentService,
	}
}

func (s *RiskService) CreateRisk(ctx context.Context, tenantID, title string, description, category *string, likelihood, impact int, ownerUserID, assetID *string, methodology, strategy *string, dueDate *time.Time) (*repo.Risk, error) {
	// Calculate risk level automatically
	level, _ := dto.CalculateRiskLevel(likelihood, impact)

	risk := repo.Risk{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Title:       title,
		Description: description,
		Category:    category,
		Likelihood:  &likelihood,
		Impact:      &impact,
		Level:       &level,
		Status:      dto.RiskStatusNew,
		OwnerUserID: ownerUserID,
		AssetID:     assetID,
		Methodology: methodology,
		Strategy:    strategy,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.riskRepo.Create(ctx, risk)
	if err != nil {
		return nil, err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, tenantID, "system", "create", "risk", &risk.ID, risk)

	return &risk, nil
}

func (s *RiskService) GetRisk(ctx context.Context, id string) (*repo.Risk, error) {
	return s.riskRepo.GetByID(ctx, id)
}

func (s *RiskService) ListRisks(ctx context.Context, tenantID string, filters map[string]interface{}, sortField, sortDirection string) ([]repo.Risk, error) {
	return s.riskRepo.ListWithFilters(ctx, tenantID, filters, sortField, sortDirection)
}

func (s *RiskService) UpdateRisk(ctx context.Context, id, title string, description, category *string, likelihood, impact int, ownerUserID, assetID *string, methodology, strategy *string, dueDate *time.Time) error {
	risk, err := s.riskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if risk == nil {
		return nil
	}

	// Calculate new risk level if likelihood or impact changed
	oldLevel := risk.Level
	level, _ := dto.CalculateRiskLevel(likelihood, impact)

	risk.Title = title
	risk.Description = description
	risk.Category = category
	risk.Likelihood = &likelihood
	risk.Impact = &impact
	risk.Level = &level
	risk.OwnerUserID = ownerUserID
	risk.AssetID = assetID
	risk.Methodology = methodology
	risk.Strategy = strategy
	risk.DueDate = dueDate
	risk.UpdatedAt = time.Now()

	err = s.riskRepo.Update(ctx, *risk)
	if err != nil {
		return err
	}

	// Log audit with level change if applicable
	auditData := map[string]interface{}{
		"title":      title,
		"likelihood": likelihood,
		"impact":     impact,
		"level":      level,
	}
	if oldLevel != nil && *oldLevel != level {
		auditData["level_changed"] = map[string]interface{}{
			"old_level": *oldLevel,
			"new_level": level,
		}

		// Check for high/critical risk level escalation
		s.checkRiskLevelEscalation(ctx, risk, *oldLevel, level)
	}

	s.auditRepo.LogAction(ctx, risk.TenantID, "system", "update", "risk", &id, auditData)

	return nil
}

func (s *RiskService) UpdateRiskStatus(ctx context.Context, id, status string) error {
	risk, err := s.riskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if risk == nil {
		return nil
	}

	err = s.riskRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, risk.TenantID, "system", "update_status", "risk", &id, map[string]string{"status": status})

	return nil
}

func (s *RiskService) DeleteRisk(ctx context.Context, id string) error {
	risk, err := s.riskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if risk == nil {
		return nil
	}

	err = s.riskRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, risk.TenantID, "system", "delete", "risk", &id, nil)

	return nil
}

// Risk Controls methods
func (s *RiskService) AddControl(ctx context.Context, riskID string, controlID, controlName, controlType, implementationStatus string, effectiveness, description *string, createdBy string) error {
	control := repo.RiskControl{
		ID:                   uuid.New().String(),
		RiskID:               riskID,
		ControlID:            controlID,
		ControlName:          controlName,
		ControlType:          controlType,
		ImplementationStatus: implementationStatus,
		Effectiveness:        effectiveness,
		Description:          description,
		CreatedBy:            &createdBy,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	err := s.riskRepo.AddControl(ctx, control)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "add_control", "risk", &riskID, control)

	return nil
}

func (s *RiskService) GetControls(ctx context.Context, riskID string) ([]repo.RiskControl, error) {
	return s.riskRepo.GetControls(ctx, riskID)
}

func (s *RiskService) UpdateControl(ctx context.Context, controlID, controlName, controlType, implementationStatus string, effectiveness, description *string) error {
	control := repo.RiskControl{
		ID:                   controlID,
		ControlName:          controlName,
		ControlType:          controlType,
		ImplementationStatus: implementationStatus,
		Effectiveness:        effectiveness,
		Description:          description,
		UpdatedAt:            time.Now(),
	}

	err := s.riskRepo.UpdateControl(ctx, control)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "update_control", "risk", &controlID, control)

	return nil
}

func (s *RiskService) DeleteControl(ctx context.Context, controlID string) error {
	err := s.riskRepo.DeleteControl(ctx, controlID)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "delete_control", "risk", &controlID, nil)

	return nil
}

// Risk Comments methods
func (s *RiskService) AddComment(ctx context.Context, riskID, userID, comment string, isInternal bool) error {
	riskComment := repo.RiskComment{
		ID:         uuid.New().String(),
		RiskID:     riskID,
		UserID:     userID,
		Comment:    comment,
		IsInternal: isInternal,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := s.riskRepo.AddComment(ctx, riskComment)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "add_comment", "risk", &riskID, riskComment)

	return nil
}

func (s *RiskService) GetComments(ctx context.Context, riskID string, includeInternal bool) ([]repo.RiskComment, error) {
	return s.riskRepo.GetComments(ctx, riskID, includeInternal)
}

// Risk History methods
func (s *RiskService) AddHistory(ctx context.Context, riskID, fieldChanged string, oldValue, newValue, changeReason *string, changedBy string) error {
	history := repo.RiskHistory{
		ID:           uuid.New().String(),
		RiskID:       riskID,
		FieldChanged: fieldChanged,
		OldValue:     oldValue,
		NewValue:     newValue,
		ChangeReason: changeReason,
		ChangedBy:    changedBy,
		ChangedAt:    time.Now(),
	}

	err := s.riskRepo.AddHistory(ctx, history)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "add_history", "risk", &riskID, history)

	return nil
}

func (s *RiskService) GetHistory(ctx context.Context, riskID string) ([]repo.RiskHistory, error) {
	return s.riskRepo.GetHistory(ctx, riskID)
}

// Risk Attachments methods
func (s *RiskService) AddAttachment(ctx context.Context, riskID, fileName, filePath string, fileSize int64, mimeType string, fileHash, description *string, uploadedBy string) error {
	attachment := repo.RiskAttachment{
		ID:          uuid.New().String(),
		RiskID:      riskID,
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    fileSize,
		MimeType:    mimeType,
		FileHash:    fileHash,
		Description: description,
		UploadedBy:  uploadedBy,
		UploadedAt:  time.Now(),
	}

	err := s.riskRepo.AddAttachment(ctx, attachment)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "add_attachment", "risk", &riskID, attachment)

	return nil
}

func (s *RiskService) GetAttachments(ctx context.Context, riskID string) ([]repo.RiskAttachment, error) {
	return s.riskRepo.GetAttachments(ctx, riskID)
}

func (s *RiskService) DeleteAttachment(ctx context.Context, attachmentID string) error {
	err := s.riskRepo.DeleteAttachment(ctx, attachmentID)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "delete_attachment", "risk", &attachmentID, nil)

	return nil
}

// Risk Tags methods
func (s *RiskService) AddTag(ctx context.Context, riskID, tagName, tagColor string, createdBy *string) error {
	tag := repo.RiskTag{
		ID:        uuid.New().String(),
		RiskID:    riskID,
		TagName:   tagName,
		TagColor:  tagColor,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}

	err := s.riskRepo.AddTag(ctx, tag)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "add_tag", "risk", &riskID, tag)

	return nil
}

func (s *RiskService) GetTags(ctx context.Context, riskID string) ([]repo.RiskTag, error) {
	return s.riskRepo.GetTags(ctx, riskID)
}

func (s *RiskService) DeleteTag(ctx context.Context, tagID string) error {
	err := s.riskRepo.DeleteTag(ctx, tagID)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "delete_tag", "risk", &tagID, nil)

	return nil
}

func (s *RiskService) DeleteTagByName(ctx context.Context, riskID, tagName string) error {
	err := s.riskRepo.DeleteTagByName(ctx, riskID, tagName)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, "", "system", "delete_tag_by_name", "risk", &riskID, map[string]string{"tag_name": tagName})

	return nil
}

// checkRiskLevelEscalation checks if risk level has escalated to High or Critical and sends notifications
func (s *RiskService) checkRiskLevelEscalation(ctx context.Context, risk *repo.Risk, oldLevel, newLevel int) {
	// Define risk level thresholds
	const (
		HighLevel     = 6 // 3*2 or 2*3 or higher combinations
		CriticalLevel = 8 // 4*2 or 2*4 or higher combinations
	)

	// Check if escalation occurred
	var escalationType string
	var shouldNotify bool

	if oldLevel < HighLevel && newLevel >= HighLevel && newLevel < CriticalLevel {
		escalationType = "high"
		shouldNotify = true
	} else if oldLevel < CriticalLevel && newLevel >= CriticalLevel {
		escalationType = "critical"
		shouldNotify = true
	}

	if shouldNotify {
		// Log the escalation event
		log.Printf("WARNING: Risk level escalation detected - Risk ID: %s, Old Level: %d, New Level: %d, Type: %s",
			risk.ID, oldLevel, newLevel, escalationType)

		// Create notification data
		notificationData := map[string]interface{}{
			"risk_id":         risk.ID,
			"risk_title":      risk.Title,
			"old_level":       oldLevel,
			"new_level":       newLevel,
			"escalation_type": escalationType,
			"tenant_id":       risk.TenantID,
		}

		// Log audit for escalation
		s.auditRepo.LogAction(ctx, risk.TenantID, "system", "risk_escalation", "risk", &risk.ID, notificationData)

		// TODO: In a real implementation, you would:
		// 1. Send email notifications to risk owners and management
		// 2. Create in-app notifications
		// 3. Send Slack/Teams notifications
		// 4. Trigger automated workflows

		log.Printf("NOTIFICATION: Risk '%s' escalated to %s level (%d). Consider immediate action.",
			risk.Title, escalationType, newLevel)
	}
}

// Risk Document methods - интеграция с модулем документов
func (s *RiskService) UploadRiskDocument(ctx context.Context, riskID, tenantID string, file multipart.File, header *multipart.FileHeader, req dto.UploadDocumentDTO, uploadedBy string) (*dto.DocumentDTO, error) {
	// Загружаем документ в модуль документов
	document, err := s.documentService.UploadDocument(ctx, tenantID, file, header, req, uploadedBy)
	if err != nil {
		return nil, err
	}

	// Создаем связь с риском
	link := dto.CreateDocumentLinkDTO{
		DocumentID: document.ID,
		Module:     "risks",
		EntityID:   riskID,
		LinkType:   "attachment",
		Description: stringPtr("Risk attachment"),
	}

	err = s.documentService.AddDocumentLink(ctx, document.ID, link)
	if err != nil {
		// Если не удалось создать связь, удаляем документ
		s.documentService.DeleteDocument(ctx, document.ID, tenantID, uploadedBy)
		return nil, err
	}

	// Логируем аудит
	s.auditRepo.LogAction(ctx, tenantID, uploadedBy, "upload_document", "risk", &riskID, map[string]interface{}{
		"document_id": document.ID,
		"file_name":   document.OriginalName,
		"file_size":   document.FileSize,
		"mime_type":   document.MimeType,
	})

	return document, nil
}

func (s *RiskService) GetRiskDocuments(ctx context.Context, riskID, tenantID string) ([]dto.DocumentDTO, error) {
	// Получаем документы, связанные с риском
	filters := dto.FileDocumentFiltersDTO{
		Module:   stringPtr("risks"),
		EntityID: stringPtr(riskID),
	}

	return s.documentService.ListDocuments(ctx, tenantID, filters)
}

func (s *RiskService) DeleteRiskDocument(ctx context.Context, riskID, documentID, tenantID, deletedBy string) error {
	// Проверяем, что документ действительно связан с этим риском
	documents, err := s.GetRiskDocuments(ctx, riskID, tenantID)
	if err != nil {
		return err
	}

	var found bool
	for _, doc := range documents {
		if doc.ID == documentID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("document not found or not linked to this risk")
	}

	// Удаляем документ
	err = s.documentService.DeleteDocument(ctx, documentID, tenantID, deletedBy)
	if err != nil {
		return err
	}

	// Логируем аудит
	s.auditRepo.LogAction(ctx, tenantID, deletedBy, "delete_document", "risk", &riskID, map[string]interface{}{
		"document_id": documentID,
	})

	return nil
}

func (s *RiskService) LinkExistingDocument(ctx context.Context, riskID, documentID, tenantID, linkedBy string) error {
	// Проверяем, что документ существует
	document, err := s.documentService.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return err
	}

	// Создаем связь с риском
	link := dto.CreateDocumentLinkDTO{
		DocumentID:  documentID,
		Module:      "risks",
		EntityID:    riskID,
		LinkType:    "attachment",
		Description: stringPtr("Risk attachment"),
	}

	err = s.documentService.AddDocumentLink(ctx, documentID, link)
	if err != nil {
		return err
	}

	// Логируем аудит
	s.auditRepo.LogAction(ctx, tenantID, linkedBy, "link_document", "risk", &riskID, map[string]interface{}{
		"document_id": documentID,
		"file_name":   document.OriginalName,
	})

	return nil
}

func (s *RiskService) UnlinkDocument(ctx context.Context, riskID, documentID, tenantID, unlinkedBy string) error {
	// Удаляем связь с риском
	err := s.documentService.RemoveDocumentLink(ctx, documentID, "risks", riskID)
	if err != nil {
		return err
	}

	// Логируем аудит
	s.auditRepo.LogAction(ctx, tenantID, unlinkedBy, "unlink_document", "risk", &riskID, map[string]interface{}{
		"document_id": documentID,
	})

	return nil
}

