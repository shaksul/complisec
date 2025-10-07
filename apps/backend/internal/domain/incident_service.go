package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

type IncidentService struct {
	incidentRepo           IncidentRepoInterface
	userRepo               UserRepoInterface
	assetRepo              AssetRepoInterface
	riskRepo               RiskRepoInterface
	documentStorageService DocumentStorageServiceInterface
}

func NewIncidentService(incidentRepo IncidentRepoInterface, userRepo UserRepoInterface, assetRepo AssetRepoInterface, riskRepo RiskRepoInterface, documentStorageService DocumentStorageServiceInterface) *IncidentService {
	return &IncidentService{
		incidentRepo:           incidentRepo,
		userRepo:               userRepo,
		assetRepo:              assetRepo,
		riskRepo:               riskRepo,
		documentStorageService: documentStorageService,
	}
}

func (s *IncidentService) CreateIncident(ctx context.Context, tenantID string, req dto.CreateIncidentRequest, reportedBy string) (*repo.Incident, error) {
	log.Printf("DEBUG: incident_service.CreateIncident tenant=%s title=%s", tenantID, req.Title)

	// Validate assigned user exists if provided
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		user, err := s.userRepo.GetByID(ctx, *req.AssignedTo)
		if err != nil {
			log.Printf("ERROR: incident_service.CreateIncident GetByID assigned user: %v", err)
			return nil, err
		}
		if user == nil {
			log.Printf("WARN: incident_service.CreateIncident assigned user not found id=%s", *req.AssignedTo)
			return nil, errors.New("assigned user not found")
		}
	}

	// Validate assets exist if provided
	for _, assetID := range req.AssetIDs {
		asset, err := s.assetRepo.GetByID(ctx, assetID)
		if err != nil {
			log.Printf("ERROR: incident_service.CreateIncident GetByID asset: %v", err)
			return nil, err
		}
		if asset == nil {
			log.Printf("WARN: incident_service.CreateIncident asset not found id=%s", assetID)
			return nil, errors.New("asset not found")
		}
	}

	// Validate risks exist if provided
	for _, riskID := range req.RiskIDs {
		risk, err := s.riskRepo.GetByIDWithTenant(ctx, riskID, tenantID)
		if err != nil {
			log.Printf("ERROR: incident_service.CreateIncident GetByIDWithTenant risk: %v", err)
			return nil, err
		}
		if risk == nil {
			log.Printf("WARN: incident_service.CreateIncident risk not found id=%s", riskID)
			return nil, errors.New("risk not found")
		}
	}

	// Set detected time
	detectedAt := time.Now()
	if req.DetectedAt != nil {
		detectedAt = *req.DetectedAt
	}

	// Create incident
	incident := repo.Incident{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Status:      dto.IncidentStatusNew,
		Criticality: req.Criticality,
		Source:      req.Source,
		ReportedBy:  reportedBy,
		AssignedTo:  req.AssignedTo,
		DetectedAt:  detectedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.incidentRepo.Create(ctx, &incident)
	if err != nil {
		log.Printf("ERROR: incident_service.CreateIncident Create: %v", err)
		return nil, err
	}

	// Add asset relations
	for _, assetID := range req.AssetIDs {
		err := s.incidentRepo.AddAsset(ctx, incident.ID, assetID)
		if err != nil {
			log.Printf("ERROR: incident_service.CreateIncident AddAsset: %v", err)
			// Continue with other assets, don't fail the entire operation
		}
	}

	// Add risk relations
	for _, riskID := range req.RiskIDs {
		err := s.incidentRepo.AddRisk(ctx, incident.ID, riskID)
		if err != nil {
			log.Printf("ERROR: incident_service.CreateIncident AddRisk: %v", err)
			// Continue with other risks, don't fail the entire operation
		}
	}

	log.Printf("INFO: incident_service.CreateIncident created id=%s", incident.ID)
	return &incident, nil
}

func (s *IncidentService) GetIncident(ctx context.Context, id, tenantID string) (*repo.Incident, error) {
	log.Printf("DEBUG: incident_service.GetIncident id=%s tenant=%s", id, tenantID)

	incident, err := s.incidentRepo.GetByID(ctx, id, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.GetIncident GetByID: %v", err)
		return nil, err
	}

	return incident, nil
}

func (s *IncidentService) UpdateIncident(ctx context.Context, id, tenantID string, req dto.UpdateIncidentRequest, updatedBy string) (*repo.Incident, error) {
	log.Printf("DEBUG: incident_service.UpdateIncident id=%s tenant=%s", id, tenantID)

	// Get existing incident
	incident, err := s.incidentRepo.GetByID(ctx, id, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.UpdateIncident GetByID: %v", err)
		return nil, err
	}

	// Validate assigned user exists if provided
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		user, err := s.userRepo.GetByID(ctx, *req.AssignedTo)
		if err != nil {
			log.Printf("ERROR: incident_service.UpdateIncident GetByID assigned user: %v", err)
			return nil, err
		}
		if user == nil {
			log.Printf("WARN: incident_service.UpdateIncident assigned user not found id=%s", *req.AssignedTo)
			return nil, errors.New("assigned user not found")
		}
	}

	// Update fields
	if req.Title != nil {
		incident.Title = *req.Title
	}
	if req.Description != nil {
		incident.Description = req.Description
	}
	if req.Category != nil {
		incident.Category = *req.Category
	}
	if req.Criticality != nil {
		incident.Criticality = *req.Criticality
	}
	if req.Status != nil {
		incident.Status = *req.Status
		// Set resolved/closed timestamps
		if *req.Status == dto.IncidentStatusResolved && incident.ResolvedAt == nil {
			now := time.Now()
			incident.ResolvedAt = &now
		}
		if *req.Status == dto.IncidentStatusClosed && incident.ClosedAt == nil {
			now := time.Now()
			incident.ClosedAt = &now
		}
	}
	if req.AssignedTo != nil {
		incident.AssignedTo = req.AssignedTo
	}
	if req.DetectedAt != nil {
		incident.DetectedAt = *req.DetectedAt
	}

	incident.UpdatedAt = time.Now()

	err = s.incidentRepo.Update(ctx, incident)
	if err != nil {
		log.Printf("ERROR: incident_service.UpdateIncident Update: %v", err)
		return nil, err
	}

	// Update asset relations if provided
	if req.AssetIDs != nil {
		// Get current assets
		currentAssets, err := s.incidentRepo.GetAssets(ctx, incident.ID)
		if err != nil {
			log.Printf("ERROR: incident_service.UpdateIncident GetAssets: %v", err)
			return nil, err
		}

		// Remove assets not in new list
		for _, asset := range currentAssets {
			found := false
			for _, assetID := range req.AssetIDs {
				if asset.ID == assetID {
					found = true
					break
				}
			}
			if !found {
				err := s.incidentRepo.RemoveAsset(ctx, incident.ID, asset.ID)
				if err != nil {
					log.Printf("ERROR: incident_service.UpdateIncident RemoveAsset: %v", err)
				}
			}
		}

		// Add new assets
		for _, assetID := range req.AssetIDs {
			found := false
			for _, asset := range currentAssets {
				if asset.ID == assetID {
					found = true
					break
				}
			}
			if !found {
				err := s.incidentRepo.AddAsset(ctx, incident.ID, assetID)
				if err != nil {
					log.Printf("ERROR: incident_service.UpdateIncident AddAsset: %v", err)
				}
			}
		}
	}

	// Update risk relations if provided
	if req.RiskIDs != nil {
		// Get current risks
		currentRisks, err := s.incidentRepo.GetRisks(ctx, incident.ID)
		if err != nil {
			log.Printf("ERROR: incident_service.UpdateIncident GetRisks: %v", err)
			return nil, err
		}

		// Remove risks not in new list
		for _, risk := range currentRisks {
			found := false
			for _, riskID := range req.RiskIDs {
				if risk.ID == riskID {
					found = true
					break
				}
			}
			if !found {
				err := s.incidentRepo.RemoveRisk(ctx, incident.ID, risk.ID)
				if err != nil {
					log.Printf("ERROR: incident_service.UpdateIncident RemoveRisk: %v", err)
				}
			}
		}

		// Add new risks
		for _, riskID := range req.RiskIDs {
			found := false
			for _, risk := range currentRisks {
				if risk.ID == riskID {
					found = true
					break
				}
			}
			if !found {
				err := s.incidentRepo.AddRisk(ctx, incident.ID, riskID)
				if err != nil {
					log.Printf("ERROR: incident_service.UpdateIncident AddRisk: %v", err)
				}
			}
		}
	}

	log.Printf("INFO: incident_service.UpdateIncident updated id=%s", incident.ID)
	return incident, nil
}

func (s *IncidentService) DeleteIncident(ctx context.Context, id, tenantID string) error {
	log.Printf("DEBUG: incident_service.DeleteIncident id=%s tenant=%s", id, tenantID)

	err := s.incidentRepo.Delete(ctx, id, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.DeleteIncident Delete: %v", err)
		return err
	}

	log.Printf("INFO: incident_service.DeleteIncident deleted id=%s", id)
	return nil
}

func (s *IncidentService) ListIncidents(ctx context.Context, tenantID string, req dto.IncidentListRequest) ([]*repo.Incident, int, error) {
	log.Printf("DEBUG: incident_service.ListIncidents tenant=%s page=%d page_size=%d", tenantID, req.Page, req.PageSize)

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.Criticality != "" {
		filters["criticality"] = req.Criticality
	}
	if req.Category != "" {
		filters["category"] = req.Category
	}
	if req.AssignedTo != "" {
		filters["assigned_to"] = req.AssignedTo
	}
	if req.Search != "" {
		filters["search"] = req.Search
	}

	offset := (req.Page - 1) * req.PageSize
	incidents, total, err := s.incidentRepo.List(ctx, tenantID, filters, req.PageSize, offset)
	if err != nil {
		log.Printf("ERROR: incident_service.ListIncidents List: %v", err)
		return nil, 0, err
	}

	log.Printf("INFO: incident_service.ListIncidents found %d incidents", len(incidents))
	return incidents, total, nil
}

func (s *IncidentService) AddComment(ctx context.Context, incidentID, tenantID string, req dto.IncidentCommentRequest, userID string) (*repo.IncidentComment, error) {
	log.Printf("DEBUG: incident_service.AddComment incident=%s user=%s", incidentID, userID)

	// Verify incident exists
	_, err := s.incidentRepo.GetByID(ctx, incidentID, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.AddComment GetByID: %v", err)
		return nil, err
	}

	comment := repo.IncidentComment{
		ID:         uuid.New().String(),
		IncidentID: incidentID,
		UserID:     userID,
		Comment:    req.Comment,
		IsInternal: req.IsInternal,
		CreatedAt:  time.Now(),
	}

	err = s.incidentRepo.AddComment(ctx, &comment)
	if err != nil {
		log.Printf("ERROR: incident_service.AddComment AddComment: %v", err)
		return nil, err
	}

	log.Printf("INFO: incident_service.AddComment added comment id=%s", comment.ID)
	return &comment, nil
}

func (s *IncidentService) GetComments(ctx context.Context, incidentID, tenantID string) ([]*repo.IncidentComment, error) {
	log.Printf("DEBUG: incident_service.GetComments incident=%s", incidentID)

	// Verify incident exists
	_, err := s.incidentRepo.GetByID(ctx, incidentID, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.GetComments GetByID: %v", err)
		return nil, err
	}

	comments, err := s.incidentRepo.GetComments(ctx, incidentID)
	if err != nil {
		log.Printf("ERROR: incident_service.GetComments GetComments: %v", err)
		return nil, err
	}

	return comments, nil
}

func (s *IncidentService) AddAction(ctx context.Context, incidentID, tenantID string, req dto.IncidentActionRequest, createdBy string) (*repo.IncidentAction, error) {
	log.Printf("DEBUG: incident_service.AddAction incident=%s user=%s", incidentID, createdBy)

	// Verify incident exists
	_, err := s.incidentRepo.GetByID(ctx, incidentID, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.AddAction GetByID: %v", err)
		return nil, err
	}

	// Validate assigned user exists if provided
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		user, err := s.userRepo.GetByID(ctx, *req.AssignedTo)
		if err != nil {
			log.Printf("ERROR: incident_service.AddAction GetByID assigned user: %v", err)
			return nil, err
		}
		if user == nil {
			log.Printf("WARN: incident_service.AddAction assigned user not found id=%s", *req.AssignedTo)
			return nil, errors.New("assigned user not found")
		}
	}

	action := repo.IncidentAction{
		ID:          uuid.New().String(),
		IncidentID:  incidentID,
		ActionType:  req.ActionType,
		Title:       req.Title,
		Description: req.Description,
		AssignedTo:  req.AssignedTo,
		DueDate:     req.DueDate,
		Status:      dto.ActionStatusPending,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.incidentRepo.AddAction(ctx, &action)
	if err != nil {
		log.Printf("ERROR: incident_service.AddAction AddAction: %v", err)
		return nil, err
	}

	log.Printf("INFO: incident_service.AddAction added action id=%s", action.ID)
	return &action, nil
}

func (s *IncidentService) UpdateAction(ctx context.Context, actionID, tenantID string, req dto.IncidentActionRequest, updatedBy string) (*repo.IncidentAction, error) {
	log.Printf("DEBUG: incident_service.UpdateAction action=%s user=%s", actionID, updatedBy)

	// Get existing action
	actions, err := s.incidentRepo.GetActions(ctx, "")
	if err != nil {
		log.Printf("ERROR: incident_service.UpdateAction GetActions: %v", err)
		return nil, err
	}

	var action *repo.IncidentAction
	for _, a := range actions {
		if a.ID == actionID {
			action = a
			break
		}
	}

	if action == nil {
		return nil, errors.New("action not found")
	}

	// Verify incident belongs to tenant
	incident, err := s.incidentRepo.GetByID(ctx, action.IncidentID, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.UpdateAction GetByID incident: %v", err)
		return nil, err
	}
	if incident == nil {
		return nil, errors.New("incident not found")
	}

	// Validate assigned user exists if provided
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		user, err := s.userRepo.GetByID(ctx, *req.AssignedTo)
		if err != nil {
			log.Printf("ERROR: incident_service.UpdateAction GetByID assigned user: %v", err)
			return nil, err
		}
		if user == nil {
			log.Printf("WARN: incident_service.UpdateAction assigned user not found id=%s", *req.AssignedTo)
			return nil, errors.New("assigned user not found")
		}
	}

	// Update fields
	action.ActionType = req.ActionType
	action.Title = req.Title
	action.Description = req.Description
	action.AssignedTo = req.AssignedTo
	action.DueDate = req.DueDate
	action.UpdatedAt = time.Now()

	err = s.incidentRepo.UpdateAction(ctx, action)
	if err != nil {
		log.Printf("ERROR: incident_service.UpdateAction UpdateAction: %v", err)
		return nil, err
	}

	log.Printf("INFO: incident_service.UpdateAction updated action id=%s", action.ID)
	return action, nil
}

func (s *IncidentService) GetActions(ctx context.Context, incidentID, tenantID string) ([]*repo.IncidentAction, error) {
	log.Printf("DEBUG: incident_service.GetActions incident=%s", incidentID)

	// Verify incident exists
	_, err := s.incidentRepo.GetByID(ctx, incidentID, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.GetActions GetByID: %v", err)
		return nil, err
	}

	actions, err := s.incidentRepo.GetActions(ctx, incidentID)
	if err != nil {
		log.Printf("ERROR: incident_service.GetActions GetActions: %v", err)
		return nil, err
	}

	return actions, nil
}

func (s *IncidentService) GetIncidentMetrics(ctx context.Context, tenantID string) (*repo.IncidentMetricsSummary, error) {
	log.Printf("DEBUG: incident_service.GetIncidentMetrics tenant=%s", tenantID)

	metrics, err := s.incidentRepo.GetIncidentMetrics(ctx, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.GetIncidentMetrics GetIncidentMetrics: %v", err)
		return nil, err
	}

	return metrics, nil
}

func (s *IncidentService) UpdateIncidentStatus(ctx context.Context, id, tenantID string, req dto.IncidentStatusUpdateRequest, updatedBy string) (*repo.Incident, error) {
	log.Printf("DEBUG: incident_service.UpdateIncidentStatus id=%s status=%s", id, req.Status)

	// Get existing incident
	incident, err := s.incidentRepo.GetByID(ctx, id, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.UpdateIncidentStatus GetByID: %v", err)
		return nil, err
	}

	// Update status
	incident.Status = req.Status
	incident.UpdatedAt = time.Now()

	// Set resolved/closed timestamps
	if req.Status == dto.IncidentStatusResolved && incident.ResolvedAt == nil {
		now := time.Now()
		incident.ResolvedAt = &now
	}
	if req.Status == dto.IncidentStatusClosed && incident.ClosedAt == nil {
		now := time.Now()
		incident.ClosedAt = &now
	}

	err = s.incidentRepo.Update(ctx, incident)
	if err != nil {
		log.Printf("ERROR: incident_service.UpdateIncidentStatus Update: %v", err)
		return nil, err
	}

	log.Printf("INFO: incident_service.UpdateIncidentStatus updated id=%s status=%s", id, req.Status)
	return incident, nil
}

// Incident Document methods - интеграция с централизованным хранилищем документов
func (s *IncidentService) UploadIncidentDocument(ctx context.Context, incidentID, tenantID string, file multipart.File, header *multipart.FileHeader, req dto.UploadDocumentDTO, uploadedBy string) (*dto.DocumentDTO, error) {
	log.Printf("DEBUG: incident_service.UploadIncidentDocument incidentID=%s", incidentID)

	// Проверяем, что инцидент существует
	incident, err := s.incidentRepo.GetByID(ctx, incidentID, tenantID)
	if err != nil {
		return nil, err
	}
	if incident == nil {
		return nil, errors.New("incident not found")
	}

	// Обновляем запрос с информацией об инциденте
	req.Name = fmt.Sprintf("%s - %s", incident.Title, req.Name)
	req.Tags = append(req.Tags, "#инциденты")
	if req.LinkedTo == nil {
		req.LinkedTo = &dto.DocumentLinkDTO{
			Module:   "incidents",
			EntityID: incidentID,
		}
	}
	if req.Description == nil {
		req.Description = incidentStringPtr(fmt.Sprintf("Document for incident: %s", incident.Title))
	}

	// Загружаем документ в централизованное хранилище
	document, err := s.documentStorageService.UploadDocument(ctx, tenantID, file, header, req, uploadedBy)
	if err != nil {
		log.Printf("ERROR: incident_service.UploadIncidentDocument UploadDocument: %v", err)
		return nil, err
	}

	log.Printf("DEBUG: incident_service.UploadIncidentDocument success documentID=%s", document.ID)
	return document, nil
}

func (s *IncidentService) GetIncidentDocuments(ctx context.Context, incidentID, tenantID string) ([]dto.DocumentDTO, error) {
	log.Printf("DEBUG: incident_service.GetIncidentDocuments incidentID=%s", incidentID)

	// Проверяем, что инцидент существует
	incident, err := s.incidentRepo.GetByID(ctx, incidentID, tenantID)
	if err != nil {
		return nil, err
	}
	if incident == nil {
		return nil, errors.New("incident not found")
	}

	// Получаем документы из централизованного хранилища
	documents, err := s.documentStorageService.GetModuleDocuments(ctx, "incidents", incidentID, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_service.GetIncidentDocuments GetModuleDocuments: %v", err)
		return nil, err
	}

	log.Printf("DEBUG: incident_service.GetIncidentDocuments found %d documents", len(documents))
	return documents, nil
}

func (s *IncidentService) LinkExistingDocumentToIncident(ctx context.Context, incidentID, documentID, tenantID, linkedBy string) error {
	log.Printf("DEBUG: incident_service.LinkExistingDocumentToIncident incidentID=%s documentID=%s", incidentID, documentID)

	// Проверяем, что инцидент существует
	incident, err := s.incidentRepo.GetByID(ctx, incidentID, tenantID)
	if err != nil {
		return err
	}
	if incident == nil {
		return errors.New("incident not found")
	}

	// Проверяем, что документ существует
	_, err = s.documentStorageService.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return err
	}

	// Связываем документ с инцидентом
	err = s.documentStorageService.LinkDocumentToModule(ctx, documentID, "incidents", incidentID, "attachment", fmt.Sprintf("Linked to incident: %s", incident.Title), linkedBy)
	if err != nil {
		log.Printf("ERROR: incident_service.LinkExistingDocumentToIncident LinkDocumentToModule: %v", err)
		return err
	}

	log.Printf("DEBUG: incident_service.LinkExistingDocumentToIncident success")
	return nil
}

func (s *IncidentService) UnlinkDocumentFromIncident(ctx context.Context, incidentID, documentID, tenantID, unlinkedBy string) error {
	log.Printf("DEBUG: incident_service.UnlinkDocumentFromIncident incidentID=%s documentID=%s", incidentID, documentID)

	// Проверяем, что инцидент существует
	incident, err := s.incidentRepo.GetByID(ctx, incidentID, tenantID)
	if err != nil {
		return err
	}
	if incident == nil {
		return errors.New("incident not found")
	}

	// Отвязываем документ от инцидента
	err = s.documentStorageService.UnlinkDocumentFromModule(ctx, documentID, "incidents", incidentID, unlinkedBy)
	if err != nil {
		log.Printf("ERROR: incident_service.UnlinkDocumentFromIncident UnlinkDocumentFromModule: %v", err)
		return err
	}

	log.Printf("DEBUG: incident_service.UnlinkDocumentFromIncident success")
	return nil
}

func (s *IncidentService) DeleteIncidentDocument(ctx context.Context, incidentID, documentID, tenantID, deletedBy string) error {
	log.Printf("DEBUG: incident_service.DeleteIncidentDocument incidentID=%s documentID=%s", incidentID, documentID)

	// Проверяем, что инцидент существует
	incident, err := s.incidentRepo.GetByID(ctx, incidentID, tenantID)
	if err != nil {
		return err
	}
	if incident == nil {
		return errors.New("incident not found")
	}

	// Получаем информацию о документе
	_, err = s.documentStorageService.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return err
	}

	// Удаляем документ из централизованного хранилища
	err = s.documentStorageService.DeleteDocument(ctx, documentID, tenantID, deletedBy)
	if err != nil {
		log.Printf("ERROR: incident_service.DeleteIncidentDocument DeleteDocument: %v", err)
		return err
	}

	log.Printf("DEBUG: incident_service.DeleteIncidentDocument success")
	return nil
}

// Helper function to create string pointer
func incidentStringPtr(s string) *string {
	return &s
}
