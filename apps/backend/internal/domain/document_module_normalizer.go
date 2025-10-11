package domain

import (
	"strings"

	"risknexus/backend/internal/dto"
)

// normalizeModuleName maps different module aliases to canonical storage names.
func normalizeModuleName(raw string) string {
	name := strings.TrimSpace(strings.ToLower(raw))
	if name == "" {
		return ""
	}

	// Strip optional leading hash (tags can include it).
	if strings.HasPrefix(name, "#") {
		name = strings.TrimPrefix(name, "#")
	}

	switch name {
	case "asset", "assets", "актив", "активы":
		return "assets"
	case "risk", "risks", "риск", "риски", "управление рисками":
		return "risks"
	case "incident", "incidents", "инцидент", "инциденты":
		return "incidents"
	case "training", "trainings", "обучение":
		return "training"
	case "compliance", "соответствие", "политики":
		return "compliance"
	case "audit", "audits", "аудит", "аудиты":
		return "audits"
	case "document", "documents", "docs", "документ", "документы":
		return "documents"
	case "general", "общие":
		return "general"
	default:
		return name
	}
}

// resolveModuleFromContext determines the canonical module for storage and linking.
func (s *DocumentService) resolveModuleFromContext(req dto.UploadDocumentDTO) string {
	if req.LinkedTo != nil && req.LinkedTo.Module != "" {
		if normalized := normalizeModuleName(req.LinkedTo.Module); normalized != "" {
			return normalized
		}
		return req.LinkedTo.Module
	}

	for _, tag := range req.Tags {
		if normalized := normalizeModuleName(tag); normalized != "" && normalized != "general" {
			return normalized
		}
	}

	return "general"
}
