//go:build legacy_tests
// +build legacy_tests

package main

import (
	"context"
	"testing"
	"time"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAssetRepo - РјРѕРє РґР»СЏ AssetRepo
type MockAssetRepo struct {
	mock.Mock
}

// РЈР±РµР¶РґР°РµРјСЃСЏ, С‡С‚Рѕ MockAssetRepo СЂРµР°Р»РёР·СѓРµС‚ РёРЅС‚РµСЂС„РµР№СЃ
var _ domain.AssetRepoInterface = (*MockAssetRepo)(nil)

func (m *MockAssetRepo) Create(ctx context.Context, asset repo.Asset) error {
	arguments := m.Called(ctx, asset)
	return arguments.Error(0)
}

func (m *MockAssetRepo) GetByID(ctx context.Context, id string) (*repo.Asset, error) {
	arguments := m.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.Asset), arguments.Error(1)
}

func (m *MockAssetRepo) List(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.Asset, error) {
	arguments := m.Called(ctx, tenantID, filters)
	return arguments.Get(0).([]repo.Asset), arguments.Error(1)
}

func (m *MockAssetRepo) ListPaginated(ctx context.Context, tenantID string, page, pageSize int, filters map[string]interface{}) ([]repo.Asset, int64, error) {
	arguments := m.Called(ctx, tenantID, page, pageSize, filters)
	return arguments.Get(0).([]repo.Asset), arguments.Get(1).(int64), arguments.Error(2)
}

func (m *MockAssetRepo) Update(ctx context.Context, asset repo.Asset) error {
	arguments := m.Called(ctx, asset)
	return arguments.Error(0)
}

func (m *MockAssetRepo) SoftDelete(ctx context.Context, id string) error {
	arguments := m.Called(ctx, id)
	return arguments.Error(0)
}

func (m *MockAssetRepo) GetWithDetails(ctx context.Context, id string) (*repo.AssetWithDetails, error) {
	arguments := m.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.AssetWithDetails), arguments.Error(1)
}

func (m *MockAssetRepo) AddDocument(ctx context.Context, assetID, documentType, filePath, createdBy string) error {
	arguments := m.Called(ctx, assetID, documentType, filePath, createdBy)
	return arguments.Error(0)
}

func (m *MockAssetRepo) GetAssetDocuments(ctx context.Context, assetID string) ([]repo.AssetDocument, error) {
	arguments := m.Called(ctx, assetID)
	return arguments.Get(0).([]repo.AssetDocument), arguments.Error(1)
}

func (m *MockAssetRepo) AddSoftware(ctx context.Context, assetID, softwareName, version string, installedAt *time.Time) error {
	arguments := m.Called(ctx, assetID, softwareName, version, installedAt)
	return arguments.Error(0)
}

func (m *MockAssetRepo) GetAssetSoftware(ctx context.Context, assetID string) ([]repo.AssetSoftware, error) {
	arguments := m.Called(ctx, assetID)
	return arguments.Get(0).([]repo.AssetSoftware), arguments.Error(1)
}

func (m *MockAssetRepo) AddHistory(ctx context.Context, assetID, fieldChanged, oldValue, newValue, changedBy string) error {
	arguments := m.Called(ctx, assetID, fieldChanged, oldValue, newValue, changedBy)
	return arguments.Error(0)
}

func (m *MockAssetRepo) GetAssetHistory(ctx context.Context, assetID string) ([]repo.AssetHistory, error) {
	arguments := m.Called(ctx, assetID)
	return arguments.Get(0).([]repo.AssetHistory), arguments.Error(1)
}

// MockUserRepo - РјРѕРє РґР»СЏ UserRepo
type MockUserRepo struct {
	mock.Mock
}

// РЈР±РµР¶РґР°РµРјСЃСЏ, С‡С‚Рѕ MockUserRepo СЂРµР°Р»РёР·СѓРµС‚ РёРЅС‚РµСЂС„РµР№СЃ
var _ domain.UserRepoInterface = (*MockUserRepo)(nil)

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*repo.User, error) {
	arguments := m.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.User), arguments.Error(1)
}

func TestAssetService_CreateAsset(t *testing.T) {
	mockAssetRepo := new(MockAssetRepo)
	mockUserRepo := new(MockUserRepo)
	service := domain.NewAssetService(mockAssetRepo, mockUserRepo)

	ctx := context.Background()
	tenantID := "tenant-123"
	createdBy := "user-123"
	req := dto.CreateAssetRequest{
		Name:            "Test Server",
		Type:            "server",
		Class:           "hardware",
		OwnerID:         "owner-123",
		Location:        "Data Center A",
		Criticality:     "high",
		Confidentiality: "high",
		Integrity:       "high",
		Availability:    "high",
		Status:          "active",
	}

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєРѕРІ
	owner := &repo.User{
		ID:        "owner-123",
		FirstName: stringPtr("John"),
		LastName:  stringPtr("Doe"),
	}
	mockUserRepo.On("GetByID", ctx, "owner-123").Return(owner, nil)
	mockAssetRepo.On("Create", ctx, mock.AnythingOfType("repo.Asset")).Return(nil)
	mockAssetRepo.On("AddHistory", ctx, mock.AnythingOfType("string"), "created", "", "Asset created", createdBy).Return(nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	asset, err := service.CreateAsset(ctx, tenantID, req, createdBy)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, req.Name, asset.Name)
	assert.Equal(t, req.Type, asset.Type)
	assert.Equal(t, req.Class, asset.Class)
	assert.Equal(t, req.OwnerID, *asset.OwnerID)
	assert.Equal(t, req.Criticality, asset.Criticality)
	assert.Equal(t, "active", asset.Status) // РЎС‚Р°С‚СѓСЃ РїРѕ СѓРјРѕР»С‡Р°РЅРёСЋ

	mockAssetRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestAssetService_CreateAsset_OwnerNotFound(t *testing.T) {
	mockAssetRepo := new(MockAssetRepo)
	mockUserRepo := new(MockUserRepo)
	service := domain.NewAssetService(mockAssetRepo, mockUserRepo)

	ctx := context.Background()
	tenantID := "tenant-123"
	createdBy := "user-123"
	req := dto.CreateAssetRequest{
		Name:            "Test Server",
		Type:            "server",
		Class:           "hardware",
		OwnerID:         "owner-123",
		Criticality:     "high",
		Confidentiality: "high",
		Integrity:       "high",
		Availability:    "high",
	}

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєР° - РІР»Р°РґРµР»РµС† РЅРµ РЅР°Р№РґРµРЅ
	mockUserRepo.On("GetByID", ctx, "owner-123").Return(nil, nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	asset, err := service.CreateAsset(ctx, tenantID, req, createdBy)

	// РџСЂРѕРІРµСЂРєРё
	assert.Error(t, err)
	assert.Nil(t, asset)
	assert.Contains(t, err.Error(), "owner not found")

	mockUserRepo.AssertExpectations(t)
	mockAssetRepo.AssertNotCalled(t, "Create")
}

func TestAssetService_UpdateAsset(t *testing.T) {
	mockAssetRepo := new(MockAssetRepo)
	mockUserRepo := new(MockUserRepo)
	service := domain.NewAssetService(mockAssetRepo, mockUserRepo)

	ctx := context.Background()
	assetID := uuid.New().String()
	updatedBy := "user-123"

	// РЎСѓС‰РµСЃС‚РІСѓСЋС‰РёР№ Р°РєС‚РёРІ
	existingAsset := &repo.Asset{
		ID:              assetID,
		Name:            "Old Name",
		Type:            "server",
		Class:           "hardware",
		Criticality:     "medium",
		Confidentiality: "medium",
		Integrity:       "medium",
		Availability:    "medium",
		Status:          "active",
	}

	req := dto.UpdateAssetRequest{
		Name:        stringPtr("New Name"),
		Criticality: stringPtr("high"),
	}

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєРѕРІ
	mockAssetRepo.On("GetByID", ctx, assetID).Return(existingAsset, nil)
	mockAssetRepo.On("Update", ctx, mock.AnythingOfType("repo.Asset")).Return(nil)
	mockAssetRepo.On("AddHistory", ctx, assetID, "name", "Old Name", "New Name", updatedBy).Return(nil)
	mockAssetRepo.On("AddHistory", ctx, assetID, "criticality", "medium", "high", updatedBy).Return(nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	err := service.UpdateAsset(ctx, assetID, req, updatedBy)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
}

func TestAssetService_UpdateAsset_NotFound(t *testing.T) {
	mockAssetRepo := new(MockAssetRepo)
	mockUserRepo := new(MockUserRepo)
	service := domain.NewAssetService(mockAssetRepo, mockUserRepo)

	ctx := context.Background()
	assetID := uuid.New().String()
	updatedBy := "user-123"
	req := dto.UpdateAssetRequest{
		Name: stringPtr("New Name"),
	}

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєР° - Р°РєС‚РёРІ РЅРµ РЅР°Р№РґРµРЅ
	mockAssetRepo.On("GetByID", ctx, assetID).Return(nil, nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	err := service.UpdateAsset(ctx, assetID, req, updatedBy)

	// РџСЂРѕРІРµСЂРєРё
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "asset not found")
	mockAssetRepo.AssertExpectations(t)
}

func TestAssetService_DeleteAsset(t *testing.T) {
	mockAssetRepo := new(MockAssetRepo)
	mockUserRepo := new(MockUserRepo)
	service := domain.NewAssetService(mockAssetRepo, mockUserRepo)

	ctx := context.Background()
	assetID := uuid.New().String()
	deletedBy := "user-123"

	// РЎСѓС‰РµСЃС‚РІСѓСЋС‰РёР№ Р°РєС‚РёРІ
	existingAsset := &repo.Asset{
		ID:     assetID,
		Name:   "Test Asset",
		Status: "active",
	}

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєРѕРІ
	mockAssetRepo.On("GetByID", ctx, assetID).Return(existingAsset, nil)
	mockAssetRepo.On("SoftDelete", ctx, assetID).Return(nil)
	mockAssetRepo.On("AddHistory", ctx, assetID, "deleted", "", "Asset deleted", deletedBy).Return(nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	err := service.DeleteAsset(ctx, assetID, deletedBy)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
}

func TestAssetService_AddDocument(t *testing.T) {
	mockAssetRepo := new(MockAssetRepo)
	mockUserRepo := new(MockUserRepo)
	service := domain.NewAssetService(mockAssetRepo, mockUserRepo)

	ctx := context.Background()
	assetID := uuid.New().String()
	createdBy := "user-123"
	req := dto.AssetDocumentRequest{
		DocumentType: "passport",
		FilePath:     "/path/to/document.pdf",
	}

	// РЎСѓС‰РµСЃС‚РІСѓСЋС‰РёР№ Р°РєС‚РёРІ
	existingAsset := &repo.Asset{
		ID:     assetID,
		Name:   "Test Asset",
		Status: "active",
	}

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєРѕРІ
	mockAssetRepo.On("GetByID", ctx, assetID).Return(existingAsset, nil)
	mockAssetRepo.On("AddDocument", ctx, assetID, "passport", "/path/to/document.pdf", createdBy).Return(nil)
	mockAssetRepo.On("AddHistory", ctx, assetID, "document_added", "", "passport", createdBy).Return(nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	err := service.AddDocument(ctx, assetID, req, createdBy)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
}

func TestAssetService_AddSoftware(t *testing.T) {
	mockAssetRepo := new(MockAssetRepo)
	mockUserRepo := new(MockUserRepo)
	service := domain.NewAssetService(mockAssetRepo, mockUserRepo)

	ctx := context.Background()
	assetID := uuid.New().String()
	addedBy := "user-123"
	installedAt := time.Now()
	req := dto.AssetSoftwareRequest{
		SoftwareName: "Windows Server 2019",
		Version:      stringPtr("10.0.17763"),
		InstalledAt:  &installedAt,
	}

	// РЎСѓС‰РµСЃС‚РІСѓСЋС‰РёР№ Р°РєС‚РёРІ
	existingAsset := &repo.Asset{
		ID:     assetID,
		Name:   "Test Asset",
		Status: "active",
	}

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєРѕРІ
	mockAssetRepo.On("GetByID", ctx, assetID).Return(existingAsset, nil)
	mockAssetRepo.On("AddSoftware", ctx, assetID, "Windows Server 2019", "10.0.17763", &installedAt).Return(nil)
	mockAssetRepo.On("AddHistory", ctx, assetID, "software_added", "", "Windows Server 2019", addedBy).Return(nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	err := service.AddSoftware(ctx, assetID, req, addedBy)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
}

func TestAssetService_PerformInventory(t *testing.T) {
	mockAssetRepo := new(MockAssetRepo)
	mockUserRepo := new(MockUserRepo)
	service := domain.NewAssetService(mockAssetRepo, mockUserRepo)

	ctx := context.Background()
	tenantID := "tenant-123"
	performedBy := "user-123"
	req := dto.AssetInventoryRequest{
		AssetIDs: []string{"asset-1", "asset-2"},
		Action:   "update_status",
		Status:   stringPtr("in_repair"),
		Notes:    stringPtr("Maintenance required"),
	}

	// РЎСѓС‰РµСЃС‚РІСѓСЋС‰РёРµ Р°РєС‚РёРІС‹
	asset1 := &repo.Asset{ID: "asset-1", Name: "Asset 1", Status: "active"}
	asset2 := &repo.Asset{ID: "asset-2", Name: "Asset 2", Status: "active"}

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєРѕРІ
	mockAssetRepo.On("GetByID", ctx, "asset-1").Return(asset1, nil)
	mockAssetRepo.On("GetByID", ctx, "asset-2").Return(asset2, nil)
	mockAssetRepo.On("Update", ctx, mock.AnythingOfType("repo.Asset")).Return(nil).Twice()
	mockAssetRepo.On("AddHistory", ctx, "asset-1", "status", "active", "in_repair", performedBy).Return(nil)
	mockAssetRepo.On("AddHistory", ctx, "asset-2", "status", "active", "in_repair", performedBy).Return(nil)
	mockAssetRepo.On("AddHistory", ctx, "asset-1", "inventory_update_status", "", "Maintenance required", performedBy).Return(nil)
	mockAssetRepo.On("AddHistory", ctx, "asset-2", "inventory_update_status", "", "Maintenance required", performedBy).Return(nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	err := service.PerformInventory(ctx, tenantID, req, performedBy)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	mockAssetRepo.AssertExpectations(t)
}
