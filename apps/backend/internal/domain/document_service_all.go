package domain

import (
	"context"
	"fmt"

	"risknexus/backend/internal/dto"
)

// ListAllDocuments gets ALL documents including those linked to modules (for file storage)
func (s *DocumentService) ListAllDocuments(ctx context.Context, tenantID string, filters dto.FileDocumentFiltersDTO) ([]dto.DocumentDTO, error) {
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

	fmt.Printf("DEBUG: DocumentService.ListAllDocuments calling repo with tenantID=%s, filters=%v\n", tenantID, filterMap)
	documents, err := s.documentRepo.ListAllDocuments(ctx, tenantID, filterMap)
	if err != nil {
		fmt.Printf("ERROR: DocumentService.ListAllDocuments repo error: %v\n", err)
		return nil, fmt.Errorf("failed to list all documents: %w", err)
	}
	fmt.Printf("DEBUG: DocumentService.ListAllDocuments got %d documents\n", len(documents))

	result := make([]dto.DocumentDTO, 0)
	for _, document := range documents {
		// Get tags and links for each document
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
			Title:        document.Title,
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
