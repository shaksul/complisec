package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

type TenantService struct {
	tenantRepo *repo.TenantRepo
	auditRepo  *repo.AuditRepo
}

func NewTenantService(tenantRepo *repo.TenantRepo, auditRepo *repo.AuditRepo) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		auditRepo:  auditRepo,
	}
}

func (s *TenantService) CreateTenant(ctx context.Context, req dto.CreateTenantDTO, createdBy string) (*dto.TenantResponse, error) {
	// Проверяем уникальность домена, если он указан
	if req.Domain != nil && *req.Domain != "" {
		exists, err := s.tenantRepo.ExistsByDomain(ctx, *req.Domain, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check domain uniqueness: %w", err)
		}
		if exists {
			return nil, errors.New("domain already exists")
		}
	}

	// Создаем организацию
	tenant := repo.Tenant{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Domain:    req.Domain,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	// Логируем создание
	if err := s.auditRepo.LogAction(ctx, tenant.ID, createdBy, "created", "tenant", &tenant.ID, map[string]interface{}{
		"name":   tenant.Name,
		"domain": tenant.Domain,
	}); err != nil {
		// Не критично, продолжаем
	}

	return &dto.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Domain:    tenant.Domain,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}, nil
}

func (s *TenantService) GetTenant(ctx context.Context, id string) (*dto.TenantResponse, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	if tenant == nil {
		return nil, errors.New("tenant not found")
	}

	return &dto.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Domain:    tenant.Domain,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}, nil
}

func (s *TenantService) GetTenantByDomain(ctx context.Context, domain string) (*dto.TenantResponse, error) {
	tenant, err := s.tenantRepo.GetByDomain(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by domain: %w", err)
	}
	if tenant == nil {
		return nil, errors.New("tenant not found")
	}

	return &dto.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Domain:    tenant.Domain,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}, nil
}

func (s *TenantService) ListTenants(ctx context.Context, page, pageSize int) (*dto.TenantListResponse, error) {
	offset := (page - 1) * pageSize

	// Получаем список организаций
	tenants, err := s.tenantRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	// Получаем общее количество
	total, err := s.tenantRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count tenants: %w", err)
	}

	// Преобразуем в DTO
	tenantResponses := make([]dto.TenantResponse, len(tenants))
	for i, tenant := range tenants {
		tenantResponses[i] = dto.TenantResponse{
			ID:        tenant.ID,
			Name:      tenant.Name,
			Domain:    tenant.Domain,
			CreatedAt: tenant.CreatedAt,
			UpdatedAt: tenant.UpdatedAt,
		}
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &dto.TenantListResponse{
		Data: tenantResponses,
		Pagination: dto.PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// ListAllTenants - алиас для ListTenants для администраторов
func (s *TenantService) ListAllTenants(ctx context.Context, page, pageSize int) (*dto.TenantListResponse, error) {
	return s.ListTenants(ctx, page, pageSize)
}

func (s *TenantService) UpdateTenant(ctx context.Context, id string, req dto.UpdateTenantDTO, updatedBy string) (*dto.TenantResponse, error) {
	// Проверяем существование организации
	existingTenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	if existingTenant == nil {
		return nil, errors.New("tenant not found")
	}

	// Проверяем уникальность домена, если он изменился
	if req.Domain != nil && *req.Domain != "" {
		if existingTenant.Domain == nil || *existingTenant.Domain != *req.Domain {
			exists, err := s.tenantRepo.ExistsByDomain(ctx, *req.Domain, &id)
			if err != nil {
				return nil, fmt.Errorf("failed to check domain uniqueness: %w", err)
			}
			if exists {
				return nil, errors.New("domain already exists")
			}
		}
	}

	// Обновляем организацию
	if err := s.tenantRepo.Update(ctx, id, req.Name, req.Domain); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	// Логируем обновление
	if err := s.auditRepo.LogAction(ctx, id, updatedBy, "updated", "tenant", &id, map[string]interface{}{
		"name":   req.Name,
		"domain": req.Domain,
	}); err != nil {
		// Не критично, продолжаем
	}

	// Возвращаем обновленную организацию
	return s.GetTenant(ctx, id)
}

func (s *TenantService) DeleteTenant(ctx context.Context, id string, deletedBy string) error {
	// Проверяем существование организации
	tenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}
	if tenant == nil {
		return errors.New("tenant not found")
	}

	// Удаляем организацию
	if err := s.tenantRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	// Логируем удаление
	if err := s.auditRepo.LogAction(ctx, id, deletedBy, "deleted", "tenant", &id, map[string]interface{}{
		"name":   tenant.Name,
		"domain": tenant.Domain,
	}); err != nil {
		// Не критично, продолжаем
	}

	return nil
}
