package domain

import (
	"context"
	"time"

	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

type AssetService struct {
	assetRepo *repo.AssetRepo
	auditRepo *repo.AuditRepo
}

func NewAssetService(assetRepo *repo.AssetRepo, auditRepo *repo.AuditRepo) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
		auditRepo: auditRepo,
	}
}

func (s *AssetService) CreateAsset(ctx context.Context, tenantID, name, assetType string, invCode, ownerID, location, software *string) (*repo.Asset, error) {
	asset := repo.Asset{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		Name:      name,
		InvCode:   invCode,
		Type:      assetType,
		Status:    "active",
		OwnerID:   ownerID,
		Location:  location,
		Software:  software,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.assetRepo.Create(ctx, asset)
	if err != nil {
		return nil, err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, tenantID, "system", "create", "asset", &asset.ID, asset)

	return &asset, nil
}

func (s *AssetService) GetAsset(ctx context.Context, id string) (*repo.Asset, error) {
	return s.assetRepo.GetByID(ctx, id)
}

func (s *AssetService) ListAssets(ctx context.Context, tenantID string) ([]repo.Asset, error) {
	return s.assetRepo.List(ctx, tenantID)
}

func (s *AssetService) UpdateAsset(ctx context.Context, id, name, assetType string, invCode, ownerID, location, software *string) error {
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if asset == nil {
		return nil
	}

	asset.Name = name
	asset.InvCode = invCode
	asset.Type = assetType
	asset.OwnerID = ownerID
	asset.Location = location
	asset.Software = software
	asset.UpdatedAt = time.Now()

	err = s.assetRepo.Update(ctx, *asset)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, asset.TenantID, "system", "update", "asset", &id, asset)

	return nil
}

func (s *AssetService) DeleteAsset(ctx context.Context, id string) error {
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if asset == nil {
		return nil
	}

	err = s.assetRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Log audit
	s.auditRepo.LogAction(ctx, asset.TenantID, "system", "delete", "asset", &id, nil)

	return nil
}
