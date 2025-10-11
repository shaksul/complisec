//go:build legacy_tests
// +build legacy_tests

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"
	httpHandler "risknexus/backend/internal/http"
	"risknexus/backend/internal/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAssetService - РјРѕРє РґР»СЏ AssetService
type MockAssetService struct {
	mock.Mock
}

// РЈР±РµР¶РґР°РµРјСЃСЏ, С‡С‚Рѕ MockAssetService СЂРµР°Р»РёР·СѓРµС‚ РёРЅС‚РµСЂС„РµР№СЃ
var _ domain.AssetServiceInterface = (*MockAssetService)(nil)

// Р РµР°Р»РёР·СѓРµРј РёРЅС‚РµСЂС„РµР№СЃ AssetService
func (m *MockAssetService) CreateAsset(ctx context.Context, tenantID string, req dto.CreateAssetRequest, createdBy string) (*repo.Asset, error) {
	arguments := m.Called(ctx, tenantID, req, createdBy)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.Asset), arguments.Error(1)
}

func (m *MockAssetService) GetAsset(ctx context.Context, id string) (*repo.Asset, error) {
	arguments := m.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.Asset), arguments.Error(1)
}

func (m *MockAssetService) GetAssetWithDetails(ctx context.Context, id string) (*repo.AssetWithDetails, error) {
	arguments := m.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.AssetWithDetails), arguments.Error(1)
}

func (m *MockAssetService) ListAssets(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.Asset, error) {
	arguments := m.Called(ctx, tenantID, filters)
	return arguments.Get(0).([]repo.Asset), arguments.Error(1)
}

func (m *MockAssetService) ListAssetsPaginated(ctx context.Context, tenantID string, page, pageSize int, filters map[string]interface{}) ([]repo.Asset, int64, error) {
	arguments := m.Called(ctx, tenantID, page, pageSize, filters)
	return arguments.Get(0).([]repo.Asset), arguments.Get(1).(int64), arguments.Error(2)
}

func (m *MockAssetService) UpdateAsset(ctx context.Context, id string, req dto.UpdateAssetRequest, updatedBy string) error {
	arguments := m.Called(ctx, id, req, updatedBy)
	return arguments.Error(0)
}

func (m *MockAssetService) DeleteAsset(ctx context.Context, id string, deletedBy string) error {
	arguments := m.Called(ctx, id, deletedBy)
	return arguments.Error(0)
}

func (m *MockAssetService) AddDocument(ctx context.Context, assetID string, req dto.AssetDocumentRequest, createdBy string) error {
	arguments := m.Called(ctx, assetID, req, createdBy)
	return arguments.Error(0)
}

func (m *MockAssetService) GetAssetDocuments(ctx context.Context, assetID string) ([]repo.AssetDocument, error) {
	arguments := m.Called(ctx, assetID)
	return arguments.Get(0).([]repo.AssetDocument), arguments.Error(1)
}

func (m *MockAssetService) AddSoftware(ctx context.Context, assetID string, req dto.AssetSoftwareRequest, addedBy string) error {
	arguments := m.Called(ctx, assetID, req, addedBy)
	return arguments.Error(0)
}

func (m *MockAssetService) GetAssetSoftware(ctx context.Context, assetID string) ([]repo.AssetSoftware, error) {
	arguments := m.Called(ctx, assetID)
	return arguments.Get(0).([]repo.AssetSoftware), arguments.Error(1)
}

func (m *MockAssetService) GetAssetHistory(ctx context.Context, assetID string) ([]repo.AssetHistory, error) {
	arguments := m.Called(ctx, assetID)
	return arguments.Get(0).([]repo.AssetHistory), arguments.Error(1)
}

func (m *MockAssetService) PerformInventory(ctx context.Context, tenantID string, req dto.AssetInventoryRequest, performedBy string) error {
	arguments := m.Called(ctx, tenantID, req, performedBy)
	return arguments.Error(0)
}

func TestAssetHandler_CreateAsset(t *testing.T) {
	mockService := new(MockAssetService)
	handler := httpHandler.NewAssetHandler(mockService)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", "tenant-123")
		c.Locals("user_id", "user-123")
		c.Locals("roles", []string{"Admin"})
		return c.Next()
	})
	handler.Register(app)

	req := dto.CreateAssetRequest{
		Name:            "Test Server",
		Type:            "server",
		Class:           "hardware",
		OwnerID:         "123e4567-e89b-12d3-a456-426614174000",
		Location:        "Data Center A",
		Criticality:     "high",
		Confidentiality: "high",
		Integrity:       "high",
		Availability:    "high",
		Status:          "active",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/assets", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєР°
	expectedAsset := &repo.Asset{
		ID:              "asset-123",
		TenantID:        "tenant-123",
		InventoryNumber: "AST-20241201-ABC12345",
		Name:            "Test Server",
		Type:            "server",
		Class:           "hardware",
		OwnerID:         stringPtr("123e4567-e89b-12d3-a456-426614174000"),
		Location:        stringPtr("Data Center A"),
		Criticality:     "high",
		Confidentiality: "high",
		Integrity:       "high",
		Availability:    "high",
		Status:          "active",
	}

	mockService.On("CreateAsset", mock.Anything, "tenant-123", req, "user-123").Return(expectedAsset, nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	resp, err := app.Test(httpReq)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response, "data")

	mockService.AssertExpectations(t)
}

func TestAssetHandler_CreateAsset_ValidationError(t *testing.T) {
	mockService := new(MockAssetService)
	handler := httpHandler.NewAssetHandler(mockService)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", "tenant-123")
		c.Locals("user_id", "user-123")
		c.Locals("roles", []string{"Admin"})
		return c.Next()
	})
	handler.Register(app)

	// РќРµРІР°Р»РёРґРЅС‹Р№ Р·Р°РїСЂРѕСЃ - РѕС‚СЃСѓС‚СЃС‚РІСѓРµС‚ РѕР±СЏР·Р°С‚РµР»СЊРЅРѕРµ РїРѕР»Рµ
	req := dto.CreateAssetRequest{
		Type:        "server",
		Class:       "hardware",
		Criticality: "high",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/assets", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	resp, err := app.Test(httpReq)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response, "error")

	mockService.AssertNotCalled(t, "CreateAsset")
}

func TestAssetHandler_GetAsset(t *testing.T) {
	mockService := new(MockAssetService)
	handler := httpHandler.NewAssetHandler(mockService)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", "tenant-123")
		c.Locals("user_id", "user-123")
		c.Locals("roles", []string{"Admin"})
		return c.Next()
	})
	handler.Register(app)

	assetID := "asset-123"
	httpReq := httptest.NewRequest("GET", "/assets/"+assetID, nil)

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєР°
	expectedAsset := &repo.Asset{
		ID:              assetID,
		TenantID:        "tenant-123",
		InventoryNumber: "AST-20241201-ABC12345",
		Name:            "Test Server",
		Type:            "server",
		Class:           "hardware",
		Criticality:     "high",
		Confidentiality: "high",
		Integrity:       "high",
		Availability:    "high",
		Status:          "active",
	}

	mockService.On("GetAsset", mock.Anything, assetID).Return(expectedAsset, nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	resp, err := app.Test(httpReq)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response, "data")

	mockService.AssertExpectations(t)
}

func TestAssetHandler_GetAsset_NotFound(t *testing.T) {
	mockService := new(MockAssetService)
	handler := httpHandler.NewAssetHandler(mockService)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", "tenant-123")
		c.Locals("user_id", "user-123")
		c.Locals("roles", []string{"Admin"})
		return c.Next()
	})
	handler.Register(app)

	assetID := "asset-123"
	httpReq := httptest.NewRequest("GET", "/assets/"+assetID, nil)

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєР° - Р°РєС‚РёРІ РЅРµ РЅР°Р№РґРµРЅ
	mockService.On("GetAsset", mock.Anything, assetID).Return(nil, nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	resp, err := app.Test(httpReq)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response, "error")

	mockService.AssertExpectations(t)
}

func TestAssetHandler_ListAssets(t *testing.T) {
	mockService := new(MockAssetService)
	handler := httpHandler.NewAssetHandler(mockService)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", "tenant-123")
		c.Locals("user_id", "user-123")
		c.Locals("roles", []string{"Admin"})
		return c.Next()
	})
	handler.Register(app)

	httpReq := httptest.NewRequest("GET", "/assets?page=1&page_size=10&type=server", nil)

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєР°
	expectedAssets := []repo.Asset{
		{
			ID:              "asset-1",
			TenantID:        "tenant-123",
			InventoryNumber: "AST-20241201-ABC12345",
			Name:            "Server 1",
			Type:            "server",
			Class:           "hardware",
			Criticality:     "high",
			Confidentiality: "high",
			Integrity:       "high",
			Availability:    "high",
			Status:          "active",
		},
	}

	expectedTotal := int64(1)
	expectedFilters := map[string]interface{}{
		"type": "server",
	}

	mockService.On("ListAssetsPaginated", mock.Anything, "tenant-123", 1, 10, expectedFilters).Return(expectedAssets, expectedTotal, nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	resp, err := app.Test(httpReq)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.PaginatedResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, expectedTotal, response.Pagination.Total)

	mockService.AssertExpectations(t)
}

func TestAssetHandler_UpdateAsset(t *testing.T) {
	mockService := new(MockAssetService)
	handler := httpHandler.NewAssetHandler(mockService)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", "tenant-123")
		c.Locals("user_id", "user-123")
		c.Locals("roles", []string{"Admin"})
		return c.Next()
	})
	handler.Register(app)

	assetID := "asset-123"
	req := dto.UpdateAssetRequest{
		Name:        stringPtr("Updated Server"),
		Criticality: stringPtr("high"),
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", "/assets/"+assetID, bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєР°
	mockService.On("UpdateAsset", mock.Anything, assetID, req, "user-123").Return(nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	resp, err := app.Test(httpReq)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response, "message")

	mockService.AssertExpectations(t)
}

func TestAssetHandler_DeleteAsset(t *testing.T) {
	mockService := new(MockAssetService)
	handler := httpHandler.NewAssetHandler(mockService)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", "tenant-123")
		c.Locals("user_id", "user-123")
		c.Locals("roles", []string{"Admin"})
		return c.Next()
	})
	handler.Register(app)

	assetID := "asset-123"
	httpReq := httptest.NewRequest("DELETE", "/assets/"+assetID, nil)

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєР°
	mockService.On("DeleteAsset", mock.Anything, assetID, "user-123").Return(nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	resp, err := app.Test(httpReq)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response, "message")

	mockService.AssertExpectations(t)
}

func TestAssetHandler_ExportAssets(t *testing.T) {
	mockService := new(MockAssetService)
	handler := httpHandler.NewAssetHandler(mockService)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", "tenant-123")
		c.Locals("user_id", "user-123")
		c.Locals("roles", []string{"Admin"})
		return c.Next()
	})
	handler.Register(app)

	httpReq := httptest.NewRequest("GET", "/assets/export?type=server", nil)

	// РќР°СЃС‚СЂРѕР№РєР° РјРѕРєР°
	expectedAssets := []repo.Asset{
		{
			ID:              "asset-1",
			TenantID:        "tenant-123",
			InventoryNumber: "AST-20241201-ABC12345",
			Name:            "Server 1",
			Type:            "server",
			Class:           "hardware",
			Criticality:     "high",
			Confidentiality: "high",
			Integrity:       "high",
			Availability:    "high",
			Status:          "active",
		},
	}

	expectedFilters := map[string]interface{}{
		"type": "server",
	}

	mockService.On("ListAssets", mock.Anything, "tenant-123", expectedFilters).Return(expectedAssets, nil)

	// Р’С‹РїРѕР»РЅРµРЅРёРµ С‚РµСЃС‚Р°
	resp, err := app.Test(httpReq)

	// РџСЂРѕРІРµСЂРєРё
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))
	assert.Contains(t, resp.Header.Get("Content-Disposition"), "attachment")

	mockService.AssertExpectations(t)
}

// Р’СЃРїРѕРјРѕРіР°С‚РµР»СЊРЅР°СЏ С„СѓРЅРєС†РёСЏ
func stringPtr(s string) *string {
	return &s
}
