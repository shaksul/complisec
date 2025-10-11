package domain

import (
	"strings"

	"risknexus/backend/internal/dto"
)

// resolveCategoryFromRequest normalizes the category selection, falling back to the first tag.
func (s *DocumentService) resolveCategoryFromRequest(req dto.UploadDocumentDTO) string {
	category := s.detectCategoryFromContext(req)
	if category != "" && category != "uncategorized" {
		return category
	}

	for _, raw := range req.Tags {
		clean := strings.TrimSpace(raw)
		if clean == "" {
			continue
		}

		if strings.HasPrefix(clean, "#") {
			clean = strings.TrimSpace(strings.TrimPrefix(clean, "#"))
		}

		if clean != "" {
			return clean
		}
	}

	return category
}
