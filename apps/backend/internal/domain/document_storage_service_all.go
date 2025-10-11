package domain

import (
	"context"

	"risknexus/backend/internal/dto"
)

// ListAllDocuments gets ALL documents including those linked to modules (for file storage)
func (s *DocumentStorageService) ListAllDocuments(ctx context.Context, tenantID string, filters dto.FileDocumentFiltersDTO) ([]dto.DocumentDTO, error) {
	return s.documentService.ListAllDocuments(ctx, tenantID, filters)
}
